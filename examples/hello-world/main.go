package main

import (
	"context"
	"net/http"

	"github.com/gosom/kit/logging"
	"github.com/gosom/kit/web"
)

func main() {
	ctx := context.Background()

	mux := web.NewRouter(web.RouterConfig{})
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log := logging.Ctx(r.Context())
		log.Info("in Hello world")
		web.JSON(w, r, http.StatusOK, "hello world")
	})
	cfg := web.ServerConfig{
		Router:   mux,
		LogLevel: logging.INFO,
	}

	if err := web.ServerRun(ctx, cfg); err != nil {
		panic(err)
	}

}
