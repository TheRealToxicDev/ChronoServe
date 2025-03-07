package services

import (
	"fmt"
	"net/http"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/services/linux"
	"github.com/toxic-development/sysmanix/services/windows"
	"github.com/toxic-development/sysmanix/utils"
)

// @Summary      Stop a system service
// @Description  Stops the specified system service
// @Tags         services
// @Accept       json
// @Produce      json
// @Param        name  path  string  true  "Service name"
// @Success      200  {object}  ServiceActionResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /services/stop/{name} [post]
func StopService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	// Validate service name
	if !utils.ValidateServiceName(name) {
		utils.WriteErrorResponse(w, "Invalid service name", http.StatusBadRequest)
		return
	}

	// Call the service handler to stop the service
	err := serviceHandler.StopServiceByName(name)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to stop service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response using the standard format
	response := ServiceActionResponse{
		Status:  "success",
		Message: fmt.Sprintf("Service %s stopped successfully", name),
		Data:    nil,
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}
