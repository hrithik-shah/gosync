package httputil

import (
	"gosync/internal/api/apperror"
	"net/http"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
)

// Util function ParseUUIDParam parses a chi URL param as a UUID, returning
// a BadRequest apperror if it's not valid.
func ParseUUIDParam(r *http.Request, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, name))
	if err != nil {
		return uuid.Nil, apperror.BadRequest("invalid " + name)
	}
	return id, nil
}
