package logging

import (
	"context"
	"io"
	"os"
	"sync"
)

const missingKey = "!MISSING!"

type ctxKey struct{}

type Level int8

func (o Level) String() string {
	switch o {
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	case DEBUG:
		return "debug"
	case TRACE:
		return "trace"
	case FATAL:
		return "fatal"
	case PANIC:
		return "panic"
	default:
		return "UNKNOWN"
	}
}

var lock sync.Mutex

const (
	TRACE    Level = -1
	DEBUG    Level = 0
	INFO     Level = 1
	WARN     Level = 2
	ERROR    Level = 3
	FATAL    Level = 4
	PANIC    Level = 5
	DISABLED Level = 6
)

// Logger is the interface that wraps the basic logging methods.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Trace(msg string, args ...any)
	Fatal(msg string, args ...any)
	Panic(msg string, args ...any)
	Log(level Level, msg string, args ...any)

	With(args ...any) Logger
	Level(level Level) Logger
	NewContext(ctx context.Context) context.Context
}

// New returns a new Logger that writes to the given io.Writer.
func New(implName string, level Level, w io.Writer) Logger {
	switch implName {
	case "zerolog":
		return newZeroLogger(level, w)
	default:
		panic("unknown logger implementation")
	}
}

// std is the default logger.
var std = New("zerolog", INFO, os.Stderr)

func Get() Logger {
	return std
}

// NewContext returns a new context with the given logger.
func NewContext(ctx context.Context, l Logger) context.Context {
	return l.NewContext(ctx)
}

// Sets the default logger to l.
func SetDefault(l Logger) {
	lock.Lock()
	defer lock.Unlock()
	std = l
}

func Ctx(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return l
	}
	return std
}

// Log uses the default logger to log a message at the given level.
func Log(level Level, msg string, args ...any) {
	std.Log(level, msg, args...)
}

// Info uses the default logger to log a message at the INFO level.
func Info(msg string, args ...any) {
	std.Info(msg, args...)
}

// Warn uses the default logger to log a message at the WARN level.
func Warn(msg string, args ...any) {
	std.Warn(msg, args...)
}

// Error uses the default logger to log a message at the ERROR level.
func Error(msg string, args ...any) {
	std.Error(msg, args...)
}

// Debug uses the default logger to log a message at the DEBUG level.
func Debug(msg string, args ...any) {
	std.Debug(msg, args...)
}

// Trace uses the default logger to log a message at the TRACE level.
func Trace(msg string, args ...any) {
	std.Trace(msg, args...)
}

// Panic uses the default logger to log a message at the PANIC level.
func Panic(msg string, args ...any) {
	std.Panic(msg, args...)
}

// Fatal uses the default logger to log a message at the FATAL level.
func Fatal(msg string, args ...any) {
	std.Fatal(msg, args...)
}

type stackTracer interface {
	StackTrace() string
}
