package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

// IsAdmin checks if the current process has administrative privileges
func IsAdmin() bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("net", "session")
		err := cmd.Run()
		return err == nil
	}

	// For Linux/Unix, check if user is root (UID 0)
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return string(output) == "0\n"
}

// GetOperatingSystem returns the current OS in a standardized format
func GetOperatingSystem() string {
	return runtime.GOOS
}

// IsProcessRunning checks if a process with the given PID is running
func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		// On Windows, FindProcess always succeeds, so we need to check if the process is actually running
		err = process.Signal(syscall.Signal(0))
		return err == nil
	}

	// On Unix-like systems, FindProcess only returns a valid process if it exists
	return true
}

// GenerateUUID creates a new UUID v4 (random) and returns it as a string
func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
