package web_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gosom/kit/logging"
	"github.com/gosom/kit/web"
	"github.com/stretchr/testify/require"
)

func TestRequestIDProvider(t *testing.T) {
	id1 := web.RequestIDProvider()
	require.NotEmpty(t, id1)
	require.IsType(t, "", id1)
	id2 := web.RequestIDProvider()
	require.NotEmpty(t, id2)
	require.NotEqual(t, id1, id2)
}

func TestTimeProvider(t *testing.T) {
	t1 := web.TimeProvider()
	require.NotEmpty(t, t1)
	require.IsType(t, time.Time{}, t1)
	t2 := web.TimeProvider()
	require.NotEmpty(t, t2)
	require.Less(t, t1, t2)
}

func TestDecodeBody(t *testing.T) {
	body := `{"foo":"bar"}`
	t.Run("CanReadBodyToMap", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(body))
		require.NoError(t, err)
		var data map[string]string
		err = web.DecodeBody(req, &data, false)
		require.NoError(t, err)
		require.Equal(t, "bar", data["foo"])
	})
	t.Run("CannotReadBodyToMapWhenUsingValidation", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(body))
		require.NoError(t, err)
		var data map[string]string
		err = web.DecodeBody(req, &data, true)
		require.Error(t, err)
	})
	t.Run("CanReadBodyToStructWithValidation", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(body))
		require.NoError(t, err)
		var data struct {
			Foo string `json:"foo" validate:"required"`
		}
		err = web.DecodeBody(req, &data, true)
		require.NoError(t, err)
		require.Equal(t, "bar", data.Foo)
	})
}

func TestStringUrlParam(t *testing.T) {
	var b bytes.Buffer
	b.Reset()
	defaultLogger := logging.Get()
	defer func() {
		logging.SetDefault(defaultLogger)
	}()
	logger := logging.New("zerolog", logging.DEBUG, &b)
	logging.SetDefault(logger)
	r := web.NewRouter(web.RouterConfig{})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := web.StringURLParam(r, "id")
		require.Equal(t, "123", id)
	})

	req := httptest.NewRequest("GET", "/123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)

}
