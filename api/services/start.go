package services

import (
	"net/http"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// @Summary      Start a system service
// @Description  Starts the specified system service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        name  path  string  true  "Service name"
// @Success      200  {object}  ServiceActionResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /services/start/{name} [post]
func StartService(w http.ResponseWriter, r *http.Request) {
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

	serviceHandler.StartService(w, r)
}
