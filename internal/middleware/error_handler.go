package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"gosync/internal/apperror"
)

// HandlerFunc is like http.HandlerFunc but can return an error.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ErrorHandler wraps a HandlerFunc, converting any returned error into
// the correct HTTP status code and JSON error body.
func ErrorHandler(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			writeError(w, err)
		}
	}
}

// writeError converts any error into the correct HTTP status code and
// JSON error body. Shared by ErrorHandler and any middleware (like
// RequireAuth) that needs to reject a request before a HandlerFunc runs.
func writeError(w http.ResponseWriter, err error) {
	appErr := apperror.Wrap(err)

	if appErr.Code >= 500 {
		log.Printf("internal error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": appErr.Message,
	})
}
