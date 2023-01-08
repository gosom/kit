package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/xid"
)

// RequestIDProvider is a function that returns a request ID.
var RequestIDProvider = func() string {
	return xid.New().String()
}

// TimeProvider is a function that returns the current time.
var TimeProvider = func() time.Time {
	return time.Now().UTC()
}

type logResponseWriter struct {
	http.ResponseWriter
	status       int
	bytesWritten int
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if err != nil {
		return n, err
	}
	w.bytesWritten += n
	return n, err
}

func StringURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
