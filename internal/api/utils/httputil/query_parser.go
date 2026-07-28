package httputil

import (
	"gosync/internal/api/apperror"
	"gosync/internal/config"
	"net/http"
	"strconv"
)

func ParsePagination(r *http.Request) (cursor uint64, limit int, err error) {
	q := r.URL.Query()

	if v := q.Get("cursor"); v != "" {
		parsed, parseErr := strconv.ParseUint(v, 10, 64)
		if parseErr != nil {
			return 0, 0, apperror.BadRequest("invalid cursor")
		}
		cursor = parsed
	} else {
		return 0, 0, apperror.BadRequest("cursor is required")
	}

	limit = config.DefaultEventLimit
	if v := q.Get("limit"); v != "" {
		parsed, parseErr := strconv.Atoi(v)
		if parseErr != nil || parsed <= 0 {
			return 0, 0, apperror.BadRequest("invalid limit")
		}
		limit = parsed
	}
	if limit > config.MaxEventLimit {
		limit = config.MaxEventLimit
	}

	return cursor, limit, nil
}
