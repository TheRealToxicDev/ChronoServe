package services

import (
	"net/http"
	"regexp"

	"github.com/toxic-development/sysmanix/utils"
)

type ServiceHandler interface {
	ListServices(w http.ResponseWriter, r *http.Request)
	StartService(w http.ResponseWriter, r *http.Request)
	StopService(w http.ResponseWriter, r *http.Request)
	ViewServiceLogs(w http.ResponseWriter, r *http.Request)
	GetServiceStatus(w http.ResponseWriter, r *http.Request)
}

type BaseServiceHandler struct{}

func (h *BaseServiceHandler) ValidateServiceName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_.]+$`, name)
	return matched
}

func (h *BaseServiceHandler) HandleError(w http.ResponseWriter, message string, statusCode int) {
	utils.WriteErrorResponse(w, message, statusCode)
}
