package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"gosync/internal/api/apperror"
)

var validate = validator.New()

// DecodeAndValidate decodes a JSON request body into dst, then runs
// struct-tag validation on it. Returns a BadRequest apperror describing
// exactly what's wrong if either step fails.
func DecodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperror.BadRequest("invalid request body")
	}

	if err := validate.Struct(dst); err != nil {
		return apperror.BadRequest(formatValidationError(err))
	}

	return nil
}

func formatValidationError(err error) string {
	verrs, ok := err.(validator.ValidationErrors)
	if !ok || len(verrs) == 0 {
		return "validation failed"
	}
	// Report the first failing field, e.g. "Name is required"
	fe := verrs[0]
	return fe.Field() + " failed validation: " + fe.Tag()
}
