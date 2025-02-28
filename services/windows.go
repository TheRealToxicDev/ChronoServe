package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/therealtoxicdev/chronoserve/utils"
)

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

// ListServices lists all Windows services
func (s *WindowsService) ListServices(w http.ResponseWriter, r *http.Request) {
	script := `
        Get-Service | Select-Object Name, DisplayName, Status | ConvertTo-Json
    `
	out, err := s.executePowershell(script)
	if err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to list services: %v", err), http.StatusInternalServerError)
		return
	}

	var services []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &services); err != nil {
		s.HandleError(w, "Failed to parse service data", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Services retrieved successfully", services)
}

// StartService starts a Windows service
func (s *WindowsService) StartService(w http.ResponseWriter, r *http.Request) {
	name := utils.ExtractServiceName(r.URL.Path)
	if !s.ValidateServiceName(name) {
		s.HandleError(w, "Invalid service name", http.StatusBadRequest)
		return
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

	out, err := s.executePowershell(script)
	if err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to start service %s: %v", name, err), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, out.String(), nil)
}

// StopService stops a Windows service
func (s *WindowsService) StopService(w http.ResponseWriter, r *http.Request) {
	name := utils.ExtractServiceName(r.URL.Path)
	if !s.ValidateServiceName(name) {
		s.HandleError(w, "Invalid service name", http.StatusBadRequest)
		return
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

	out, err := s.executePowershell(script)
	if err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to stop service %s: %v", name, err), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, out.String(), nil)
}

// ViewServiceLogs retrieves Windows service logs
func (s *WindowsService) ViewServiceLogs(w http.ResponseWriter, r *http.Request) {
	name := utils.ExtractServiceName(r.URL.Path)
	if !s.ValidateServiceName(name) {
		s.HandleError(w, "Invalid service name", http.StatusBadRequest)
		return
	}

	script := fmt.Sprintf(`
        Get-WinEvent -FilterHashtable @{
            LogName = 'Application'
            ProviderName = '%s'
        } -MaxEvents 100 | Select-Object TimeCreated, LevelDisplayName, Message | ConvertTo-Json
    `, name)

	out, err := s.executePowershell(script)
	if err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to retrieve logs for service %s: %v", name, err), http.StatusInternalServerError)
		return
	}

	var logs interface{}
	if err := json.Unmarshal(out.Bytes(), &logs); err != nil {
		s.HandleError(w, "Failed to parse log data", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Logs retrieved successfully", logs)
}

// GetServiceStatus gets the current status of a Windows service
func (s *WindowsService) GetServiceStatus(w http.ResponseWriter, r *http.Request) {
	name := utils.ExtractServiceName(r.URL.Path)
	if !s.ValidateServiceName(name) {
		s.HandleError(w, "Invalid service name", http.StatusBadRequest)
		return
	}

	script := fmt.Sprintf(`
        Get-Service -Name "%s" | Select-Object Name, DisplayName, Status | ConvertTo-Json
    `, name)

	out, err := s.executePowershell(script)
	if err != nil {
		s.HandleError(w, fmt.Sprintf("Failed to get status for service %s: %v", name, err), http.StatusInternalServerError)
		return
	}

	var status interface{}
	if err := json.Unmarshal(out.Bytes(), &status); err != nil {
		s.HandleError(w, "Failed to parse service status", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Service status retrieved successfully", status)
}

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
