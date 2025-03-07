package linux

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// GetServiceLogs retrieves logs for a systemd service
func (s *SystemdService) GetServiceLogs(name string, lines int) ([]services.LogEntry, error) {
	if err := s.ValidateServiceOperation(name); err != nil {
		return nil, err
	}

	if lines <= 0 {
		lines = 100 // Default to 100 lines
	}

	cmd := exec.Command("journalctl", "-u", name, "--no-pager", "-n", fmt.Sprintf("%d", lines), "--output=short")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs for service %s: %v", name, err)
	}

	// Parse the output into log entries
	logLines := strings.Split(string(output), "\n")
	entries := make([]services.LogEntry, 0, len(logLines))

	for _, line := range logLines {
		if line == "" {
			continue
		}

		// Basic log parsing - can be improved for systemd journal format
		parts := strings.SplitN(line, ":", 2)
		var timestamp, level, message string

		if len(parts) >= 2 {
			timestampParts := strings.Fields(parts[0])
			if len(timestampParts) > 0 {
				timestamp = strings.Join(timestampParts[0:3], " ")
				if len(timestampParts) > 3 {
					level = timestampParts[3]
				}
			}
			message = strings.TrimSpace(parts[1])
		} else {
			message = line
		}

		entries = append(entries, services.LogEntry{
			Time:    timestamp,
			Level:   level,
			Message: message,
		})
	}

	return entries, nil
}

// ViewServiceLogs retrieves systemd service logs
func (s *SystemdService) ViewServiceLogs(w http.ResponseWriter, r *http.Request) {
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
