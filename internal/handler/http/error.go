package http

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
} 