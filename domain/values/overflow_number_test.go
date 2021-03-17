package values

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOverflowNumberFromInt(t *testing.T) {
	overflowNumber := OverflowNumberFromInt(15)
	assert.Equal(t, OverflowNumber{0x00, 0x0F}, overflowNumber)
}

func TestOverflowNumber_Int(t *testing.T) {
	overflowNumber := OverflowNumber{0x00, 0x0F}
	assert.Equal(t, uint16(15), overflowNumber.Int())
}

func TestOverflowNumber_Empty(t *testing.T) {
	assert.True(t, OverflowNumber{}.Empty())
	assert.False(t, OverflowNumber{0x1}.Empty())
}

func TestOverflowNumber_Equal(t *testing.T) {
	assert.True(t, OverflowNumber{}.Equal(OverflowNumber{}))
	assert.False(t, OverflowNumber{0x1}.Equal(OverflowNumber{}))
	assert.True(t, OverflowNumber{0x01}.Equal(OverflowNumber{0x01}))
}
