package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendResponse sends a JSON response with the given status code
func SendResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := APIResponse{
		Success: success,
		Message: message,
		Data:    data,
		Error:   err,
	}
	json.NewEncoder(w).Encode(response)
}
