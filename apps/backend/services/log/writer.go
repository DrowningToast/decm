package log

import (
	"io"
	"sync"

	"github.com/cockroachdb/errors"
)

// Write Hook System
//
// This file implements a write hook mechanism that allows intercepting all log writes
// before they are written to the underlying writer (stdout + file).
//
// Architecture:
//   User Code → Logger → HookWriter.Write() → [Hook Function] → Actual Writers (stdout/file)
//
// The hook is used for log rotation: before each write, we check if the log file
// needs to be rotated based on size, and rotate it if necessary.
//
// Thread Safety:
//   - hookMutex (RWMutex): Protects concurrent access to the writeHook function pointer
//     Multiple readers can check the hook concurrently, but updates are exclusive.
//   - rotationMutex (Mutex): Serializes the entire check-rotate sequence to prevent
//     race conditions where multiple goroutines try to rotate simultaneously.

var (
	writeHook     func([]byte) error // Hook function called on every write (currently: BeforeLogFileWrite)
	hookMutex     sync.RWMutex       // Protects concurrent access to writeHook function pointer
	rotationMutex sync.Mutex         // Serializes file rotation operations (check size + rotate)
)

// HookWriter wraps an io.Writer and calls a hook function on every write.
// This is the core interceptor that enables log rotation without modifying slog internals.
type HookWriter struct {
	writer io.Writer
}

// Write implements io.Writer interface.
// It calls the registered hook (if any) before writing to the underlying writer.
//
// Flow:
//  1. Read current hook (using RLock for concurrent safety)
//  2. If hook exists, call it with a defensive copy of the data
//  3. If hook returns error, abort write and return error (log is lost)
//  4. Otherwise, write to underlying writer (stdout + file)
//
// Thread Safety:
//  - Uses RLock/RUnlock to safely read the hook function pointer
//  - Makes a defensive copy of data to prevent race if hook modifies it
//
// Current Behavior Note:
//  If the hook returns an error, the write is aborted and the log entry is lost.
//  This ensures rotation failures are visible, but may need reconsideration
//  for production resilience (alternative: log hook errors but write anyway).
func (hw *HookWriter) Write(p []byte) (n int, err error) {
	// Read current hook safely (allows concurrent log writes)
	hookMutex.RLock()
	hook := writeHook
	hookMutex.RUnlock()

	if hook != nil {
		// Make a copy to avoid race conditions if hook modifies the slice
		hookCopy := make([]byte, len(p))
		copy(hookCopy, p)
		err = hook(hookCopy)
		if err != nil {
			return 0, errors.Wrap(err, "failed to call write hook")
		}
	}

	// Write to underlying writer (stdout + file)
	return hw.writer.Write(p)
}

// SetWriteHook registers a function to be called on every log write.
// The function receives a copy of the log data being written.
// Set to nil to remove the hook.
//
// Thread Safety:
//  Uses full Lock (not RLock) to ensure exclusive access when updating the hook.
//  This prevents torn reads where a writer might see a partially-updated pointer.
//
// Usage:
//  SetWriteHook(BeforeLogFileWrite)  // Enable rotation checks
//  SetWriteHook(nil)                 // Disable hook
func SetWriteHook(hook func([]byte) error) {
	hookMutex.Lock()
	defer hookMutex.Unlock()
	writeHook = hook
}
