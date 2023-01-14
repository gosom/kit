package lib_test

import (
	"errors"
	"testing"

	"github.com/gosom/kit/lib"

	"github.com/stretchr/testify/require"
)

func TestNewApiError(t *testing.T) {
	t.Run("TestThatNewApiErrorReturnsErrorWithStatusAndMessage", func(t *testing.T) {
		want := "foo"
		ae := lib.NewApiError(400, want)
		require.Equal(t, want, ae.Error())
		code, message := ae.ApiError()
		require.Equal(t, 400, code)
		require.Equal(t, want, message)
	})
}

func TestWrapError(t *testing.T) {
	t.Run("TestThatWrapErrorReturnsErrorWithStatusAndMessage", func(t *testing.T) {
		ae := lib.NewApiError(400, "foo")
		err := errors.New("bar")
		we := lib.WrapError(err, ae)
		require.Equal(t, "bar", we.Error())
		require.Implements(t, (*lib.ApiError)(nil), we)
	})
}
