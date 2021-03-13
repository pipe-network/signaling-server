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
			fmt.Sprintf("Cookie: %s", hex.EncodeToString(n.Cookie[:])),
			fmt.Sprintf("Source: %s", strconv.Itoa(n.Source.Int())),
			fmt.Sprintf("Destination: %s", strconv.Itoa(n.Destination.Int())),
			fmt.Sprintf("OverflowNumber: %s", strconv.Itoa(n.OverflowNumber.Int())),
			fmt.Sprintf("SequenceNumber: %s", strconv.Itoa(n.SequenceNumber.Int())),
		},
		"|",
	)
}
