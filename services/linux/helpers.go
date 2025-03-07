package linux

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

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
func parseSystemdStatus(output string) services.ServiceStatus {
	lines := strings.Split(output, "\n")
	status := services.ServiceStatus{
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
