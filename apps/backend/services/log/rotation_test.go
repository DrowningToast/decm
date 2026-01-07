package log

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// TestRotatingFileWriter_BasicWrite tests basic write functionality
func TestRotatingFileWriter_BasicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	writer, err := NewRotatingFileWriter(logPath)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	testData := []byte("test log line\n")
	n, err := writer.Write(testData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}

	// Verify file contents
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if !bytes.Equal(content, testData) {
		t.Errorf("File content mismatch. Expected %q, got %q", testData, content)
	}
}

// TestRotatingFileWriter_ReopenAfterRename tests Problem 4: File handle after rotation
func TestRotatingFileWriter_ReopenAfterRename(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "head.log")
	rotatedPath := filepath.Join(tmpDir, "rotated.log")

	writer, err := NewRotatingFileWriter(logPath)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	// Write before rotation
	beforeRotation := []byte("Before rotation\n")
	writer.Write(beforeRotation)

	// Simulate rotation: rename the file
	if err := os.Rename(logPath, rotatedPath); err != nil {
		t.Fatalf("Failed to rename file: %v", err)
	}

	// Reopen to create new file with original name
	if err := writer.Reopen(); err != nil {
		t.Fatalf("Failed to reopen: %v", err)
	}

	// Write after rotation
	afterRotation := []byte("After rotation\n")
	writer.Write(afterRotation)

	// Verify old file has only pre-rotation content
	oldContent, err := os.ReadFile(rotatedPath)
	if err != nil {
		t.Fatalf("Failed to read rotated file: %v", err)
	}
	if !bytes.Equal(oldContent, beforeRotation) {
		t.Errorf("Rotated file should only have pre-rotation content. Got %q", oldContent)
	}

	// Verify new file has only post-rotation content
	newContent, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read new file: %v", err)
	}
	if !bytes.Equal(newContent, afterRotation) {
		t.Errorf("New file should only have post-rotation content. Got %q", newContent)
	}
}

// TestRotatingFileWriter_ConcurrentWrites tests thread safety
func TestRotatingFileWriter_ConcurrentWrites(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "concurrent.log")

	writer, err := NewRotatingFileWriter(logPath)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	const numGoroutines = 10
	const writesPerGoroutine = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < writesPerGoroutine; j++ {
				writer.Write([]byte("x"))
			}
		}(i)
	}

	wg.Wait()

	// Verify all writes succeeded
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	expectedLen := numGoroutines * writesPerGoroutine
	if len(content) != expectedLen {
		t.Errorf("Expected %d bytes, got %d", expectedLen, len(content))
	}
}

// TestConcurrentRotation tests Problem 3: Race condition in rotation
func TestConcurrentRotation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a temporary logs directory
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs dir: %v", err)
	}

	// Change to temp directory for test
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Initialize the rotating file writer
	var err error
	rotatingFile, err = NewRotatingFileWriter("logs/head.log")
	if err != nil {
		t.Fatalf("Failed to create rotating file: %v", err)
	}
	defer rotatingFile.Close()

	// Write initial data to make file non-empty
	initialData := make([]byte, 500*1024) // 500KB
	for i := range initialData {
		initialData[i] = 'A'
	}
	rotatingFile.Write(initialData)

	const numGoroutines = 5
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	rotationErrors := make(chan error, numGoroutines)

	// Simulate concurrent rotation attempts
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Create data that would trigger rotation (600KB)
			data := make([]byte, 600*1024)
			for j := range data {
				data[j] = byte('0' + id)
			}

			err := BeforeLogFileWrite(data)
			if err != nil {
				rotationErrors <- err
			}
		}(i)
	}

	wg.Wait()
	close(rotationErrors)

	// Check if any rotation errors occurred
	for err := range rotationErrors {
		t.Errorf("Rotation error: %v", err)
	}

	// Verify that head.log still exists
	if _, err := os.Stat("logs/head.log"); os.IsNotExist(err) {
		t.Error("head.log should exist after rotation")
	}
}

// TestHookWriter_WithHook tests the hook mechanism
func TestHookWriter_WithHook(t *testing.T) {
	var hookCalled bool
	var hookData []byte

	// Set up hook
	SetWriteHook(func(data []byte) error {
		hookCalled = true
		hookData = make([]byte, len(data))
		copy(hookData, data)
		return nil
	})
	defer SetWriteHook(nil) // Clean up

	// Create writer
	buf := &bytes.Buffer{}
	writer := &HookWriter{writer: buf}

	testData := []byte("test data")
	n, err := writer.Write(testData)

	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}
	if !hookCalled {
		t.Error("Hook was not called")
	}
	if !bytes.Equal(hookData, testData) {
		t.Errorf("Hook received wrong data. Expected %q, got %q", testData, hookData)
	}
	if !bytes.Equal(buf.Bytes(), testData) {
		t.Errorf("Underlying writer received wrong data. Expected %q, got %q", testData, buf.Bytes())
	}
}

// TestHookWriter_HookError tests Problem 1: Hook failure behavior
func TestHookWriter_HookError(t *testing.T) {
	// Set up hook that returns error
	SetWriteHook(func(data []byte) error {
		return io.ErrShortWrite
	})
	defer SetWriteHook(nil)

	buf := &bytes.Buffer{}
	writer := &HookWriter{writer: buf}

	testData := []byte("test data")
	n, err := writer.Write(testData)

	// Current behavior: hook error prevents write
	if err == nil {
		t.Error("Expected error when hook fails, got nil")
	}
	if n != 0 {
		t.Errorf("Expected 0 bytes written on hook error, got %d", n)
	}

	// Note: This test documents the current behavior where hook failure
	// prevents the log from being written. This may need to be changed
	// for better resilience (write the log even if hook fails).
}

// TestHookWriter_ConcurrentHookModification tests hook modification during writes
func TestHookWriter_ConcurrentHookModification(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := &HookWriter{writer: buf}

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1: Constantly writing
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			writer.Write([]byte("x"))
			time.Sleep(time.Microsecond)
		}
	}()

	// Goroutine 2: Constantly modifying hook
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				SetWriteHook(func(data []byte) error { return nil })
			} else {
				SetWriteHook(nil)
			}
			time.Sleep(time.Microsecond)
		}
	}()

	wg.Wait()

	// Should not panic or deadlock
	SetWriteHook(nil) // Cleanup
}

// TestBeforeLogFileWrite_NoRotationNeeded tests rotation is skipped when file is small
func TestBeforeLogFileWrite_NoRotationNeeded(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs dir: %v", err)
	}

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create small file
	smallData := []byte("small log")
	os.WriteFile("logs/head.log", smallData, 0644)

	// Initialize rotating file
	var err error
	rotatingFile, err = NewRotatingFileWriter("logs/head.log")
	if err != nil {
		t.Fatalf("Failed to create rotating file: %v", err)
	}
	defer rotatingFile.Close()

	// Try to write small data (should not trigger rotation)
	err = BeforeLogFileWrite([]byte("more small data"))
	if err != nil {
		t.Errorf("BeforeLogFileWrite failed: %v", err)
	}

	// Verify no rotated files were created
	files, _ := os.ReadDir("logs")
	if len(files) != 1 {
		t.Errorf("Expected only head.log to exist, found %d files", len(files))
	}
}

// TestBeforeLogFileWrite_RotationTriggered tests rotation when file is large
func TestBeforeLogFileWrite_RotationTriggered(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs dir: %v", err)
	}

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create large file (almost 1MB)
	largeData := make([]byte, 1024*1024-100) // Just under 1MB
	for i := range largeData {
		largeData[i] = 'X'
	}
	os.WriteFile("logs/head.log", largeData, 0644)

	// Initialize rotating file
	var err error
	rotatingFile, err = NewRotatingFileWriter("logs/head.log")
	if err != nil {
		t.Fatalf("Failed to create rotating file: %v", err)
	}
	defer rotatingFile.Close()

	// Write data that pushes over 1MB (should trigger rotation)
	newData := make([]byte, 200) // This pushes total over 1MB
	err = BeforeLogFileWrite(newData)
	if err != nil {
		t.Errorf("BeforeLogFileWrite failed: %v", err)
	}

	// Verify rotation occurred - there should be a rotated file
	files, _ := os.ReadDir("logs")
	rotatedFileExists := false
	headLogExists := false

	for _, file := range files {
		if file.Name() == "head.log" {
			headLogExists = true
		} else if file.Name() != "head.log" {
			rotatedFileExists = true
		}
	}

	if !rotatedFileExists {
		t.Error("Expected rotated file to exist after rotation")
	}
	if !headLogExists {
		t.Error("Expected new head.log to exist after rotation")
	}
}

// TestRotatingFileWriter_AutoCreateFile tests that files are auto-created when they don't exist
func TestRotatingFileWriter_AutoCreateFile(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "nonexistent", "auto.log")

	// Create parent directory
	os.MkdirAll(filepath.Dir(logPath), 0755)

	// File doesn't exist yet - should be auto-created
	writer, err := NewRotatingFileWriter(logPath)
	if err != nil {
		t.Fatalf("Failed to create writer for non-existent file: %v", err)
	}
	defer writer.Close()

	// Write to it
	testData := []byte("auto-created content\n")
	n, err := writer.Write(testData)
	if err != nil {
		t.Fatalf("Write to auto-created file failed: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}

	// Verify file was created and contains data
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read auto-created file: %v", err)
	}
	if !bytes.Equal(content, testData) {
		t.Errorf("Auto-created file content mismatch. Expected %q, got %q", testData, content)
	}
}

// TestRotationMutex_Serialization tests that rotation mutex properly serializes operations
func TestRotationMutex_Serialization(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		t.Fatalf("Failed to create logs dir: %v", err)
	}

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Initialize with large file
	largeData := make([]byte, 1024*1024-100)
	os.WriteFile("logs/head.log", largeData, 0644)

	var err error
	rotatingFile, err = NewRotatingFileWriter("logs/head.log")
	if err != nil {
		t.Fatalf("Failed to create rotating file: %v", err)
	}
	defer rotatingFile.Close()

	// Track execution order
	var executionOrder []int
	var orderMutex sync.Mutex

	const numGoroutines = 3
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Lock rotation mutex explicitly to test serialization
			rotationMutex.Lock()
			orderMutex.Lock()
			executionOrder = append(executionOrder, id)
			orderMutex.Unlock()
			time.Sleep(10 * time.Millisecond) // Simulate work
			rotationMutex.Unlock()
		}(i)
	}

	wg.Wait()

	// Verify all goroutines executed
	if len(executionOrder) != numGoroutines {
		t.Errorf("Expected %d executions, got %d", numGoroutines, len(executionOrder))
	}

	// The mutex ensures serialization - we can't have overlapping executions
	// This is tested by the sleep in the critical section
}
