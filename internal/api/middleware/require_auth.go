package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"gosync/internal/api/apperror"
	"gosync/internal/api/utils/jwtutil"
	"gosync/internal/config"
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
func UserIDFromContext(r *http.Request) (uuid.UUID, bool) {
	id, ok := r.Context().Value(userIDCtxKey).(uuid.UUID)
	return id, ok
}

// MustUserID retrieves the authenticated user's ID, panicking if called
// on a request that wasn't processed by RequireAuth.
func MustUserID(r *http.Request) uuid.UUID {
	id, ok := UserIDFromContext(r)
	if !ok {
		panic("MustUserID called on a request with no authenticated user")
	}
	return id
}

// validateToken is a placeholder — replace with real JWT parsing/verification
// or a session store lookup.
func validateToken(token string) (userID uuid.UUID, err error) {
	cfg := config.Get()
	return jwtutil.Decode(token, cfg.JWTSecret())
}
