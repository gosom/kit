package web_test

import (
	"context"
	"testing"
	"time"

	"github.com/gosom/kit/web"
	"github.com/stretchr/testify/require"
)

func TestNewHttpServer(t *testing.T) {
	t.Run("CanCreateHttpServer", func(t *testing.T) {
		cfg := web.ServerConfig{}
		s := web.NewHttpServer(cfg)
		require.NotNil(t, s)
		require.IsType(t, &web.HttpServer{}, s)
	})
	t.Run("CannotListenAndServeWithoutHandler", func(t *testing.T) {
		cfg := web.ServerConfig{}
		s := web.NewHttpServer(cfg)
		err := s.ListenAndServe(context.Background())
		require.Error(t, err)
	})
	t.Run("CanListenAndServeWithHandler", func(t *testing.T) {
		cfg := web.ServerConfig{
			Host:   "127.0.0.1:0",
			Router: web.NewRouter(web.RouterConfig{}),
		}
		s := web.NewHttpServer(cfg)
		ctx, cancel := context.WithDeadline(context.Background(), web.TimeProvider().Add(100*time.Millisecond))
		defer cancel()
		errc := make(chan error, 1)
		go func() {
			errc <- s.ListenAndServe(ctx)
		}()
		err := <-errc
		require.NoError(t, err)
	})
}
