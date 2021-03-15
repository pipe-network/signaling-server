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

func AddressFromInt(value uint8) Address {
	var address Address
	address[0] = value
	return address
}

func (a Address) Int() uint8 {
	var value uint8
	buf := bytes.NewBuffer(a[:])
	err := binary.Read(buf, binary.BigEndian, &value)
	if err != nil {
		log.Fatalf("Decode failed: %s", err)
	}
	return value
}
