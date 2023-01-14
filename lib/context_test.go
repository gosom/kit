package lib_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gosom/kit/lib"
	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	t.Run("TestThatYouCanSetAndFetchRequestIDFromContext", func(t *testing.T) {
		ctx := lib.NewContextWithRequestID(context.Background(), "123")
		want := "123"
		got := lib.RequestIDFromContext(ctx)
		require.Equal(t, want, got)
	})
	t.Run("TestThatYouCanSetAndFetchUserIDFromContext", func(t *testing.T) {
		u := User{}
		ctx := lib.NewContextWithUser(context.Background(), &u)
		want := "123"
		got := lib.UserFromContext(ctx).GetID()
		require.Equal(t, want, got)
	})
	t.Run("TestThatYouCanSetAndFetchErrorFromContext", func(t *testing.T) {
		want := errors.New("foo")
		ctx := lib.NewContextWithErr(context.Background(), want)
		got := lib.ErrorFromContext(ctx)
		require.Equal(t, want, got)
	})
	t.Run("TestThatYouCanSetAndFetchClientIPFromContext", func(t *testing.T) {
		want := "127.0.0.1"
		ctx := lib.NewContextWithClientIP(context.Background(), want)
		got := lib.IPFromContext(ctx)
		require.Equal(t, want, got)
	})
}
