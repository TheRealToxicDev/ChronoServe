package linux

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/toxic-development/sysmanix/utils"
)

// StopServiceByName stops a systemd service by name
func (s *SystemdService) StopServiceByName(name string) error {
	if err := s.ValidateServiceOperation(name); err != nil {
		return err
	}

	// Check if service is already stopped
	status, err := s.getServiceActiveState(name)
	if err != nil {
		return fmt.Errorf("failed to check service status: %v", err)
	}

	if status == "inactive" {
		return nil // Already stopped, considered success
	}

	// Stop the service asynchronously
	cmd := exec.Command("systemctl", "stop", name)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to stop service %s: %v", name, err)
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
				return fmt.Errorf("failed to check service status after stop: %v", err)
			}

			if currentStatus == "inactive" {
				return nil // Successfully stopped
			}

			if time.Since(startTime) > timeout {
				return fmt.Errorf("timeout waiting for service to stop (status: %s)", currentStatus)
			}
		}
	}
}

// StopService stops a systemd service
func (s *SystemdService) StopService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := utils.ExtractServiceName(r.URL.Path)
	err := s.StopServiceByName(name)
	if err != nil {
		// Return forbidden status for protected services instead of internal error
		if strings.Contains(err.Error(), "protected system service") {
			utils.WriteErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, fmt.Sprintf("Service %s stopped successfully", name), nil)
}
