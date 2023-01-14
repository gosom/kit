package lib_test

import (
	"os"
	"testing"

	"github.com/gosom/kit/lib"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("TestThatNewConfigReturnsConfig", func(t *testing.T) {
		type TestConfig struct {
			Dummy string
		}
		os.Setenv("DUMMY", "foo")
		defer os.Unsetenv("DUMMY")
		c, err := lib.NewConfig[TestConfig]("")
		require.NoError(t, err)
		require.Equal(t, "foo", c.Dummy)
	})
}
