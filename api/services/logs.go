package services

import (
	"net/http"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/services/linux"
	"github.com/toxic-development/sysmanix/services/windows"
	"github.com/toxic-development/sysmanix/utils"
)

// @Summary      Get service logs
// @Description  Returns logs for the specified service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        name   path    string  true   "Service name"
// @Param        lines  query   int     false  "Number of log lines to return (default: 100)"
// @Success      200  {object}  ServiceLogsResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /services/logs/{name} [get]
func ViewServiceLogs(w http.ResponseWriter, r *http.Request) {
	var serviceHandler services.ServiceHandler

	// Choose the appropriate handler based on OS
	switch utils.GetOperatingSystem() {
	case "linux":
		serviceHandler = linux.NewSystemdService()
	case "windows":
		serviceHandler = windows.NewWindowsService()
	default:
		utils.WriteErrorResponse(w, "Unsupported operating system", http.StatusInternalServerError)
		return
	}

	// Extract service name from path
	name := utils.ExtractServiceName(r.URL.Path)
	if name == "" {
		utils.WriteErrorResponse(w, "Missing or invalid service name", http.StatusBadRequest)
		return
	}

	// Validate service name - use the utility function directly
	if err := utils.ValidateServiceName(name); err != nil {
		utils.WriteErrorResponse(w, "Invalid service name: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the number of lines from the query string
	lines := 100 // Default
	if linesStr := r.URL.Query().Get("lines"); linesStr != "" {
		if parsedLines, err := utils.ParseInt(linesStr); err == nil && parsedLines > 0 {
			lines = parsedLines
		} else if err != nil {
			utils.WriteErrorResponse(w, "Invalid lines parameter: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Get logs from the service handler
	logs, err := serviceHandler.GetServiceLogs(name, lines)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to retrieve logs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert from services.LogEntry to our LogEntry type
	convertedLogs := make([]LogEntry, 0, len(logs))
	for _, log := range logs {
		convertedLogs = append(convertedLogs, LogEntry{
			Time:    log.Time,
			Level:   log.Level,
			Message: log.Message,
		})
	}

	// Write the response using our defined type
	response := ServiceLogsResponse{
		Status:  "success",
		Message: "Service logs retrieved successfully",
		Data:    convertedLogs,
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}
