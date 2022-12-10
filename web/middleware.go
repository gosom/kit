package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gosom/kit/core"
	"github.com/gosom/kit/logging"
	"github.com/realclientip/realclientip-go"
	"github.com/rs/cors"
)

// CorsConfig is the configuration for cors
type CorsConfig struct {
	NotAllowAll bool
	Options     cors.Options
}

// NewCors returns a cors middleware
func NewCors(opts CorsConfig) *cors.Cors {
	switch opts.NotAllowAll {
	case false:
		return cors.AllowAll()
	default:
		return cors.New(opts.Options)
	}
}

// Timeout is a middleware that sets a timeout on the request context.
func RequestLogger(report core.ErrorReporter) func(http.Handler) http.Handler {
	log := logging.Get()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := TimeProvider()
			reqID := RequestIDProvider()
			ctxLogger := log.With("request_id", reqID)
			r = r.WithContext(core.NewContextWithRequestID(r.Context(), reqID))
			r = r.WithContext(logging.NewContext(r.Context(), ctxLogger))
			lrw := &logResponseWriter{ResponseWriter: w, status: http.StatusOK}
			defer func() {
				level := logging.INFO
				switch {
				case lrw.status >= 400 && lrw.status < 500:
					level = logging.WARN
				case lrw.status >= 500 || lrw.status < 200:
					level = logging.ERROR
				}
				fields := []any{
					"status", lrw.status,
					"bytes", lrw.bytesWritten,
					"method", r.Method,
					"path", r.URL.Path,
					"query", r.URL.RawQuery,
					"ip", core.IPFromContext(r.Context()),
					"user-agent", r.UserAgent(),
					"latency", TimeProvider().Sub(start),
				}
				if lrw.status >= http.StatusInternalServerError {
					err := core.ErrorFromContext(r.Context())
					if err != nil {
						fields = append(fields, "error", err)
					}
					report.ReportError(r.Context(), r, err)
				}
				ctxLogger.Log(level, http.StatusText(lrw.status), fields...)
			}()
			next.ServeHTTP(lrw, r)
		})
	}
}

// Timeout is a middleware that sets a timeout on the request context.
func Timeout(dur time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), dur)
			defer cancel()
			*r = *r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// Recover is a middleware that catches panics
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				switch err := rec.(type) {
				case error:
					JSONError(w, r, core.WrapErrorWithStack(err, debug.Stack()))
				default:
					toErr := core.WrapErrorWithStack(fmt.Errorf("%v", err), debug.Stack())
					JSONError(w, r, toErr)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RealIP is a middleware that sets a context value with the real ip of the request.
// it uses the defined stragey to determine the real ip.
// See documentation for more information on how to pick the right strategy
func RealIP(strat realclientip.Strategy) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := strat.ClientIP(r.Header, r.RemoteAddr)
			if clientIP == "" {
				err := core.WrapError(
					errors.New("Could not determine client ip"),
					core.ErrBadRequest,
				)
				JSONError(w, r, err)
				return
			}
			// I don't like that. Is there any better way to pass the context
			// to the middleware?
			*r = *r.WithContext(core.NewContextWithClientIP(r.Context(), clientIP))
			next.ServeHTTP(w, r)
		})
	}
}
