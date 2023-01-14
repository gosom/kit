package lib

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

// MustNewULID returns a new ULID.
// It panics if it fails to generate a new ULID.
// The ULID is a 26-character string and it's monotonically increasing
// It has millisecond precision.
func MustNewULID() string {
	return ulid.Make().String()
}

// NewUUID returns a new UUID.
func NewUUID() string {
	return uuid.New().String()
}
