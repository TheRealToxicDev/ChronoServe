package services

import (
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

// NewSystemdService creates a new systemd service handler
func NewSystemdService() *SystemdService {
	return &SystemdService{
		cache:    make(map[string]ServiceStatus),
		cacheTTL: 5 * time.Minute,
	}
}

// ---- DATA RETURNING METHODS (NEW) ----

// GetServices returns a list of all systemd services
func (s *SystemdService) GetServices() ([]ServiceInfo, error) {
	if utils.GetOperatingSystem() != "linux" {
		return nil, fmt.Errorf("systemd is only supported on Linux")
	}

	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %v", err)
	}

	// Parse the output into service info
	services := []ServiceInfo{}
	lines := strings.Split(string(output), "\n")

	// Skip header and process each line
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Split by whitespace
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		serviceName := strings.TrimSuffix(fields[0], ".service")
		displayName := serviceName
		status := fields[3] // LOADED, ACTIVE, etc.

		services = append(services, ServiceInfo{
			Name:        serviceName,
			DisplayName: displayName,
			Status:      status,
		})
	}

	return services, nil
}

// GetServiceStatusByName gets the current status of a systemd service by name
func (s *SystemdService) GetServiceStatusByName(name string) (*ServiceStatus, error) {
	if !s.ValidateServiceName(name) {
		return nil, fmt.Errorf("invalid service name")
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

// StartServiceByName starts a systemd service by name
func (s *SystemdService) StartServiceByName(name string) error {
	if !s.ValidateServiceName(name) {
		return fmt.Errorf("invalid service name")
	}

	// Check if service is already running
	status, err := s.getServiceActiveState(name)
	if err != nil {
		return fmt.Errorf("failed to check service status: %v", err)
	}

	if status == "active" {
		return nil // Already running, considered success
	}

	cmd := exec.Command("systemctl", "start", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start service %s: %v", name, err)
	}

	// Invalidate cache
	s.cacheMutex.Lock()
	delete(s.cache, name)
	s.cacheMutex.Unlock()

	return nil
}

// StopServiceByName stops a systemd service by name
func (s *SystemdService) StopServiceByName(name string) error {
	if !s.ValidateServiceName(name) {
		return fmt.Errorf("invalid service name")
	}

	// Check if service is already stopped
	status, err := s.getServiceActiveState(name)
	if err != nil {
		return fmt.Errorf("failed to check service status: %v", err)
	}

	if status == "inactive" {
		return nil // Already stopped, considered success
	}

	cmd := exec.Command("systemctl", "stop", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop service %s: %v", name, err)
	}

	// Invalidate cache
	s.cacheMutex.Lock()
	delete(s.cache, name)
	s.cacheMutex.Unlock()

	return nil
}

// GetServiceLogs retrieves logs for a systemd service
func (s *SystemdService) GetServiceLogs(name string, lines int) ([]LogEntry, error) {
	if !s.ValidateServiceName(name) {
		return nil, fmt.Errorf("invalid service name")
	}

	if lines <= 0 {
		lines = 100 // Default to 100 lines
	}

	cmd := exec.Command("journalctl", "-u", name, "--no-pager", "-n", fmt.Sprintf("%d", lines), "--output=short")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs for service %s: %v", name, err)
	}

	// Parse the output into log entries
	logLines := strings.Split(string(output), "\n")
	entries := make([]LogEntry, 0, len(logLines))

	for _, line := range logLines {
		if line == "" {
			continue
		}

		// Basic log parsing - can be improved for systemd journal format
		parts := strings.SplitN(line, ":", 2)
		var timestamp, level, message string

		if len(parts) >= 2 {
			timestampParts := strings.Fields(parts[0])
			if len(timestampParts) > 0 {
				timestamp = strings.Join(timestampParts[0:3], " ")
				if len(timestampParts) > 3 {
					level = timestampParts[3]
				}
			}
			message = strings.TrimSpace(parts[1])
		} else {
			message = line
		}

		entries = append(entries, LogEntry{
			Time:    timestamp,
			Level:   level,
			Message: message,
		})
	}

	return entries, nil
}

// ---- HTTP HANDLER METHODS (UPDATED) ----

// ListServices lists all systemd services
func (s *SystemdService) ListServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	services, err := s.GetServices()
	if err != nil {
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Services retrieved successfully", services)
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
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, fmt.Sprintf("Service %s started successfully", name), nil)
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
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, fmt.Sprintf("Service %s stopped successfully", name), nil)
}

// ViewServiceLogs retrieves systemd service logs
func (s *SystemdService) ViewServiceLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := utils.ExtractServiceName(r.URL.Path)

	// Parse the number of lines from the query string
	lines := 100 // Default
	if linesStr := r.URL.Query().Get("lines"); linesStr != "" {
		if parsedLines, err := utils.ParseInt(linesStr); err == nil && parsedLines > 0 {
			lines = parsedLines
		}
	}

	logs, err := s.GetServiceLogs(name, lines)
	if err != nil {
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Service logs retrieved successfully", logs)
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
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Service status retrieved successfully", status)
}

// ---- HELPER METHODS (EXISTING) ----

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

// HandleError handles error responses for the service
func (s *SystemdService) HandleError(w http.ResponseWriter, message string, statusCode int) {
	utils.WriteErrorResponse(w, message, statusCode)
}
