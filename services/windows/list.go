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

// GetServices returns a list of all Windows services with current status
func (s *WindowsService) GetServices() ([]services.ServiceInfo, error) {
	// Enhance PowerShell script to provide more status information
	script := `
		Get-Service | ForEach-Object {
			$status = $_.Status.ToString()
			$isActive = ($status -eq "Running")

			[PSCustomObject]@{
				Name = $_.Name
				DisplayName = $_.DisplayName
				Status = $status
				IsActive = $isActive
			}
		} | ConvertTo-Json -Compress
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

	// Convert to our service info format and filter out protected services
	servicesList := make([]services.ServiceInfo, 0, len(rawServices))
	for _, svc := range rawServices {
		name, _ := svc["Name"].(string)

		// Skip protected services
		if s.IsProtectedService(name) {
			continue
		}

		displayName, _ := svc["DisplayName"].(string)
		statusStr, _ := svc["Status"].(string)

		// Get isActive directly from PowerShell output
		isActive := false
		if activeVal, ok := svc["IsActive"].(bool); ok {
			isActive = activeVal
		} else {
			// Fallback to string comparison if boolean parsing fails
			isActive = strings.EqualFold(statusStr, "Running")
		}

		// Create service info with current status
		servicesList = append(servicesList, services.ServiceInfo{
			Name:        name,
			DisplayName: displayName,
			Status:      statusStr,
			IsActive:    isActive,
			UpdatedAt:   time.Now(),
		})
	}

	return servicesList, nil
}

// GetServicesWithDetails returns a list of all Windows services with detailed status
// This is a separate function for cases when you need extra details
func (s *WindowsService) GetServicesWithDetails() ([]services.ServiceInfo, error) {
	// First get the basic list of services with current status
	servicesList, err := s.GetServices()
	if err != nil {
		return nil, err
	}

	// For services where we need more detailed information
	return servicesList, nil
}

// ListServices lists all Windows services
func (s *WindowsService) ListServices(w http.ResponseWriter, r *http.Request) {
	// Force refresh of service list with each request
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	servicesList, err := s.GetServices()
	if err != nil {
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Services retrieved successfully", servicesList)
}
