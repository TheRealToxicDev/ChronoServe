package utils

import (
	"encoding/json"
	"net/http"
)

// Standard response structures
type StandardResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// WriteJSONResponse writes a JSON response with the given data and status code
func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// WriteSuccessResponse writes a success response with the given message and data
func WriteSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	response := StandardResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	WriteJSONResponse(w, response, http.StatusOK)
}

// WriteErrorResponse writes an error response with the given message and status code
func WriteErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    statusCode,
	}

	WriteJSONResponse(w, response, statusCode)
}

// WriteMethodNotAllowedResponse writes a standard method not allowed response
func WriteMethodNotAllowedResponse(w http.ResponseWriter) {
	WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// WriteBadRequestResponse writes a standard bad request response
func WriteBadRequestResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Bad request"
	}
	WriteErrorResponse(w, message, http.StatusBadRequest)
}

// WriteUnauthorizedResponse writes a standard unauthorized response
func WriteUnauthorizedResponse(w http.ResponseWriter) {
	WriteErrorResponse(w, "Authentication required", http.StatusUnauthorized)
}

// WriteForbiddenResponse writes a standard forbidden response
func WriteForbiddenResponse(w http.ResponseWriter) {
	WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
}

// WriteNotFoundResponse writes a standard not found response
func WriteNotFoundResponse(w http.ResponseWriter, resource string) {
	message := "Not found"
	if resource != "" {
		message = resource + " not found"
	}
	WriteErrorResponse(w, message, http.StatusNotFound)
}

// WriteInternalErrorResponse writes a standard internal server error response
func WriteInternalErrorResponse(w http.ResponseWriter, err error) {
	WriteErrorResponse(w, "Internal server error", http.StatusInternalServerError)
}

// WriteTimeoutResponse writes a standard gateway timeout response
func WriteTimeoutResponse(w http.ResponseWriter, operation string) {
	message := "Request timed out"
	if operation != "" {
		message = "Request timed out during " + operation
	}
	WriteErrorResponse(w, message, http.StatusGatewayTimeout)
}

// WriteServiceProtectedResponse writes a response for protected service operations
func WriteServiceProtectedResponse(w http.ResponseWriter, serviceName string) {
	message := "Operation not allowed on protected system service"
	if serviceName != "" {
		message = "Operation not allowed on protected system service: " + serviceName
	}
	WriteErrorResponse(w, message, http.StatusForbidden)
}
