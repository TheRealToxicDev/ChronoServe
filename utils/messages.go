package utils

import (
	"encoding/json"
	"net/http"
)

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
	response := SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	WriteJSON(w, response, http.StatusOK)
}

// WriteErrorResponse writes a JSON error response
func WriteErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    statusCode,
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

// WriteJSONResponse writes a raw JSON response to the client
// Useful when you need to send a completely custom structure
func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	WriteJSON(w, data, statusCode)
}
