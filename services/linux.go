package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/toxic-development/sysmanix/utils"
)

// Ensure SystemdService implements ServiceHandler
var _ ServiceHandler = (*SystemdService)(nil)

// SystemdService implements the ServiceHandler interface for Linux
type SystemdService struct {
	BaseServiceHandler
	cache      map[string]ServiceStatus
	cacheMutex sync.RWMutex
	cacheTTL   time.Duration
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsActive  bool      `json:"isActive"`
}

// NewSystemdService creates a new systemd service handler
func NewSystemdService() *SystemdService {
	return &SystemdService{
		cache:    make(map[string]ServiceStatus),
		cacheTTL: 5 * time.Minute,
	}
}

// ListServices lists all systemd services
func (s *SystemdService) ListServices(w http.ResponseWriter, r *http.Request) {
	if utils.GetOperatingSystem() != "linux" {
		s.HandleError(w, "Systemd is only supported on Linux", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--output=json")
	output, err := cmd.Output()
	if err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to list services: %v", err), http.StatusInternalServerError)
		return
	}

	var services []map[string]interface{}
	if err := json.Unmarshal(output, &services); err != nil {
		s.HandleError(w, "Failed to parse service data", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Services retrieved successfully", services)
}

// StartService starts a systemd service
func (s *SystemdService) StartService(w http.ResponseWriter, r *http.Request) {
	name := utils.ExtractServiceName(r.URL.Path)
	if !s.ValidateServiceName(name) {
		s.HandleError(w, "Invalid service name", http.StatusBadRequest)
		return
	}

	// Check if service is already running
	status, err := s.getServiceActiveState(name)
	if err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to check service status: %v", err), http.StatusInternalServerError)
		return
	}

	if status == "active" {
		utils.WriteSuccessResponse(w, fmt.Sprintf("Service %s is already running", name), nil)
		return
	}

	cmd := exec.Command("systemctl", "start", name)
	if err := cmd.Run(); err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to start service %s: %v", name, err), http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	s.cacheMutex.Lock()
	delete(s.cache, name)
	s.cacheMutex.Unlock()

	utils.WriteSuccessResponse(w, fmt.Sprintf("Service %s started successfully", name), nil)
}

// StopService stops a systemd service
func (s *SystemdService) StopService(w http.ResponseWriter, r *http.Request) {
	name := utils.ExtractServiceName(r.URL.Path)
	if !s.ValidateServiceName(name) {
		s.HandleError(w, "Invalid service name", http.StatusBadRequest)
		return
	}

	// Check if service is already stopped
	status, err := s.getServiceActiveState(name)
	if err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to check service status: %v", err), http.StatusInternalServerError)
		return
	}

	if status == "inactive" {
		utils.WriteSuccessResponse(w, fmt.Sprintf("Service %s is already stopped", name), nil)
		return
	}

	cmd := exec.Command("systemctl", "stop", name)
	if err := cmd.Run(); err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to stop service %s: %v", name, err), http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	s.cacheMutex.Lock()
	delete(s.cache, name)
	s.cacheMutex.Unlock()

	utils.WriteSuccessResponse(w, fmt.Sprintf("Service %s stopped successfully", name), nil)
}

// ViewServiceLogs retrieves systemd service logs
func (s *SystemdService) ViewServiceLogs(w http.ResponseWriter, r *http.Request) {
	name := utils.ExtractServiceName(r.URL.Path)
	if !s.ValidateServiceName(name) {
		s.HandleError(w, "Invalid service name", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("journalctl", "-u", name, "--no-pager", "-n", "100", "--output=json")
	output, err := cmd.Output()
	if err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to retrieve logs for service %s: %v", name, err), http.StatusInternalServerError)
		return
	}

	var logs interface{}
	if err := json.Unmarshal(output, &logs); err != nil {
		s.HandleError(w, "Failed to parse log data", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Service logs retrieved successfully", logs)
}

// GetServiceStatus gets the current status of a systemd service
func (s *SystemdService) GetServiceStatus(w http.ResponseWriter, r *http.Request) {
	name := utils.ExtractServiceName(r.URL.Path)
	if !s.ValidateServiceName(name) {
		s.HandleError(w, "Invalid service name", http.StatusBadRequest)
		return
	}

	// Check cache first
	s.cacheMutex.RLock()
	if status, ok := s.cache[name]; ok {
		if time.Since(status.UpdatedAt) < s.cacheTTL {
			s.cacheMutex.RUnlock()
			utils.WriteSuccessResponse(w, "Service status retrieved from cache", status)
			return
		}
	}
	s.cacheMutex.RUnlock()

	cmd := exec.Command("systemctl", "show", name, "--property=ActiveState,SubState,UnitFileState")
	output, err := cmd.Output()
	if err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to get status for service %s: %v", name, err), http.StatusInternalServerError)
		return
	}

	status := parseSystemdStatus(string(output))
	status.Name = name
	status.UpdatedAt = time.Now()

	// Update cache
	s.cacheMutex.Lock()
	s.cache[name] = status
	s.cacheMutex.Unlock()

	utils.WriteSuccessResponse(w, "Service status retrieved successfully", status)
}

// Helper function to get service active state
func (s *SystemdService) getServiceActiveState(name string) (string, error) {
	cmd := exec.Command("systemctl", "show", name, "--property=ActiveState")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	parts := strings.Split(string(output), "=")
	if len(parts) != 2 {
		return "", fmt.Errorf("unexpected output format")
	}
	return strings.TrimSpace(parts[1]), nil
}

// Helper function to parse systemd status output
func parseSystemdStatus(output string) ServiceStatus {
	lines := strings.Split(output, "\n")
	status := ServiceStatus{
		UpdatedAt: time.Now(),
	}

	for _, line := range lines {
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "ActiveState":
			status.Status = value
			status.IsActive = value == "active"
		}
	}

	return status
}
