package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gosom/kit/lib"
)

// ErrResponse is the response body for an error.
type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// JSON writes the given data to the response as JSON.
func JSON(w http.ResponseWriter, r *http.Request, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// I don't know what to do here. The headers have already been
		// written, so we can't write a status code.
		panic(err)
	}
}

// JSONError writes the given error to the response as JSON.
func JSONError(w http.ResponseWriter, r *http.Request, err error) {
	var resp ErrResponse
	var e lib.ApiError
	var ve validator.ValidationErrors
	switch {
	case errors.As(err, &e):
		resp.Code, resp.Message = e.ApiError()
	case errors.As(err, &ve):
		resp.Code = http.StatusBadRequest
		resp.Message = ve.Error()
	default:
		resp.Code = http.StatusInternalServerError
		resp.Message = http.StatusText(resp.Code)
	}
	// I don't like that. Is there any better way to pass the context
	// to the middleware?
	*r = *r.WithContext(lib.NewContextWithErr(r.Context(), err))
	JSON(w, r, resp.Code, resp)
}
