package services

import (
	"net/http"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// @Summary      Get service status
// @Description  Returns the current status of the specified service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        name  path  string  true  "Service name"
// @Success      200  {object}  ServiceStatusResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /services/status/{name} [get]
func GetServiceStatus(w http.ResponseWriter, r *http.Request) {
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

	serviceHandler.GetServiceStatus(w, r)
}
