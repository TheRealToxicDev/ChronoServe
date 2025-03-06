package services

import (
	"net/http"
	"regexp"
	"time"
)

// ServiceInfo represents basic information about a service
type ServiceInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
}

// ServiceStatus represents detailed status of a service
type ServiceStatus struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	IsActive  bool      `json:"isActive"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LogEntry represents a service log entry
type LogEntry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type ServiceHandler interface {
	// Original HTTP handler methods
	ListServices(w http.ResponseWriter, r *http.Request)
	StartService(w http.ResponseWriter, r *http.Request)
	StopService(w http.ResponseWriter, r *http.Request)
	ViewServiceLogs(w http.ResponseWriter, r *http.Request)
	GetServiceStatus(w http.ResponseWriter, r *http.Request)

	// New data-returning methods
	GetServices() ([]ServiceInfo, error)
	GetServiceStatusByName(name string) (*ServiceStatus, error)
	StartServiceByName(name string) error
	StopServiceByName(name string) error
	GetServiceLogs(name string, lines int) ([]LogEntry, error)
}

type BaseServiceHandler struct{}

func (h *BaseServiceHandler) ValidateServiceName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_.]+$`, name)
	return matched
}
