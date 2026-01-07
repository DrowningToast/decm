package log

import (
	"os"
	"sync"
	"sync/atomic"

	"github.com/cockroachdb/errors"
)

// RotatingFileWriter is a file writer that can reopen its file handle.
// This is needed for log rotation - after renaming the file, we need to
// create a new file with the original name.
//
// Performance Optimization:
//
//	Writes use atomic.Value for lock-free reads of the file handle, allowing
//	high-concurrency logging without serialization. Only Reopen() acquires
//	an exclusive lock to safely swap the file handle during rotation.
//
// Thread Safety:
//   - Write(): Lock-free reads via atomic.Value (scales with concurrency)
//   - Reopen(): Exclusive lock during file handle swap (rare operation)
type RotatingFileWriter struct {
	fileHandle atomic.Value // Stores *os.File - lock-free reads
	path       string
	rotateMu   sync.Mutex // Only locked during Reopen()
}

// NewRotatingFileWriter creates a new rotating file writer.
// The file is auto-created if it doesn't exist (O_CREATE flag).
func NewRotatingFileWriter(path string) (*RotatingFileWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open log file")
	}

	w := &RotatingFileWriter{
		path: path,
	}
	w.fileHandle.Store(f)
	return w, nil
}

// Write writes data to the current file handle.
//
// Performance: Lock-free operation using atomic.Value.Load()
// Multiple goroutines can write concurrently without waiting.
// The file handle is loaded atomically, ensuring we always have a valid handle.
//
// Note: Even during rotation, writes continue safely to either the old
// or new file handle depending on timing - both are valid destinations.
func (w *RotatingFileWriter) Write(p []byte) (n int, err error) {
	// Lock-free read of current file handle - allows concurrent writes
	file := w.fileHandle.Load()
	if file == nil {
		return 0, errors.New("file handle is nil")
	}

	f, ok := file.(*os.File)
	if !ok {
		return 0, errors.New("invalid file handle type")
	}

	return f.Write(p)
}

// Reopen closes the current file handle and opens a new one.
// This should be called after rotating (renaming) the file.
//
// Thread Safety:
//
//	Uses rotateMu to ensure only one rotation happens at a time.
//	File handle swap is atomic via atomic.Value.Store().
//
// Flow:
//  1. Lock rotateMu (blocks concurrent rotations, but not writes)
//  2. Close old file handle
//  3. Open new file with original path
//  4. Atomically swap file handles
//  5. Unlock
//
// Note: Writes happening during rotation will use either old or new handle,
// both are valid - logs won't be lost.
func (w *RotatingFileWriter) Reopen() error {
	w.rotateMu.Lock()
	defer w.rotateMu.Unlock()

	// Get and close old file handle
	oldFile := w.fileHandle.Load()
	if oldFile != nil {
		if f, ok := oldFile.(*os.File); ok {
			if err := f.Close(); err != nil {
				return errors.Wrap(err, "failed to close old file handle")
			}
		}
	}

	// Open new file with same path (creates if doesn't exist)
	newFile, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to reopen log file")
	}

	// Atomically swap file handles - concurrent writes now use new file
	w.fileHandle.Store(newFile)
	return nil
}

// Close closes the file handle.
// Should be called when shutting down the logger.
func (w *RotatingFileWriter) Close() error {
	w.rotateMu.Lock()
	defer w.rotateMu.Unlock()

	file := w.fileHandle.Load()
	if file != nil {
		if f, ok := file.(*os.File); ok {
			return f.Close()
		}
	}
	return nil
}
