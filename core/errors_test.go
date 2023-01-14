package core_test

import (
	"errors"
	"testing"

	"github.com/gosom/kit/core"

	"github.com/stretchr/testify/require"
)

func TestNewApiError(t *testing.T) {
	t.Run("TestThatNewApiErrorReturnsErrorWithStatusAndMessage", func(t *testing.T) {
		want := "foo"
		ae := core.NewApiError(400, want)
		require.Equal(t, want, ae.Error())
		code, message := ae.ApiError()
		require.Equal(t, 400, code)
		require.Equal(t, want, message)
	})
}

func TestWrapError(t *testing.T) {
	t.Run("TestThatWrapErrorReturnsErrorWithStatusAndMessage", func(t *testing.T) {
		ae := core.NewApiError(400, "foo")
		err := errors.New("bar")
		we := core.WrapError(err, ae)
		require.Equal(t, "bar", we.Error())
		require.Implements(t, (*core.ApiError)(nil), we)
	})
}
