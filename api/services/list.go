package services

import (
	"net/http"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/utils"
)

// @Summary      List system services
// @Description  Returns a list of all available system services
// @Tags         services
// @Accept       json
// @Produce      json
// @Success      200  {object}  ServiceListResponse
// @Failure      500  {object}  utils.ErrorResponse
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

	// Get services from the service handler
	servicesList, err := serviceHandler.GetServices()
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to retrieve services: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to the expected response format
	summaries := make([]ServiceSummary, 0, len(servicesList))
	for _, svc := range servicesList {
		summaries = append(summaries, ServiceSummary{
			Name:        svc.Name,
			DisplayName: svc.DisplayName,
			Status:      svc.Status,
		})
	}

	// Write the response using the standard format
	response := ServiceListResponse{
		Status:  "success",
		Message: "Services retrieved successfully",
		Data:    summaries,
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}
