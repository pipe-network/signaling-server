package dtos

import (
	"github.com/pipe-network/signaling-server/domain/domain_services"
)

var (
	Unassigned = Address{0x00}
	Server     = Address{0x00}
	Initiator  = Address{0x01}
)

type Address [1]byte

func (a Address) Int() int {
	return domain_services.BytesToInteger(a[:])
}
