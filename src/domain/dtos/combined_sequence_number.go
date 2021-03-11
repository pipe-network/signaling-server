package dtos

import "crypto/rand"

type CombinedSequenceNumber [6]byte

func NewCombinedSequenceNumber() (*CombinedSequenceNumber, error) {
	sequenceNumber := CombinedSequenceNumber{}
	_, err := rand.Read(sequenceNumber[2:])
	if err != nil {
		return nil, err
	}
	return &sequenceNumber, nil
}
