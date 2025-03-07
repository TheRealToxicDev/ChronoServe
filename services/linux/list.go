package linux

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// GetServices returns a list of all systemd services
func (s *SystemdService) GetServices() ([]services.ServiceInfo, error) {
	if utils.GetOperatingSystem() != "linux" {
		return nil, fmt.Errorf("systemd is only supported on Linux")
	}

	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %v", err)
	}

	// Parse the output into service info
	servicesList := []services.ServiceInfo{}
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

		// Skip protected services
		if s.IsProtectedService(serviceName) {
			continue
		}

		displayName := serviceName
		status := fields[3] // LOADED, ACTIVE, etc.

		servicesList = append(servicesList, services.ServiceInfo{
			Name:        serviceName,
			DisplayName: displayName,
			Status:      status,
		})
	}

	return servicesList, nil
}

// ListServices lists all systemd services
func (s *SystemdService) ListServices(w http.ResponseWriter, r *http.Request) {
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
