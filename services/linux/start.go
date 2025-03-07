package linux

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/toxic-development/sysmanix/utils"
)

// StartServiceByName starts a systemd service by name
func (s *SystemdService) StartServiceByName(name string) error {
	if err := s.ValidateServiceOperation(name); err != nil {
		return err
	}

	// Check if service is already running
	status, err := s.getServiceActiveState(name)
	if err != nil {
		return fmt.Errorf("failed to check service status: %v", err)
	}

	if status == "active" {
		return nil // Already running, considered success
	}

	// Start the service asynchronously (don't wait in the command execution)
	cmd := exec.Command("systemctl", "start", name)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start service %s: %v", name, err)
	}

	// Allow the process to run in the background
	go func() {
		_ = cmd.Wait() // Prevent zombie process
	}()

	// Invalidate cache immediately
	s.cacheMutex.Lock()
	delete(s.cache, name)
	s.cacheMutex.Unlock()

	// Poll for service status with a reasonable timeout
	startTime := time.Now()
	timeout := 10 * time.Second
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			currentStatus, err := s.getServiceActiveState(name)
			if err != nil {
				return fmt.Errorf("failed to check service status after start: %v", err)
			}

			if currentStatus == "active" {
				return nil // Successfully started
			}

			if time.Since(startTime) > timeout {
				return fmt.Errorf("timeout waiting for service to start (status: %s)", currentStatus)
			}
		}
	}
}

// StartService starts a systemd service
func (s *SystemdService) StartService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := utils.ExtractServiceName(r.URL.Path)
	err := s.StartServiceByName(name)
	if err != nil {
		// Return forbidden status for protected services instead of internal error
		if strings.Contains(err.Error(), "protected system service") {
			utils.WriteErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, fmt.Sprintf("Service %s started successfully", name), nil)
}
