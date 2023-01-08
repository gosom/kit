package lib

import (
	"hash/fnv"

	"github.com/oklog/ulid/v2"
)

// HashToUInt32 returns a hash of the given string.
// Uses the FNV-1a algorithm.
func HashToUInt32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// Int32Ring returns a number in the range [0, 2^31).
func Int32Ring(v uint32) int32 {
	return int32(v & 0x7fffffff)
}

// MustNewULID returns a new ULID.
// It panics if it fails to generate a new ULID.
// The ULID is a 26-character string and it's monotonically increasing
// It has millisecond precision.
func MustNewULID() string {
	return ulid.Make().String()
}
