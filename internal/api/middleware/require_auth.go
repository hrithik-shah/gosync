package middleware

import (
	"context"
	"net/http"
	"strings"

	"gosync/internal/apperror"
)

type ctxKey int

const userIDCtxKey ctxKey = 0

// RequireAuth checks for a valid Bearer token and rejects the request
// with a 401 if missing or invalid. On success, it stores the resolved
// user ID in the request context for downstream handlers.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			writeError(w, apperror.Unauthorized("missing authorization header"))
			return
		}

		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || token == "" {
			writeError(w, apperror.Unauthorized("invalid authorization header"))
			return
		}

		userID, err := validateToken(token)
		if err != nil {
			writeError(w, apperror.Unauthorized("invalid or expired token"))
			return
		}

		ctx := context.WithValue(r.Context(), userIDCtxKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext retrieves the authenticated user's ID, set by RequireAuth.
func UserIDFromContext(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(userIDCtxKey).(string)
	return id, ok
}

// validateToken is a placeholder — replace with real JWT parsing/verification
// or a session store lookup.
func validateToken(token string) (userID string, err error) {
	// TODO: parse/verify JWT, or look up session
	return "", apperror.Unauthorized("token validation not implemented")
}
