package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/toxic-development/sysmanix/utils"
)

// Ensure WindowsService implements ServiceHandler
var _ ServiceHandler = (*WindowsService)(nil)

// WindowsService implements the ServiceHandler interface for Windows
type WindowsService struct {
	BaseServiceHandler
	cache      map[string]ServiceStatus
	cacheMutex sync.RWMutex
	cacheTTL   time.Duration
}

// NewWindowsService creates a new Windows service handler
func NewWindowsService() *WindowsService {
	return &WindowsService{
		cache:    make(map[string]ServiceStatus),
		cacheTTL: 5 * time.Minute,
	}
}

// ---- DATA RETURNING METHODS (NEW) ----

// GetServices returns a list of all Windows services
func (s *WindowsService) GetServices() ([]ServiceInfo, error) {
	script := `
        Get-Service | Select-Object Name, DisplayName, Status | ConvertTo-Json -Compress
    `
	out, err := s.executePowershell(script)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %v", err)
	}

	// Debug the output
	outStr := out.String()
	if outStr == "" {
		return nil, fmt.Errorf("empty response from PowerShell")
	}

	// Trim any whitespace or BOM characters
	outStr = strings.TrimSpace(outStr)

	var rawServices []map[string]interface{}
	if err := json.Unmarshal([]byte(outStr), &rawServices); err != nil {
		return nil, fmt.Errorf("failed to parse service data (%s): %v", outStr, err)
	}

	// Convert to our service info format
	services := make([]ServiceInfo, 0, len(rawServices))
	for _, svc := range rawServices {
		name, _ := svc["Name"].(string)
		displayName, _ := svc["DisplayName"].(string)
		status, _ := svc["Status"].(string)

		services = append(services, ServiceInfo{
			Name:        name,
			DisplayName: displayName,
			Status:      status,
		})
	}

	return services, nil
}

// GetServiceStatusByName gets the current status of a Windows service by name
func (s *WindowsService) GetServiceStatusByName(name string) (*ServiceStatus, error) {
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

	// Improved script with diagnostic information
	script := fmt.Sprintf(`
        try {
            $service = Get-Service -Name "%s" -ErrorAction Stop
            $obj = @{
                "Name" = $service.Name
                "DisplayName" = $service.DisplayName
                "Status" = $service.Status.ToString()
            }
            $json = ConvertTo-Json -InputObject $obj -Compress
            Write-Output $json
        } catch {
            Write-Error "Error getting service: $_"
            exit 1
        }
    `, name)

	out, err := s.executePowershell(script)
	if err != nil {
		return nil, fmt.Errorf("failed to get status for service %s: %v", name, err)
	}

	// Debug the output
	outStr := out.String()
	if outStr == "" {
		return nil, fmt.Errorf("empty response from PowerShell for service %s", name)
	}

	// Trim any whitespace or BOM characters
	outStr = strings.TrimSpace(outStr)

	var rawStatus map[string]interface{}
	if err := json.Unmarshal([]byte(outStr), &rawStatus); err != nil {
		return nil, fmt.Errorf("failed to parse service status (%s): %v", outStr, err)
	}

	// Extract values with defensive coding
	nameStr, _ := rawStatus["Name"].(string)
	statusStr, _ := rawStatus["Status"].(string)

	// Use the actual value if available, otherwise fall back to the requested name
	if nameStr == "" {
		nameStr = name
	}

	isActive := statusStr == "Running"

	status := ServiceStatus{
		Name:      nameStr,
		Status:    statusStr,
		IsActive:  isActive,
		UpdatedAt: time.Now(),
	}

	// Update cache
	s.cacheMutex.Lock()
	s.cache[name] = status
	s.cacheMutex.Unlock()

	return &status, nil
}

// StartServiceByName starts a Windows service by name
func (s *WindowsService) StartServiceByName(name string) error {
	if !s.ValidateServiceName(name) {
		return fmt.Errorf("invalid service name")
	}

	script := fmt.Sprintf(`
        $service = Get-Service -Name "%s"
        if ($service.Status -eq "Running") {
            Write-Output "Service is already running"
            exit 0
        }
        Start-Service -Name "%s"
        $service.WaitForStatus("Running", "00:00:30")
        Write-Output "Service started successfully"
    `, name, name)

	_, err := s.executePowershell(script)
	if err != nil {
		return fmt.Errorf("failed to start service %s: %v", name, err)
	}

	// Invalidate cache
	s.cacheMutex.Lock()
	delete(s.cache, name)
	s.cacheMutex.Unlock()

	return nil
}

// StopServiceByName stops a Windows service by name
func (s *WindowsService) StopServiceByName(name string) error {
	if !s.ValidateServiceName(name) {
		return fmt.Errorf("invalid service name")
	}

	script := fmt.Sprintf(`
        $service = Get-Service -Name "%s"
        if ($service.Status -eq "Stopped") {
            Write-Output "Service is already stopped"
            exit 0
        }
        Stop-Service -Name "%s"
        $service.WaitForStatus("Stopped", "00:00:30")
        Write-Output "Service stopped successfully"
    `, name, name)

	_, err := s.executePowershell(script)
	if err != nil {
		return fmt.Errorf("failed to stop service %s: %v", name, err)
	}

	// Invalidate cache
	s.cacheMutex.Lock()
	delete(s.cache, name)
	s.cacheMutex.Unlock()

	return nil
}

// GetServiceLogs retrieves logs for a Windows service
func (s *WindowsService) GetServiceLogs(name string, lines int) ([]LogEntry, error) {
	if !s.ValidateServiceName(name) {
		return nil, fmt.Errorf("invalid service name")
	}

	if lines <= 0 {
		lines = 100 // Default to 100 lines
	}

	script := fmt.Sprintf(`
        $service = Get-WmiObject -Class Win32_Service -Filter "Name='%s'"
        if ($service -eq $null) {
            Write-Error "Service not found"
            exit 1
        }

        $logs = Get-WinEvent -FilterHashtable @{
            LogName = 'System'
            ID = @(7036, 7045, 7040)  # Common service-related event IDs
            StartTime = (Get-Date).AddDays(-7)  # Last 7 days
        } -MaxEvents %d -ErrorAction SilentlyContinue | 
        Where-Object { $_.Message -like "*$($service.DisplayName)*" -or $_.Message -like "*%s*" } |
        Select-Object @{Name='Time';Expression={$_.TimeCreated.ToString('yyyy-MM-dd HH:mm:ss')}}, 
                      @{Name='Level';Expression={$_.LevelDisplayName}},
                      @{Name='Message';Expression={$_.Message}} |
        ConvertTo-Json
        
        if ($logs -eq $null) {
            Write-Output "[]"  # Return empty array if no logs found
        } else {
            Write-Output $logs
        }
    `, name, lines, name)

	out, err := s.executePowershell(script)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs for service %s: %v", name, err)
	}

	// Handle both single log entry and array of logs
	var entry LogEntry
	if err := json.Unmarshal(out.Bytes(), &entry); err == nil {
		return []LogEntry{entry}, nil
	}

	var logs []LogEntry
	if err := json.Unmarshal(out.Bytes(), &logs); err != nil {
		return nil, fmt.Errorf("failed to parse log data: %v", err)
	}

	return logs, nil
}

// ---- HTTP HANDLER METHODS (UPDATED) ----

// ListServices lists all Windows services
func (s *WindowsService) ListServices(w http.ResponseWriter, r *http.Request) {
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

// StartService starts a Windows service
func (s *WindowsService) StartService(w http.ResponseWriter, r *http.Request) {
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

// StopService stops a Windows service
func (s *WindowsService) StopService(w http.ResponseWriter, r *http.Request) {
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

// ViewServiceLogs retrieves Windows service logs
func (s *WindowsService) ViewServiceLogs(w http.ResponseWriter, r *http.Request) {
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

// GetServiceStatus gets the current status of a Windows service
func (s *WindowsService) GetServiceStatus(w http.ResponseWriter, r *http.Request) {
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

// ---- HELPER METHODS ----

// executePowershell executes a PowerShell script and returns its output
func (s *WindowsService) executePowershell(script string) (*bytes.Buffer, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%v: %s", err, stderr.String())
	}

	return &out, nil
}

// HandleError handles error responses for the service
func (s *WindowsService) HandleError(w http.ResponseWriter, message string, statusCode int) {
	utils.WriteErrorResponse(w, message, statusCode)
}
