package lib_test

import (
	"testing"

	"github.com/gosom/kit/lib"
	"github.com/stretchr/testify/require"
)

func TestHashToUint32(t *testing.T) {
	t.Run("TestThatHashToUint32ReturnsUint32", func(t *testing.T) {
		hash := "foo"
		var want uint32 = 2851307223
		got := lib.HashToUInt32(hash)
		require.Equal(t, want, got)
	})
}

func TestInt32Ring(t *testing.T) {
	t.Run("TestThatInt32RingReturnsInt32", func(t *testing.T) {
		var input uint32 = 4294967295
		var want int32 = 2147483647
		got := lib.Int32Ring(input)
		require.Equal(t, want, got)
	})
}
