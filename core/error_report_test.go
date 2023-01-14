package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gosom/kit/core"
)

func TestStubErrorReporter(t *testing.T) {
	t.Run("TestThatStubErrorReporterImplementsErrorReporter", func(t *testing.T) {
		var _ core.ErrorReporter = &core.StubErrorReporter{}
	})
	t.Run("TestThatStubErrorReporterDoesNothing", func(t *testing.T) {
		er := core.StubErrorReporter{}
		er.ReportError(context.Background(), errors.New("foo"))
		er.ReportPanic(context.Background(), errors.New("foo"))
		er.Close()
	})
}
