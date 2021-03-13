package values

import (
	"bytes"
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

func SequenceNumberFromInt(value uint32) SequenceNumber {
	var byteValue [4]byte
	binary.BigEndian.PutUint32(byteValue[:], value)
	return byteValue
}

func (n SequenceNumber) Int() uint32 {
	return binary.BigEndian.Uint32(n[:])
}

func (n SequenceNumber) Empty() bool {
	emptySequenceNumber := SequenceNumber{}
	return bytes.Equal(n[:], emptySequenceNumber[:])
}

func (n SequenceNumber) Equal(sequenceNumber SequenceNumber) bool {
	return bytes.Equal(n[:], sequenceNumber[:])
}
