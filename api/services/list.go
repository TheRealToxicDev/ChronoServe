package services

import (
	"net/http"
	"time"

	"github.com/toxic-development/sysmanix/services"
	"github.com/toxic-development/sysmanix/services/linux"
	"github.com/toxic-development/sysmanix/services/windows"
	"github.com/toxic-development/sysmanix/utils"
)

// @Summary      List system services
// @Description  Returns a list of all available system services
// @Tags         services
// @Accept       json
// @Produce      json
// @Success      200  {object}  ServiceListResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /services [get]
func ListServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	// Always get fresh data with a timeout
	resultChan := make(chan []services.ServiceInfo, 1)
	errChan := make(chan error, 1)

	go func() {
		// Get fresh service data
		servicesList, err := serviceHandler.GetServices()
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- servicesList
	}()

	// Wait for result or timeout after 10 seconds
	select {
	case servicesList := <-resultChan:
		// Convert to the expected response format
		summaries := make([]ServiceSummary, 0, len(servicesList))
		for _, svc := range servicesList {
			summaries = append(summaries, ServiceSummary{
				Name:        svc.Name,
				DisplayName: svc.DisplayName,
				Status:      svc.Status,
				IsActive:    svc.IsActive,
				UpdatedAt:   svc.UpdatedAt,
			})
		}

		// Write the response using the standard format
		response := ServiceListResponse{
			Status:  "success",
			Message: "Services retrieved successfully",
			Data:    summaries,
		}

		utils.WriteJSONResponse(w, response, http.StatusOK)

	case err := <-errChan:
		utils.WriteErrorResponse(w, "Failed to retrieve services: "+err.Error(), http.StatusInternalServerError)

	case <-time.After(10 * time.Second):
		utils.WriteErrorResponse(w, "Request timed out while retrieving services", http.StatusGatewayTimeout)
	}
}
