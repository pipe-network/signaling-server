package values

import (
	"bytes"
	"encoding/binary"
)

type OverflowNumber [2]byte

func NewOverflowNumber() OverflowNumber {
	return [2]byte{0x00, 0x00}
}

func OverflowNumberFromInt(value uint16) OverflowNumber {
	var byteValue [2]byte
	binary.BigEndian.PutUint16(byteValue[:], value)
	return byteValue
}

func (n OverflowNumber) Int() uint16 {
	return binary.BigEndian.Uint16(n[:])
}

func (n OverflowNumber) Empty() bool {
	emptyOverflowNumber := OverflowNumber{}
	return bytes.Equal(n[:], emptyOverflowNumber[:])
}

func (n OverflowNumber) Equal(overflowNumber OverflowNumber) bool {
	return bytes.Equal(n[:], overflowNumber[:])
}
