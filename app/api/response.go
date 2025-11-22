package api

import (
	"encoding/json"
	"net/http"
)

// OKResponse sends a 200 OK KSON response
func OKResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, send a 500 Internal Server Error
		ErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

// ErrorResponse sends a JSON response for a given HTTP status code
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := map[string]string{
		"error": message,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
