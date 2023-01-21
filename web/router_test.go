package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gosom/kit/web"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	t.Run("test that router is created with default config", func(t *testing.T) {
		r := web.NewRouter(web.RouterConfig{})
		r.Get("/getOnly", func(w http.ResponseWriter, r *http.Request) {})
		require.NotNil(t, r)
		require.Implements(t, (*web.Router)(nil), r)

		middlewares := r.Middlewares()
		require.Len(t, middlewares, 5)
		for _, m := range middlewares {
			require.NotNil(t, m)
		}
		require.NotNil(t, r.NotFound)
		require.NotNil(t, r.MethodNotAllowed)

		t.Run("Test that route is registered", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getOnly", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("Test that route is not registered", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getOnly2", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusNotFound, w.Code)
			var resp2 web.ErrResponse
			err := json.NewDecoder(w.Body).Decode(&resp2)
			require.NoError(t, err)
			require.Equal(t, 404, resp2.Code)
			require.Equal(t, "not found", resp2.Message)
		})
		t.Run("test that method not allowed is returned", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/getOnly", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusMethodNotAllowed, w.Code)
			var resp2 web.ErrResponse
			err := json.NewDecoder(w.Body).Decode(&resp2)
			require.NoError(t, err)
			require.Equal(t, 405, resp2.Code)
			require.Equal(t, "method not allowed", resp2.Message)
		})
	})
	t.Run("test that router is created with custom config", func(t *testing.T) {
		r := web.NewRouter(web.RouterConfig{
			NotUseDefaultMiddlewares: true,
			NotFoundHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			NotAllowedHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusMethodNotAllowed)
			},
		})
		r.Get("/getOnly", func(w http.ResponseWriter, r *http.Request) {})
		require.NotNil(t, r)
		require.Implements(t, (*web.Router)(nil), r)

		middlewares := r.Middlewares()
		require.Len(t, middlewares, 0)
		require.NotNil(t, r.NotFound)
		require.NotNil(t, r.MethodNotAllowed)

		t.Run("Test that route is registered", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getOnly", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("Test that route is not registered", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getOnly2", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusNotFound, w.Code)
		})
		t.Run("test that method not allowed is returned", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/getOnly", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	})
}
