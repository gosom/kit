package web

import (
	"net/http"
	"time"

	"github.com/gosom/kit/lib"
	"github.com/realclientip/realclientip-go"

	"github.com/go-chi/chi/v5"
)

type Router interface {
	chi.Router
}

type RouterConfig struct {
	// NotUseDefaultMiddleware if true, the default middleware will not be used
	NotUseDefaultMiddlewares bool
	// NotFoundHandler is the handler to be called when no route matches.
	NotFoundHandler http.HandlerFunc
	// NotFoundHandler is the handler to be called when a method is not allowed.
	NotAllowedHandler http.HandlerFunc
	// RealIPStrategy is the strategy to be used to get the real ip address of the client.
	// see docs: https://pkg.go.dev/github.com/realclientip/realclientip-go
	// By default, it uses the RemoteAddrStrategy (assumes your server is directly connected to the internet)
	// Adjust based on your needs.
	RealIPStrategy realclientip.Strategy
	// Timeout is the timeout to be used for the context of each request
	// Default is 15 seconds. Use -1 to disable.
	Timeout time.Duration

	// CorsCfg by default allows all. Configure properly
	CorsCfg CorsConfig

	// ErrorReporter is the error reporter to be used to report errors
	// By default, it uses the default error reporter (lib.StubErrorReporter)
	ErrorReporter lib.ErrorReporter
}

func NewRouter(cfg RouterConfig) (r Router) {
	r = chi.NewRouter()
	var ipStrategy realclientip.Strategy
	switch cfg.RealIPStrategy {
	case nil:
		ipStrategy = realclientip.RemoteAddrStrategy{}
	default:
		ipStrategy = cfg.RealIPStrategy
	}
	if !cfg.NotUseDefaultMiddlewares {
		var reporter lib.ErrorReporter
		if cfg.ErrorReporter == nil {
			reporter = &lib.StubErrorReporter{}
		} else {
			reporter = cfg.ErrorReporter
		}
		r.Use(RequestLogger(reporter))
		r.Use(Recover)
		r.Use(RealIP(ipStrategy))
		if cfg.Timeout > -1 {
			if cfg.Timeout == 0 {
				cfg.Timeout = 30 * time.Second
			}
			r.Use(Timeout(cfg.Timeout))
		}
		corsMiddleware := NewCors(cfg.CorsCfg)
		r.Use(corsMiddleware.Handler)
	}
	switch {
	case cfg.NotFoundHandler != nil:
		r.NotFound(cfg.NotFoundHandler)
	default:
		r.NotFound(defaultNotFoundHandler)
	}
	switch {
	case cfg.NotAllowedHandler != nil:
		r.MethodNotAllowed(cfg.NotAllowedHandler)
	default:
		r.MethodNotAllowed(defaultMethdoNotAllowed)
	}
	return r
}

func defaultNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	JSONError(w, r, lib.ErrNotFound)
}

func defaultMethdoNotAllowed(w http.ResponseWriter, r *http.Request) {
	JSONError(w, r, lib.ErrMethodNotAllowed)
}
