package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response format
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// WriteJSON writes a JSON response with proper headers
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If JSON encoding fails, send a plain text error
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode JSON response"))
	}
}

// WriteSuccessResponse writes a successful JSON response
func WriteSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	WriteJSON(w, response, http.StatusOK)
}

// WriteErrorResponse writes a JSON error response
func WriteErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := Response{
		Success: false,
		Error:   message,
	}
	WriteJSON(w, response, statusCode)
}

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, message, http.StatusBadRequest)
}

// WriteInternalError writes an internal server error response
func WriteInternalError(w http.ResponseWriter, err error) {
	message := "An internal server error occurred"
	if err != nil {
		// Only include error details in development
		message = err.Error()
	}
	WriteErrorResponse(w, message, http.StatusInternalServerError)
}
