package core

import (
	"context"
)

// ErrorReporter is an interface for reporting errors to an external service
type ErrorReporter interface {
	ReportError(ctx context.Context, args ...any)
	ReportPanic(ctx context.Context, args ...any)
	Close()
}

type StubErrorReporter struct{}

func (*StubErrorReporter) ReportError(context.Context, ...any) {}
func (*StubErrorReporter) ReportPanic(context.Context, ...any) {}
func (*StubErrorReporter) Close()                              {}
