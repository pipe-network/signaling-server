package values

import (
	"encoding/binary"
	"math/rand"
)

type SequenceNumber [4]byte

func NewSequenceNumber() (*SequenceNumber, error) {
	sequenceNumber := SequenceNumber{}
	randomIntValue := rand.Uint32()
	binary.BigEndian.PutUint32(sequenceNumber[:], randomIntValue)
	return &sequenceNumber, nil
}

func (a SequenceNumber) Int() int {
	return int(binary.BigEndian.Uint32(a[:]))
}
