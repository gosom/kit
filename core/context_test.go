package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gosom/kit/core"
	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	t.Run("TestThatYouCanSetAndFetchRequestIDFromContext", func(t *testing.T) {
		ctx := core.NewContextWithRequestID(context.Background(), "123")
		want := "123"
		got := core.RequestIDFromContext(ctx)
		require.Equal(t, want, got)
	})
	t.Run("TestThatYouCanSetAndFetchUserIDFromContext", func(t *testing.T) {
		u := User{}
		ctx := core.NewContextWithUser(context.Background(), &u)
		want := "123"
		got := core.UserFromContext(ctx).GetID()
		require.Equal(t, want, got)
	})
	t.Run("TestThatYouCanSetAndFetchErrorFromContext", func(t *testing.T) {
		want := errors.New("foo")
		ctx := core.NewContextWithErr(context.Background(), want)
		got := core.ErrorFromContext(ctx)
		require.Equal(t, want, got)
	})
	t.Run("TestThatYouCanSetAndFetchClientIPFromContext", func(t *testing.T) {
		want := "127.0.0.1"
		ctx := core.NewContextWithClientIP(context.Background(), want)
		got := core.IPFromContext(ctx)
		require.Equal(t, want, got)
	})
}
