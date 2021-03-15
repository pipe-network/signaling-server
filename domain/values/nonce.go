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

func (n Nonce) Bytes() [NonceByteLength]byte {
	var bytes []byte
	var nonceBytes [NonceByteLength]byte

	bytes = append(bytes, n.Cookie[:]...)
	bytes = append(bytes, n.Source[:]...)
	bytes = append(bytes, n.Destination[:]...)
	bytes = append(bytes, n.OverflowNumber[:]...)
	bytes = append(bytes, n.SequenceNumber[:]...)

	copy(nonceBytes[:], bytes[:])
	return nonceBytes
}

func (n Nonce) String() string {
	return strings.Join(
		[]string{
			fmt.Sprintf("OutgoingCookie: %s", hex.EncodeToString(n.Cookie[:])),
			fmt.Sprintf("Source: %s", strconv.Itoa(int(n.Source.Int()))),
			fmt.Sprintf("Destination: %s", strconv.Itoa(int(n.Destination.Int()))),
			fmt.Sprintf("OutgoingOverflowNumber: %s", strconv.Itoa(int(n.OverflowNumber.Int()))),
			fmt.Sprintf("OutgoingSequenceNumber: %s", strconv.Itoa(int(n.SequenceNumber.Int()))),
		},
		"|",
	)
}
