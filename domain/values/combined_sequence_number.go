package values

import (
	"errors"
	"math"
)

var (
	OverflowReached = errors.New("overflow of combined sequence number reached")
)

type CombinedSequenceNumber struct {
	SequenceNumber SequenceNumber
	OverflowNumber OverflowNumber
}

func NewCombinedSequenceNumber(sequenceNumber SequenceNumber, overflowNumber OverflowNumber) CombinedSequenceNumber {
	return CombinedSequenceNumber{
		SequenceNumber: sequenceNumber,
		OverflowNumber: overflowNumber,
	}
}

func (n CombinedSequenceNumber) Equal(combinedSequenceNumber CombinedSequenceNumber) bool {
	return n.SequenceNumber.Equal(combinedSequenceNumber.SequenceNumber) &&
		n.OverflowNumber.Equal(combinedSequenceNumber.OverflowNumber)
}

func (n CombinedSequenceNumber) Empty() bool {
	return n.SequenceNumber.Empty() && n.OverflowNumber.Empty()
}

func (n CombinedSequenceNumber) Increment() (*CombinedSequenceNumber, error) {
	sequenceNumberValue := n.SequenceNumber.Int()
	overflowNumberValue := n.OverflowNumber.Int()
	if sequenceNumberValue == math.MaxUint32 {
		if overflowNumberValue == math.MaxUint16 {
			return nil, OverflowReached
		}
		sequenceNumberValue = 0
		overflowNumberValue += 1
	} else {
		sequenceNumberValue += 1
	}

	return &CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(sequenceNumberValue),
		OverflowNumber: OverflowNumberFromInt(overflowNumberValue),
	}, nil
}
