package windows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// GetServices returns a list of all Windows services
func (s *WindowsService) GetServices() ([]services.ServiceInfo, error) {
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

	// Convert to our service info format and filter out protected services
	servicesList := make([]services.ServiceInfo, 0, len(rawServices))
	for _, svc := range rawServices {
		name, _ := svc["Name"].(string)

		// Skip protected services
		if s.IsProtectedService(name) {
			continue
		}

		displayName, _ := svc["DisplayName"].(string)
		status, _ := svc["Status"].(string)

		servicesList = append(servicesList, services.ServiceInfo{
			Name:        name,
			DisplayName: displayName,
			Status:      status,
		})
	}

	return servicesList, nil
}

// ListServices lists all Windows services
func (s *WindowsService) ListServices(w http.ResponseWriter, r *http.Request) {
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
