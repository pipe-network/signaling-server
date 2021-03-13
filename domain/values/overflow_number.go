package values

import (
	"encoding/binary"
)

type OverflowNumber [2]byte

func NewOverflowNumber() OverflowNumber {
	return [2]byte{0x00, 0x00}
}

func (n OverflowNumber) Int() int {
	return int(binary.BigEndian.Uint16(n[:]))
}
