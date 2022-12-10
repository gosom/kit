package logging

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

type zeroLogger struct {
	l zerolog.Logger
}

func newZeroLogger(level Level, w io.Writer) zeroLogger {
	ans := zeroLogger{
		l: zerolog.New(w).Level(zerolog.Level(level)).
			With().Timestamp().Logger(),
	}
	return ans
}

func (o zeroLogger) Info(msg string, args ...any) {
	o.Log(INFO, msg, args...)
}

func (o zeroLogger) Warn(msg string, args ...any) {
	o.Log(WARN, msg, args...)
}

func (o zeroLogger) Error(msg string, args ...any) {
	o.Log(ERROR, msg, args...)
}

func (o zeroLogger) Debug(msg string, args ...any) {
	o.Log(DEBUG, msg, args...)
}

func (o zeroLogger) Trace(msg string, args ...any) {
	o.Log(TRACE, msg, args...)
}

func (o zeroLogger) Fatal(msg string, args ...any) {
	o.Log(FATAL, msg, args...)
}

func (o zeroLogger) Panic(msg string, args ...any) {
	o.Log(PANIC, msg, args...)
}

func (o zeroLogger) Log(level Level, msg string, args ...any) {
	ev := o.l.WithLevel(zerolog.Level(level))
	for len(args) > 0 {
		ev, args = set(ev, args)
	}
	switch msg {
	case "":
		ev.Send()
	default:
		ev.Msg(msg)
	}
}

func (o zeroLogger) With(args ...any) Logger {
	ans := zeroLogger{}
	c := o.l.With()
	for len(args) > 0 {
		c, args = setWith(c, args)
	}
	ans.l = c.Logger()
	return ans
}

func (o zeroLogger) Level(level Level) Logger {
	ans := zeroLogger{
		l: o.l.Level(zerolog.Level(level)),
	}
	return ans
}

// WithContext returns a new context with the logger attached.
func (o zeroLogger) NewContext(ctx context.Context) context.Context {
	if l, ok := ctx.Value(ctxKey{}).(*zeroLogger); ok && l == &o {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, &o)
}

// -----------------------------------------------------------------------------

func setWith(ev zerolog.Context, args []any) (zerolog.Context, []any) {
	switch k := args[0].(type) {
	case string:
		if len(args) == 1 {
			ev = mapZerologContext(ev, missingKey, k)
			return ev, nil
		}
		ev = mapZerologContext(ev, k, args[1])
		return ev, args[2:]
	default:
		ev = mapZerologContext(ev, missingKey, k)
		return ev, args[1:]
	}
}

func mapZerologContext(ev zerolog.Context, key string, value any) zerolog.Context {
	switch v := value.(type) {
	case string:
		return ev.Str(key, v)
	case int:
		return ev.Int(key, v)
	case int64:
		return ev.Int64(key, v)
	case uint:
		return ev.Uint(key, v)
	case uint64:
		return ev.Uint64(key, v)
	case float32:
		return ev.Float32(key, v)
	case float64:
		return ev.Float64(key, v)
	case bool:
		return ev.Bool(key, v)
	case time.Time:
		return ev.Time(key, v)
	case time.Duration:
		return ev.Dur(key, v)
	case error:
		return errorWithStackContext(ev, key, v)
	case []byte:
		return ev.Bytes(key, v)
	case fmt.Stringer:
		return ev.Str(key, v.String())
	case fmt.GoStringer:
		return ev.Str(key, v.GoString())
	case nil:
		return ev.Interface(key, nil)
	default:
		return ev.Interface(key, v)
	}
}

func set(ev *zerolog.Event, args []any) (*zerolog.Event, []any) {
	switch k := args[0].(type) {
	case string:
		if len(args) == 1 {
			ev = mapZerolog(ev, missingKey, k)
			return ev, nil
		}
		ev = mapZerolog(ev, k, args[1])
		return ev, args[2:]
	default:
		ev = mapZerolog(ev, missingKey, k)
		return ev, args[1:]
	}
}

func mapZerolog(ev *zerolog.Event, key string, value any) *zerolog.Event {
	switch v := value.(type) {
	case string:
		return ev.Str(key, v)
	case int:
		return ev.Int(key, v)
	case int64:
		return ev.Int64(key, v)
	case uint:
		return ev.Uint(key, v)
	case uint64:
		return ev.Uint64(key, v)
	case float32:
		return ev.Float32(key, v)
	case float64:
		return ev.Float64(key, v)
	case bool:
		return ev.Bool(key, v)
	case time.Time:
		return ev.Time(key, v)
	case time.Duration:
		return ev.Dur(key, v)
	case error:
		return errorWithStack(ev, key, v)
	case []byte:
		return ev.Bytes(key, v)
	case fmt.Stringer:
		return ev.Str(key, v.String())
	case fmt.GoStringer:
		return ev.Str(key, v.GoString())
	case nil:
		return ev.Interface(key, nil)
	default:
		return ev.Interface(key, v)
	}
}

func errorWithStack(ev *zerolog.Event, key string, err error) *zerolog.Event {
	se, ok := err.(stackTracer)
	if ok {
		ev.Interface("stack", stackFormater(se.StackTrace()))
	}
	return ev.AnErr(key, err)
}

func errorWithStackContext(ev zerolog.Context, key string, err error) zerolog.Context {
	se, ok := err.(stackTracer)
	if ok {
		ev.Interface("stack", stackFormater(se.StackTrace()))
	}
	return ev.AnErr(key, err)
}

func stackFormater(s string) []string {
	s = strings.Replace(s, "\t", "", -1)
	l := strings.Split(s, "\n")
	n := 0
	for i := range l {
		if l[i] != "" {
			l[n] = l[i]
			n++
		}
	}
	return l[:n]
}
