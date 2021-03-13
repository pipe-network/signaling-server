package values

import (
	"bytes"
	"encoding/binary"
	log "github.com/sirupsen/logrus"
)

var (
	UnassignedAddress = Address{0x00}
	ServerAddress     = Address{0x00}
	InitiatorAddress  = Address{0x01}
)

type Address [1]byte

func (a Address) Int() int {
	var value uint8
	buf := bytes.NewBuffer(a[:])
	err := binary.Read(buf, binary.BigEndian, &value)
	if err != nil {
		log.Fatalf("Decode failed: %s", err)
	}
	return int(value)
}
