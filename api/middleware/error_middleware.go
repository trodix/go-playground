package middleware

import (
	"encoding/json"
	"net/http"
	"time"
)

// ErrorResponse represents the structure of the error response
type ErrorResponse struct {
	Status    string       `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// ErrorHandlingMiddleware is the middleware for handling errors
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response recorder to capture the status code
		rec := &ResponseRecorder{ResponseWriter: w, status: http.StatusText(http.StatusOK)}
		
		// Call the next handler
		next.ServeHTTP(rec, r)		
	})
}

// responseRecorder is a custom ResponseWriter to capture the status code and error message
type ResponseRecorder struct {
	http.ResponseWriter
	status   string
	errorMessage string
}

// CaptureError sets the error message for the response
func (rr *ResponseRecorder) CaptureError(code int, msg string) {
	rr.status = http.StatusText(code)
	rr.errorMessage = msg
	rr.ResponseWriter.WriteHeader(code)

	// If the status code indicates an error, send a structured JSON response
	errorResponse := ErrorResponse{
		Status:    http.StatusText(code),
		Message:   rr.errorMessage,
		Timestamp: time.Now(),
	}

	// Write the JSON response
	rr.ResponseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rr.ResponseWriter).Encode(errorResponse)
}
