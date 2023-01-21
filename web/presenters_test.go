package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gosom/kit/lib"
	"github.com/gosom/kit/web"
	"github.com/stretchr/testify/require"
)

func TestJSONPresenter(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	data := map[string]interface{}{
		"foo": "bar",
	}
	web.JSON(w, req, 200, data)
	require.Equal(t, 200, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	var expectedBody bytes.Buffer
	err := json.NewEncoder(&expectedBody).Encode(data)
	require.NoError(t, err)
	require.Equal(t, expectedBody.String(), w.Body.String())
}

func TestJSONErrorPresenter(t *testing.T) {
	t.Run("test that error is set on context", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		err := errors.New("test error")
		web.JSONError(w, req, err)
		require.Equal(t, 500, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		ctxErr := lib.ErrorFromContext(req.Context())
		require.Equal(t, err, ctxErr)
	})
	t.Run("when error is ApiError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		ae := lib.NewApiError(400, "foo")
		web.JSONError(w, req, ae)
		require.Equal(t, 400, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		var resp web.ErrResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, 400, resp.Code)
		require.Equal(t, "foo", resp.Message)
	})
	t.Run("when error is validation error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		type testStruct struct {
			Foo string `json:"foo" validate:"required"`
		}
		err := lib.Validate(testStruct{})
		web.JSONError(w, req, err)
		require.Equal(t, 400, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp web.ErrResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, 400, resp.Code)
		require.Equal(t, "Key: 'testStruct.Foo' Error:Field validation for 'Foo' failed on the 'required' tag", resp.Message)
	})
	t.Run("when error is not ApiError and not validation error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		err := errors.New("foo")
		web.JSONError(w, req, err)
		require.Equal(t, 500, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		var resp web.ErrResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, 500, resp.Code)
		require.Equal(t, http.StatusText(500), resp.Message)
	})
}
