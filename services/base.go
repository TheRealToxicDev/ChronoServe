package services

import (
	"errors"
	"net/http"
	"regexp"
	"time"
)

// ServiceInfo represents basic information about a service
type ServiceInfo struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	Status      string    `json:"status"`
	IsActive    bool      `json:"isActive"`
	UpdatedAt   time.Time `json:"updatedAt"`
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

// ServiceHandler is an interface for handling service-related HTTP requests
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

// Common error definitions for the services package
var (
	// ErrInvalidServiceName is returned when a service name fails validation
	ErrInvalidServiceName = errors.New("invalid service name")

	// ErrProtectedService is returned when attempting operations on protected services
	ErrProtectedService = errors.New("operation not allowed on protected system service")

	// ErrServiceNotFound is returned when a requested service doesn't exist
	ErrServiceNotFound = errors.New("service not found")

	// ErrInsufficientPermissions is returned when the user lacks required permissions
	ErrInsufficientPermissions = errors.New("insufficient permissions to perform this operation")
)
