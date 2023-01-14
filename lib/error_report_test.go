package lib_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gosom/kit/lib"
)

func TestStubErrorReporter(t *testing.T) {
	t.Run("TestThatStubErrorReporterImplementsErrorReporter", func(t *testing.T) {
		var _ lib.ErrorReporter = &lib.StubErrorReporter{}
	})
	t.Run("TestThatStubErrorReporterDoesNothing", func(t *testing.T) {
		er := lib.StubErrorReporter{}
		er.ReportError(context.Background(), errors.New("foo"))
		er.ReportPanic(context.Background(), errors.New("foo"))
		er.Close()
	})
}
