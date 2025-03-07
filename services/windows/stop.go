package windows

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/toxic-development/sysmanix/utils"
)

// StopServiceByName stops a Windows service by name
func (s *WindowsService) StopServiceByName(name string) error {
	if err := s.ValidateServiceOperation(name); err != nil {
		return err
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

// StopService stops a Windows service
func (s *WindowsService) StopService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := utils.ExtractServiceName(r.URL.Path)

	err := s.StopServiceByName(name)
	if err != nil {
		// Return forbidden status for protected services instead of internal error
		if strings.Contains(err.Error(), "protected system service") {
			utils.WriteErrorResponse(w, err.Error(), http.StatusForbidden)
			return
		}
		utils.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, fmt.Sprintf("Service %s stopped successfully", name), nil)
}
