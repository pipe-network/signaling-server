package values

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRandomCookie(t *testing.T) {
	cookie, err := NewRandomCookie()
	assert.NoError(t, err)
	assert.NotEqual(t, [16]byte{}, cookie[:])
}

func TestCookie_Equal(t *testing.T) {
	assert.True(t, Cookie{}.Equal(Cookie{}))
	assert.False(t, Cookie{0x1}.Equal(Cookie{}))
}

func TestCookie_Empty(t *testing.T) {
	assert.True(t, Cookie{}.Empty())
	assert.False(t, Cookie{0x1}.Empty())
}
