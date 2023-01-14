package core_test

import (
	"testing"

	"github.com/gosom/kit/core"
	"github.com/stretchr/testify/require"
)

func TestNewUUID(t *testing.T) {
	uuid := core.NewUUID()
	require.NotEmpty(t, uuid)
	require.Len(t, uuid, 36)
	require.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", uuid)
}
