package linux

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// GetServiceStatusByName gets the current status of a systemd service by name
func (s *SystemdService) GetServiceStatusByName(name string) (*services.ServiceStatus, error) {
	if err := s.ValidateServiceOperation(name); err != nil {
		return nil, err
	}

	// Check cache first
	s.cacheMutex.RLock()
	if status, ok := s.cache[name]; ok {
		if time.Since(status.UpdatedAt) < s.cacheTTL {
			s.cacheMutex.RUnlock()
			return &status, nil
		}
	}
	s.cacheMutex.RUnlock()

	cmd := exec.Command("systemctl", "show", name, "--property=ActiveState,SubState,UnitFileState")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get status for service %s: %v", name, err)
	}

	status := parseSystemdStatus(string(output))
	status.Name = name
	status.UpdatedAt = time.Now()

	// Update cache
	s.cacheMutex.Lock()
	s.cache[name] = status
	s.cacheMutex.Unlock()

	return &status, nil
}

// GetServiceStatus gets the current status of a systemd service
func (s *SystemdService) GetServiceStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := utils.ExtractServiceName(r.URL.Path)
	status, err := s.GetServiceStatusByName(name)
	if err != nil {
		// Return forbidden status for protected services instead of internal error
		if strings.Contains(err.Error(), "protected system service") {
			utils.WriteErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Service status retrieved successfully", status)
}
