package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse defines the standard structure for API responses
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a successful JSON response
func SuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	sendResponse(w, http.StatusOK, true, message, data, nil)
}

// CreatedResponse sends a response for successful resource creation
func CreatedResponse(w http.ResponseWriter, message string, data interface{}) {
	sendResponse(w, http.StatusCreated, true, message, data, nil)
}

// ErrorResponse sends an error JSON response
func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	sendResponse(w, statusCode, false, message, nil, err)
}

// UnauthorizedResponse sends a 401 Unauthorized response
func UnauthorizedResponse(w http.ResponseWriter, message string) {
	sendResponse(w, http.StatusUnauthorized, false, message, nil, nil)
}

// Internal function to send JSON response
func sendResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	}

	if err != nil {
		response.Error = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}
