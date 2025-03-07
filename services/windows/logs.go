package windows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// GetServiceLogs retrieves logs for a Windows service
func (s *WindowsService) GetServiceLogs(name string, lines int) ([]services.LogEntry, error) {
	if err := s.ValidateServiceOperation(name); err != nil {
		return nil, err
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
	var entry services.LogEntry
	if err := json.Unmarshal(out.Bytes(), &entry); err == nil {
		return []services.LogEntry{entry}, nil
	}

	var logs []services.LogEntry
	if err := json.Unmarshal(out.Bytes(), &logs); err != nil {
		return nil, fmt.Errorf("failed to parse log data: %v", err)
	}

	return logs, nil
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
		// Return forbidden status for protected services instead of internal error
		if strings.Contains(err.Error(), "protected system service") {
			utils.WriteErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, "Service logs retrieved successfully", logs)
}
