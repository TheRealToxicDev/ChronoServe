package services

import (
	"fmt"
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
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /services/start/{name} [post]
func StartService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	// Call the service handler to start the service
	err := serviceHandler.StartServiceByName(name)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to start service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response using the standard format
	response := ServiceActionResponse{
		Status:  "success",
		Message: fmt.Sprintf("Service %s started successfully", name),
		Data:    nil,
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}
