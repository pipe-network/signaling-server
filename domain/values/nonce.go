package values

import (
	"encoding/hex"
	"fmt"
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

func NonceFromBytes(bytes [NonceByteLength]byte) Nonce {
	var cookie Cookie
	var source, destination Address
	var overflowNumber OverflowNumber
	var sequenceNumber SequenceNumber
	var from int

	sourceBytes := make([]byte, 1)
	destinationBytes := make([]byte, 1)

	copy(cookie[:], bytes[from:len(cookie)])
	from += len(cookie)
	copy(sourceBytes[:], bytes[from:from+len(sourceBytes)])
	from += len(source.Bytes())
	copy(destinationBytes[:], bytes[from:from+len(destinationBytes)])
	from += len(destination.Bytes())
	copy(overflowNumber[:], bytes[from:from+len(overflowNumber)])
	from += len(overflowNumber)
	copy(sequenceNumber[:], bytes[from:from+len(sequenceNumber)])

	if len(sourceBytes) > 0 {
		source = Address(int(sourceBytes[0]))
	}

	if len(destinationBytes) > 0 {
		destination = Address(int(destinationBytes[0]))
	}

	return Nonce{
		Cookie:         cookie,
		Source:         source,
		Destination:    destination,
		OverflowNumber: overflowNumber,
		SequenceNumber: sequenceNumber,
	}
}

func (n Nonce) Bytes() [NonceByteLength]byte {
	var bytes []byte
	var nonceBytes [NonceByteLength]byte

	bytes = append(bytes, n.Cookie[:]...)
	bytes = append(bytes, n.Source.Bytes()...)
	bytes = append(bytes, n.Destination.Bytes()...)
	bytes = append(bytes, n.OverflowNumber[:]...)
	bytes = append(bytes, n.SequenceNumber[:]...)

	copy(nonceBytes[:], bytes[:])
	return nonceBytes
}

func (n Nonce) String() string {
	return strings.Join(
		[]string{
			fmt.Sprintf("OutgoingCookie: %s", hex.EncodeToString(n.Cookie[:])),
			fmt.Sprintf("Source: %s", strconv.Itoa(int(n.Source))),
			fmt.Sprintf("Destination: %s", strconv.Itoa(int(n.Destination))),
			fmt.Sprintf("OutgoingOverflowNumber: %s", strconv.Itoa(int(n.OverflowNumber.Int()))),
			fmt.Sprintf("OutgoingSequenceNumber: %s", strconv.Itoa(int(n.SequenceNumber.Int()))),
		},
		"|",
	)
}
