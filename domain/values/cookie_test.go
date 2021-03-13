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
