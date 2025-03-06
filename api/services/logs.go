package services

import (
	"net/http"

	"github.com/toxic-development/sysmanix/services"
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
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /services/logs/{name} [get]
func ViewServiceLogs(w http.ResponseWriter, r *http.Request) {
	var serviceHandler services.ServiceHandler

	// Choose the appropriate handler based on OS
	switch utils.GetOperatingSystem() {
	case "linux":
		serviceHandler = services.NewSystemdService()
	case "windows":
		serviceHandler = services.NewWindowsService()
	default:
		utils.WriteErrorResponse(w, "Unsupported operating system", http.StatusInternalServerError)
		return
	}

	serviceHandler.ViewServiceLogs(w, r)
}
