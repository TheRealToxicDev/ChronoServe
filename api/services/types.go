package services

import (
	"time"
)

// ServiceListResponse represents the response format for listing services
type ServiceListResponse struct {
	Status  string           `json:"status" example:"success"`
	Message string           `json:"message" example:"Services retrieved successfully"`
	Data    []ServiceSummary `json:"data"`
}

// ServiceSummary represents basic information about a service
type ServiceSummary struct {
	Name        string `json:"name" example:"nginx"`
	DisplayName string `json:"displayName" example:"Nginx Web Server"`
	Status      string `json:"status" example:"running"`
}

// ServiceStatusResponse represents the response format for service status
type ServiceStatusResponse struct {
	Status  string        `json:"status" example:"success"`
	Message string        `json:"message" example:"Service status retrieved successfully"`
	Data    ServiceStatus `json:"data"`
}

// ServiceStatus represents detailed information about a service's status
type ServiceStatus struct {
	Name      string    `json:"name" example:"nginx"`
	Status    string    `json:"status" example:"active"`
	IsActive  bool      `json:"isActive" example:"true"`
	UpdatedAt time.Time `json:"updatedAt" example:"2023-01-01T12:00:00Z"`
}

// LogEntry represents a single log entry from a service
type LogEntry struct {
	Time    string `json:"time" example:"2023-01-01T12:00:00Z"`
	Level   string `json:"level" example:"info"`
	Message string `json:"message" example:"Service started successfully"`
}

// ServiceLogsResponse represents the response format for service logs
type ServiceLogsResponse struct {
	Status  string     `json:"status" example:"success"`
	Message string     `json:"message" example:"Logs retrieved successfully"`
	Data    []LogEntry `json:"data"`
}

// ServiceActionResponse represents the response format for service actions (start/stop)
type ServiceActionResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Service started successfully"`
	Data    any    `json:"data"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"An error occurred while processing your request"`
	Code    int    `json:"code,omitempty" example:"404"`
}

// SuccessResponse represents a success response structure
type SuccessResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Operation completed successfully"`
	Data    any    `json:"data,omitempty"`
}
