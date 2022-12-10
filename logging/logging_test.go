package logging_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gosom/kit/logging"
	"github.com/stretchr/testify/require"
)

type baseLogJSON struct {
	Level   string    `json:"level"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type tObject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (o tObject) String() string {
	return fmt.Sprintf("<tObject>%d:%s</tObject>", o.ID, o.Name)
}

type pill int

func (pill) GoString() string {
	return "pill"
}

const (
	Placebo pill = iota
)

func TestLog(t *testing.T) {
	var b bytes.Buffer
	b.Reset()
	logger := logging.New("zerolog", logging.TRACE, &b)
	logging.SetDefault(logger)
	funcs := map[string]func(string, ...any){
		"info":  logging.Info,
		"warn":  logging.Warn,
		"error": logging.Error,
		"debug": logging.Debug,
		"trace": logging.Trace,
		"panic": logging.Panic,
		"fatal": logging.Fatal,
	}
	for k, f := range funcs {
		b.Reset()
		f("test", "ok")
		checkLog(t, &b, k, "test")
		require.Contains(t, b.String(), `"!MISSING!":"ok"`)
		b.Reset()
		now := time.Now()
		keyval := []any{
			"int", 1, "string", "s1", "bool", true, "float", 1.1,
			"error", errors.New("test error"), "object", tObject{ID: 1, Name: "test"},
			"byte", []byte("test byte"), "ts", now,
			"int64", int64(1),
			"uint", uint(1), "uint64", uint64(1),
			"float32", float32(1.1), "float64", float64(1.1),
			"duration", time.Second, "pill", Placebo,
			"map", map[string]interface{}{"key": "value"},
			"nil", nil, "missval",
		}
		f("test", keyval...)
		checkLog(t, &b, k, "test")
		require.Contains(t, b.String(), `"int":1`)
		require.Contains(t, b.String(), `"string":"s1"`)
		require.Contains(t, b.String(), `"bool":true`)
		require.Contains(t, b.String(), `"float":1.1`)
		require.Contains(t, b.String(), `"error":"test error"`)
		require.Contains(t, b.String(), `"object":"<tObject>1:test</tObject>"`)
		require.Contains(t, b.String(), `"byte":"test byte"`)
		require.Contains(t, b.String(), `"ts":"`+now.Format(time.RFC3339Nano)+`"`)
		require.Contains(t, b.String(), `"int64":1`)
		require.Contains(t, b.String(), `"uint":1`)
		require.Contains(t, b.String(), `"uint64":1`)
		require.Contains(t, b.String(), `"uint64":1`)
		require.Contains(t, b.String(), `"float32":1.1`)
		require.Contains(t, b.String(), `"float64":1.1`)
		require.Contains(t, b.String(), `"duration":1000`)
		require.Contains(t, b.String(), `"pill":"pill"`)
		require.Contains(t, b.String(), `"map":{"key":"value"}`)
		require.Contains(t, b.String(), `"nil":null`)
		require.Contains(t, b.String(), `"!MISSING!":"missval"`)
		b.Reset()
		for lv := logging.TRACE; lv <= logging.ERROR; lv++ {
			b.Reset()
			logging.Log(lv, "test", keyval...)
			checkLog(t, &b, lv.String(), "test")
			require.Contains(t, b.String(), `"int":1`)
			require.Contains(t, b.String(), `"string":"s1"`)
			require.Contains(t, b.String(), `"bool":true`)
			require.Contains(t, b.String(), `"float":1.1`)
			require.Contains(t, b.String(), `"error":"test error"`)
			require.Contains(t, b.String(), `"object":"<tObject>1:test</tObject>"`)
			require.Contains(t, b.String(), `"byte":"test byte"`)
			require.Contains(t, b.String(), `"ts":"`+now.Format(time.RFC3339Nano)+`"`)
			require.Contains(t, b.String(), `"int64":1`)
			require.Contains(t, b.String(), `"uint":1`)
			require.Contains(t, b.String(), `"uint64":1`)
			require.Contains(t, b.String(), `"uint64":1`)
			require.Contains(t, b.String(), `"float32":1.1`)
			require.Contains(t, b.String(), `"float64":1.1`)
			require.Contains(t, b.String(), `"duration":1000`)
			require.Contains(t, b.String(), `"pill":"pill"`)
			require.Contains(t, b.String(), `"map":{"key":"value"}`)
			require.Contains(t, b.String(), `"nil":null`)
			require.Contains(t, b.String(), `"!MISSING!":"missval"`)
		}
	}
}

func checkLog(t *testing.T, buff *bytes.Buffer, level string, message string) {
	t.Helper()
	var l baseLogJSON
	err := json.Unmarshal(buff.Bytes(), &l)
	require.NoError(t, err)
	require.Equal(t, level, l.Level)
	require.Equal(t, message, l.Message)
	require.WithinDuration(t, time.Now().UTC(), l.Time, time.Second)
}

func BenchmarkLog(b *testing.B) {
	var buff bytes.Buffer
	l := logging.New("zerolog", logging.INFO, &buff)
	b.Run("message only", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Log(logging.INFO, "test message")
		}
	})
	b.Run("log a message and 15 fields", func(b *testing.B) {
		now := time.Now().UTC()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logging.Info("test message",
				"k1", "v1",
				"k2", 1,
				"k3", int64(2),
				"k4", uint(3),
				"k5", uint64(4),
				"k6", float32(3.14),
				"k7", float64(4.14),
				"k8", true,
				"k9", now,
				"k10", (1000 * time.Second),
				"k11", errors.New("test error"),
				"k12", []byte("test bytes"),
				"k13", tObject{ID: 1, Name: "test"},
				"k14", nil,
				"k15", map[string]any{"k1": "v1", "k2": 2},
			)
		}
	})
}
