package lib_test

import (
	"testing"

	"github.com/gosom/kit/lib"
	"github.com/stretchr/testify/require"
)

func TestNewUUID(t *testing.T) {
	uuid := lib.NewUUID()
	require.NotEmpty(t, uuid)
	require.Len(t, uuid, 36)
	require.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", uuid)
}

func TestMustNewUUID(t *testing.T) {
	ulid1 := lib.MustNewULID()
	require.NotEmpty(t, ulid1)
	require.Len(t, ulid1, 26)
	ulid2 := lib.MustNewULID()
	require.NotEmpty(t, ulid2)
	require.Len(t, ulid2, 26)

	require.NotEqual(t, ulid1, ulid2)
	require.Less(t, ulid1, ulid2)
}
