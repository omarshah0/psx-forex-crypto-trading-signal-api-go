package utils

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Status  string      `json:"status"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Status string       `json:"status"`
	Type   string       `json:"type"`
	Error  ErrorDetails `json:"error"`
}

// ErrorDetails contains error information
type ErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SendSuccess sends a successful JSON response
func SendSuccess(w http.ResponseWriter, statusCode int, responseType string, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := SuccessResponse{
		Status:  "success",
		Type:    responseType,
		Data:    data,
		Message: message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// SendError sends an error JSON response
func SendError(w http.ResponseWriter, statusCode int, errorType string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Status: "error",
		Type:   errorType,
		Error: ErrorDetails{
			Code:    statusCode,
			Message: message,
		},
	}
	
	json.NewEncoder(w).Encode(response)
}

// Common error type constants
const (
	ErrorTypeBadRequest          = "bad_request"
	ErrorTypeUnauthorized        = "unauthorized"
	ErrorTypeForbidden           = "forbidden"
	ErrorTypeNotFound            = "not_found"
	ErrorTypeConflict            = "conflict"
	ErrorTypeValidation          = "validation_error"
	ErrorTypeInternalServer      = "internal_server_error"
	ErrorTypeServiceUnavailable  = "service_unavailable"
)

// Common response type constants
const (
	ResponseTypeResource   = "resource"
	ResponseTypeCollection = "collection"
	ResponseTypeAction     = "action"
	ResponseTypeAuth       = "auth"
)

