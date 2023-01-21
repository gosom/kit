package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/realclientip/realclientip-go"
	"github.com/rs/cors"
	"github.com/stretchr/testify/require"

	"github.com/gosom/kit/lib"
	"github.com/gosom/kit/logging"
	"github.com/gosom/kit/web"
)

func TestNewCorsMiddleware(t *testing.T) {
	c := web.NewCors(web.CorsConfig{})
	require.NotNil(t, c)
	require.IsType(t, &cors.Cors{}, c)
}

func TestRequestLoggerMiddleware(t *testing.T) {
	var b bytes.Buffer
	b.Reset()
	defaultLogger := logging.Get()
	defer func() {
		logging.SetDefault(defaultLogger)
	}()
	logger := logging.New("zerolog", logging.DEBUG, &b)
	logging.SetDefault(logger)

	mw := web.RequestLogger(&lib.StubErrorReporter{})
	require.NotNil(t, mw)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	mw(next).ServeHTTP(w, req)

	v := make(map[string]any)
	err := json.Unmarshal(b.Bytes(), &v)
	require.NoError(t, err)
	require.Equal(t, "info", v["level"])
	require.NotEmpty(t, v["request_id"])
	require.Equal(t, 200., v["status"])
	require.Equal(t, 0., v["bytes"])
	require.Equal(t, "GET", v["method"])
	require.Equal(t, "/", v["path"])
	require.Equal(t, "", v["query"])
	require.Equal(t, "", v["ip"])
	require.Equal(t, "", v["user-agent"])
}

func TestTimeoutMiddleware(t *testing.T) {
	mw := web.Timeout(50 * time.Millisecond)
	require.NotNil(t, mw)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for {
			select {
			case <-r.Context().Done():
				ae := lib.NewApiError(http.StatusRequestTimeout, "timeout")
				web.JSONError(w, r, ae)
				return
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
		web.JSON(w, r, http.StatusOK, `{"foo":"bar"}`)
	})
	mw(next).ServeHTTP(w, req)
	require.Equal(t, http.StatusRequestTimeout, w.Code)
}

func TestRecoverMiddleware(t *testing.T) {
	var b bytes.Buffer
	b.Reset()
	defaultLogger := logging.Get()
	defer func() {
		logging.SetDefault(defaultLogger)
	}()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})
	web.Recover(next).ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRealIPMiddleware(t *testing.T) {
	ipStrategy := realclientip.RemoteAddrStrategy{}
	mw := web.RealIP(ipStrategy)
	require.NotNil(t, mw)
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		web.JSON(w, r, http.StatusOK, `{"foo":"bar"}`)
	})

	mw(next).ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	ip := lib.IPFromContext(req.Context())
	require.Equal(t, "127.0.0.1", ip)
}
