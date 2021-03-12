package dtos

import (
	"crypto/rand"
	"github.com/pipe-network/signaling-server/domain/domain_services"
	"math"
	"math/big"
)

type SequenceNumber [4]byte

func NewSequenceNumber() (*SequenceNumber, error) {
	sequenceNumber := SequenceNumber{}
	randomIntValue, err := rand.Int(rand.Reader, big.NewInt(int64(math.Pow(2, 32))))
	if err != nil {
		return nil, err
	}
	copy(sequenceNumber[:], randomIntValue.Bytes())
	return &sequenceNumber, nil
}

func (a SequenceNumber) Int() int {
	return domain_services.BytesToInteger(a[:])
}
