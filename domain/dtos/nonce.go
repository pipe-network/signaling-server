package dtos

import (
	"encoding/hex"
	"strconv"
	"strings"
)

const (
	NonceByteLength = 24
)

type Nonce struct {
	Cookie         Cookie
	Source         Address
	Destination    Address
	OverflowNumber OverflowNumber
	SequenceNumber SequenceNumber
}

func NonceFromBytes(bytes [24]byte) Nonce {
	var cookie Cookie
	var source, destination Address
	var overflowNumber OverflowNumber
	var sequenceNumber SequenceNumber
	var from int

	copy(cookie[:], bytes[from:len(cookie)])
	from += len(cookie)
	copy(source[:], bytes[from:from+len(source)])
	from += len(source)
	copy(destination[:], bytes[from:from+len(destination)])
	from += len(destination)
	copy(overflowNumber[:], bytes[from:from+len(overflowNumber)])
	from += len(overflowNumber)
	copy(sequenceNumber[:], bytes[from:from+len(sequenceNumber)])

	return Nonce{
		Cookie:         cookie,
		Source:         source,
		Destination:    destination,
		OverflowNumber: overflowNumber,
		SequenceNumber: sequenceNumber,
	}
}

func (n Nonce) Bytes() []byte {
	var bytes []byte
	bytes = append(bytes, n.Cookie[:]...)
	bytes = append(bytes, n.Source[:]...)
	bytes = append(bytes, n.Destination[:]...)
	bytes = append(bytes, n.OverflowNumber[:]...)
	bytes = append(bytes, n.SequenceNumber[:]...)
	return bytes
}

func (n Nonce) String() string {
	return strings.Join(
		[]string{
			hex.EncodeToString(n.Cookie[:]),
			strconv.Itoa(n.Source.Int()),
			strconv.Itoa(n.Destination.Int()),
			strconv.Itoa(n.OverflowNumber.Int()),
			strconv.Itoa(n.SequenceNumber.Int()),
		},
		"|",
	)
}
