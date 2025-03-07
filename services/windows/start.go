package windows

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/toxic-development/sysmanix/utils"
)

// StartServiceByName starts a Windows service by name
func (s *WindowsService) StartServiceByName(name string) error {
	if err := s.ValidateServiceOperation(name); err != nil {
		return err
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

// StartService starts a Windows service
func (s *WindowsService) StartService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := utils.ExtractServiceName(r.URL.Path)

	err := s.StartServiceByName(name)
	if err != nil {
		// Return forbidden status for protected services instead of internal error
		if strings.Contains(err.Error(), "protected system service") {
			utils.WriteErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, fmt.Sprintf("Service %s started successfully", name), nil)
}
