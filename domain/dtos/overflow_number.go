package dtos

import "github.com/pipe-network/signaling-server/domain/domain_services"

type OverflowNumber [2]byte

func NewOverflowNumber() OverflowNumber {
	return [2]byte{}
}

func (n OverflowNumber) Int() int {
	return domain_services.BytesToInteger(n[:])
}
