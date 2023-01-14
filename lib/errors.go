package lib

import "strings"

// ApiError is an interface for all errors in the API.
type ApiError interface {
	// ApiError returns an HTTP status code and a message.
	ApiError() (int, string)
}

type StackTracer interface {
	StackTrace() string
}

func StackFormater(s string) []string {
	s = strings.Replace(s, "\t", "", -1)
	l := strings.Split(s, "\n")
	n := 0
	for i := range l {
		if l[i] != "" {
			l[n] = l[i]
			n++
		}
	}
	return l[:n]
}

// some common errors
var (
	ErrBadRequest       = apiError{code: 400, message: "bad request"}
	ErrAuth             = apiError{code: 401, message: "unauthorized"}
	ErrForbidden        = apiError{code: 403, message: "forbidden"}
	ErrNotFound         = apiError{code: 404, message: "not found"}
	ErrMethodNotAllowed = apiError{code: 405, message: "method not allowed"}
	ErrTimeout          = apiError{code: 408, message: "request timeout"}
	ErrConflict         = apiError{code: 409, message: "conflict"}
	ErrUnprocessable    = apiError{code: 422, message: "unprocessable"}
	ErrInternal         = apiError{code: 500, message: "internal error"}
)

// NewApiError creates a new ApiError with the given code and message.
// prefer using the predefined errors above.
func NewApiError(code int, message string) apiError {
	return apiError{code: code, message: message}
}

// WrapError wraps an error in an ApiError.
func WrapError(err error, ae apiError) error {
	return wrappedError{error: err, ae: ae}
}

type wrappedStackError struct {
	error
	Stack []byte
}

func (e wrappedStackError) StackTrace() string {
	return string(e.Stack)
}

// WrapErrorWithStack wraps an error and adds a stack trace.
func WrapErrorWithStack(err error, stack []byte) error {
	return wrappedStackError{
		error: err,
		Stack: stack,
	}
}

type wrappedError struct {
	error
	ae apiError
}

func (w wrappedError) Is(target error) bool {
	return w.ae == target
}

func (w wrappedError) ApiError() (int, string) {
	return w.ae.ApiError()
}

type apiError struct {
	code    int
	message string
}

func (o apiError) ApiError() (int, string) {
	return o.code, o.message
}

func (o apiError) Error() string {
	return o.message
}
