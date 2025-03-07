package windows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// GetServiceStatusByName gets the current status of a Windows service by name
func (s *WindowsService) GetServiceStatusByName(name string) (*services.ServiceStatus, error) {
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

	status := services.ServiceStatus{
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

// GetServiceStatus gets the current status of a Windows service
func (s *WindowsService) GetServiceStatus(w http.ResponseWriter, r *http.Request) {
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
