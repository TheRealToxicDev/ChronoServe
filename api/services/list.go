package services

import (
	"net/http"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// ServiceInfo represents information about a system service
type ServiceInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// @Summary      List all services
// @Description  Returns a list of all system services
// @Tags         services
// @Accept       json
// @Produce      json
// @Success      200  {array}  ServiceInfo
// @Security     BearerAuth
// @Router       /services [get]
func ListServices(w http.ResponseWriter, r *http.Request) {
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
	serviceHandler.ListServices(w, r)
}
