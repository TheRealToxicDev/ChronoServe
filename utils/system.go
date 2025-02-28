package utils

import (
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
