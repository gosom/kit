package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/gosom/kit/logging"
	"github.com/gosom/kit/rollbar"
	"github.com/gosom/kit/web"
)

func main() {
	ctx := context.Background()

	rollbarReporter := rollbar.NewRollbarErrorReporter(
		os.Getenv("ROLLBAR_TOKEN"), "development", "", "", "",
	)
	defer rollbarReporter.Close()

	rCfg := web.RouterConfig{
		ErrorReporter: rollbarReporter,
	}

	mux := web.NewRouter(rCfg)
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log := logging.Ctx(r.Context())
		log.Info("in Hello world")
		web.JSON(w, r, http.StatusOK, "hello world")
	})
	mux.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("in panic")
		web.JSON(w, r, http.StatusOK, "hello world")
	})
	mux.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("in error")
		web.JSONError(w, r, err)
	})
	cfg := web.ServerConfig{
		Router:   mux,
		LogLevel: logging.INFO,
	}

	server := web.NewHttpServer(cfg)
	if err := server.Start(ctx); err != nil {
		panic(err)
	}
}
