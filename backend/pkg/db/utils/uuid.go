package utils

import (
	"bytes"

	"github.com/google/uuid"
)

var zeroUUID uuid.UUID

func EqualUUID(a, b uuid.UUID) bool {
	return bytes.Equal(a[:], b[:])
}
func IsZeroUUID(u uuid.UUID) bool {
	return EqualUUID(u, zeroUUID)
}
