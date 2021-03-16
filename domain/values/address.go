package values

var (
	UnassignedAddress Address = 0
	ServerAddress     Address = 0
	InitiatorAddress  Address = 1
	MaxAddress        Address = 255
)

type Address int

func (a Address) Bytes() []byte {
	var bytes [1]byte
	bytes[0] = byte(a)
	return bytes[:]
}
