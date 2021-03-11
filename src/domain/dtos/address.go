package dtos


type Address [1]byte

var (
	Server    = Address{0x00}
	Initiator = Address{0x00}
)
