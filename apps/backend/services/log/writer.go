package log

import (
	"io"
	"sync"

	"github.com/cockroachdb/errors"
)

var (
	writeHook func([]byte) error // Hook function called on every write
	hookMutex sync.RWMutex       // Protects writeHook

)

// HookWriter wraps an io.Writer and calls a hook function on every write
type HookWriter struct {
	writer io.Writer
}

func (hw *HookWriter) Write(p []byte) (n int, err error) {
	// Call hook function if registered
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

	// Write to underlying writer
	return hw.writer.Write(p)
}

// SetWriteHook registers a function to be called on every log write.
// The function receives a copy of the log data being written.
// Set to nil to remove the hook.
func SetWriteHook(hook func([]byte) error) {
	hookMutex.Lock()
	defer hookMutex.Unlock()
	writeHook = hook
}
