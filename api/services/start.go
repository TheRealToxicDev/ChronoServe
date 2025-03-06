package services

import (
	"net/http"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// @Summary      Start a service
// @Description  Starts the specified system service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        name path string true "Service name"
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

	// Call the actual handler
	serviceHandler.StartService(w, r)
}
