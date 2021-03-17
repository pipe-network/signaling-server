package values

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddress_Bytes(t *testing.T) {
	actualInitatorBytes := InitiatorAddress.Bytes()
	assert.Equal(t, []byte{0x1}, actualInitatorBytes)

	actualServerBytes := ServerAddress.Bytes()
	assert.Equal(t, []byte{0x0}, actualServerBytes)

	actualRandomBytes := Address(3).Bytes()
	assert.Equal(t, []byte{0x3}, actualRandomBytes)
}
