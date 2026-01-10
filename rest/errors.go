package rest

import (
	"encoding/json"
	"net/http"
)

// ErrorCode represents standard error codes
type ErrorCode string

const (
	ErrorCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrorCodeInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrorCodeInternalError ErrorCode = "INTERNAL_ERROR"
)

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// WriteErrorResponse writes a standardized error response
func WriteErrorResponse(w http.ResponseWriter, code ErrorCode, message, details string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	
	json.NewEncoder(w).Encode(response)
}

// WriteNotFoundError writes a 404 error response
func WriteNotFoundError(w http.ResponseWriter, resource, id string) {
	WriteErrorResponse(w, ErrorCodeNotFound, 
		resource+" not found", 
		"No "+resource+" exists with ID "+id, 
		http.StatusNotFound)
}

// WriteInvalidInputError writes a 400 error response
func WriteInvalidInputError(w http.ResponseWriter, message string) {
	WriteErrorResponse(w, ErrorCodeInvalidInput, 
		"Invalid input", 
		message, 
		http.StatusBadRequest)
}

// WriteInternalError writes a 500 error response
func WriteInternalError(w http.ResponseWriter, err error) {
	WriteErrorResponse(w, ErrorCodeInternalError, 
		"Internal server error", 
		err.Error(), 
		http.StatusInternalServerError)
}