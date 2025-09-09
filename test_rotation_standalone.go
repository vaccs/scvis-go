package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Simple test program to verify log rotation logic
func main() {
	fmt.Println("Testing log rotation logic...")

	// Create a test log rotator with very small threshold
	logPath := "/tmp/test-rotation.log"
	maxSize := int64(1024) // 1KB threshold for easy testing
	maxBackups := 3

	rotator, err := NewLogRotator(logPath, maxSize, maxBackups)
	if err != nil {
		log.Fatal("Failed to create rotator:", err)
	}
	defer rotator.Close()

	fmt.Printf("Created rotator with %d byte threshold\n", maxSize)

	// Write some data to trigger rotation
	for i := 1; i <= 10; i++ {
		message := fmt.Sprintf("Test log entry #%d - %s\n", i, time.Now().Format(time.RFC3339))

		// Write multiple times to reach threshold
		for j := 0; j < 50; j++ {
			rotator.Write([]byte(message))
		}

		// Check if files exist
		fmt.Printf("\nAfter writing batch %d:\n", i)

		// List all related files
		pattern := logPath + "*"
		matches, _ := filepath.Glob(pattern)
		for _, match := range matches {
			info, _ := os.Stat(match)
			fmt.Printf("  %s: %d bytes\n", match, info.Size())
		}

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nTest complete!")
}

// Copy the LogRotator from main.go for testing
type LogRotator struct {
	logFile    *os.File
	logPath    string
	maxSize    int64
	maxBackups int
}

func NewLogRotator(logPath string, maxSize int64, maxBackups int) (*LogRotator, error) {
	lr := &LogRotator{
		logPath:    logPath,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}

	err := lr.openLogFile()
	if err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) openLogFile() error {
	dir := filepath.Dir(lr.logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	lr.logFile = file
	return nil
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
	if lr.needsRotation() {
		fmt.Println("  🔄 Rotation triggered!")
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	return lr.logFile.Write(p)
}

func (lr *LogRotator) needsRotation() bool {
	if lr.logFile == nil {
		return false
	}

	info, err := lr.logFile.Stat()
	if err != nil {
		return false
	}

	return info.Size() >= lr.maxSize
}

func (lr *LogRotator) rotate() error {
	if lr.logFile != nil {
		lr.logFile.Close()
	}

	for i := lr.maxBackups; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", lr.logPath, i)
		newPath := fmt.Sprintf("%s.%d", lr.logPath, i+1)

		if i == lr.maxBackups {
			os.Remove(oldPath)
		} else {
			if _, err := os.Stat(oldPath); err == nil {
				os.Rename(oldPath, newPath)
			}
		}
	}

	if _, err := os.Stat(lr.logPath); err == nil {
		os.Rename(lr.logPath, lr.logPath+".1")
	}

	return lr.openLogFile()
}

func (lr *LogRotator) Close() error {
	if lr.logFile != nil {
		return lr.logFile.Close()
	}
	return nil
}
