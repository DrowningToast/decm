package log

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// BenchmarkRotatingFileWriter_ConcurrentWrites measures write performance under concurrency
func BenchmarkRotatingFileWriter_ConcurrentWrites(b *testing.B) {
	tmpDir := b.TempDir()
	logPath := filepath.Join(tmpDir, "bench.log")

	writer, err := NewRotatingFileWriter(logPath)
	if err != nil {
		b.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	data := []byte("benchmark log entry\n")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			writer.Write(data)
		}
	})
}

// BenchmarkRotatingFileWriter_SequentialWrites measures sequential write performance
func BenchmarkRotatingFileWriter_SequentialWrites(b *testing.B) {
	tmpDir := b.TempDir()
	logPath := filepath.Join(tmpDir, "bench.log")

	writer, err := NewRotatingFileWriter(logPath)
	if err != nil {
		b.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	data := []byte("benchmark log entry\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer.Write(data)
	}
}

// BenchmarkRotatingFileWriter_WithRotation measures performance including rotation
func BenchmarkRotatingFileWriter_WithRotation(b *testing.B) {
	tmpDir := b.TempDir()
	logPath := filepath.Join(tmpDir, "bench.log")

	writer, err := NewRotatingFileWriter(logPath)
	if err != nil {
		b.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	data := []byte("benchmark log entry\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer.Write(data)
		// Simulate rotation every 1000 writes
		if i%1000 == 999 {
			os.Rename(logPath, logPath+".rotated")
			writer.Reopen()
		}
	}
}

// BenchmarkHookWriter_WithHook measures hook overhead
func BenchmarkHookWriter_WithHook(b *testing.B) {
	tmpFile, _ := os.CreateTemp("", "bench")
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	writer := &HookWriter{writer: tmpFile}

	// Set up a simple hook
	SetWriteHook(func(data []byte) error {
		return nil // No-op hook
	})
	defer SetWriteHook(nil)

	data := []byte("benchmark log entry\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer.Write(data)
	}
}

// BenchmarkConcurrentRotation measures rotation under concurrent load
func BenchmarkConcurrentRotation(b *testing.B) {
	tmpDir := b.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	os.MkdirAll(logsDir, 0755)

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	var err error
	rotatingFile, err = NewRotatingFileWriter("logs/head.log")
	if err != nil {
		b.Fatalf("Failed to create rotating file: %v", err)
	}
	defer rotatingFile.Close()

	// Pre-fill with data
	initialData := make([]byte, 900*1024) // 900KB
	rotatingFile.Write(initialData)

	b.ResetTimer()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := make([]byte, 200*1024) // 200KB - triggers rotation
			BeforeLogFileWrite(data)
		}()
	}
	wg.Wait()
}
