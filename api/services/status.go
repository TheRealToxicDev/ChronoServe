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
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /services/status/{name} [get]
func GetServiceStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	// Call the service handler to get the service status
	status, err := serviceHandler.GetServiceStatusByName(name)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to get service status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to our API ServiceStatus type
	apiStatus := ServiceStatus{
		Name:      status.Name,
		Status:    status.Status,
		IsActive:  status.IsActive,
		UpdatedAt: status.UpdatedAt,
	}

	// Write the response using the standard format
	response := ServiceStatusResponse{
		Status:  "success",
		Message: "Service status retrieved successfully",
		Data:    apiStatus,
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}
