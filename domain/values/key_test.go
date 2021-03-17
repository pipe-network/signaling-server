package values

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeHex_Success(t *testing.T) {
	actualKey, err := FromHex("55f9fb6377899b8fa6868db0cded3c96349d650a2622d803a4887104918f0227")
	assert.NoError(t, err)
	assert.Equal(
		t,
		&Key{
			0x55, 0xf9, 0xfb, 0x63, 0x77, 0x89, 0x9b, 0x8f, 0xa6, 0x86, 0x8d, 0xb0, 0xcd, 0xed, 0x3c, 0x96,
			0x34, 0x9d, 0x65, 0xa, 0x26, 0x22, 0xd8, 0x3, 0xa4, 0x88, 0x71, 0x4, 0x91, 0x8f, 0x2, 0x27,
		},
		actualKey,
	)
}

func TestDecodeHex_NotValid(t *testing.T) {
	actualKey, err := FromHex("!test!")
	assert.Error(t, err)
	assert.Nil(t, actualKey)
}

func TestDecodeHex_Not32BytesLong(t *testing.T) {
	actualKey, err := FromHex("55f9fb6377899b8fa6868db0cded3c96349d650a2622d803a4887104918f02271234")
	assert.Error(t, err)
	assert.Equal(t, HexKeyNot32BytesLong, err)
	assert.Nil(t, actualKey)
}

func TestKey_Empty(t *testing.T) {
	assert.True(t, Key{}.Empty())
	assert.False(t, Key{0x1}.Empty())
}

func TestKey_Equals(t *testing.T) {
	assert.True(t, Key{}.Equals(Key{}))
	assert.False(t, Key{0x1}.Equals(Key{}))
	assert.True(t, Key{0x1}.Equals(Key{0x1}))
}
