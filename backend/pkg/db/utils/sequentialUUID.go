package utils

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

var generatorId uint16
var sequenceIndex atomic.Uint64

const (
	// maxSequenceIndex sequence index can only fill 6 bytes
	maxSequenceIndex uint64 = 0xFFFFFFFFFF
)

func SetGeneratorId(value uint16) { generatorId = value }
func GenerateNewSequentialUUID() uuid.UUID {
	uuid.NewUUID()

	var id uuid.UUID
	n := time.Now().UTC().UnixNano()

	// copy first 8 bytes from the time
	for i := 0; i < 8; i++ {
		id[i] = byte((n >> (56 - (i * 8))) & 0xFF)
	}

	// copy generator ID in next 2 bytes
	for i := 0; i < 2; i++ {
		id[i+8] = byte((generatorId >> (8 - (i * 8))) & 0xFF)
	}

	// and the rest will be the index, since the index should only contains 6 bytes
	// I want to increment the sequence index and wrap it on the 6 bytes capacity
	var index uint64
	for {
		index = sequenceIndex.Load()
		var next uint64
		if index == maxSequenceIndex {
			next = 0
		} else {
			next = index + 1
		}

		if sequenceIndex.CompareAndSwap(index, next) {
			break
		}
	}
	for i := 0; i < 6; i++ {
		id[i+10] = byte((index >> (40 - (i * 8))) & 0xFF)
	}

	return id
}
