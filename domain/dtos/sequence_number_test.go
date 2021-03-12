package dtos

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCombinedSequenceNumber(t *testing.T) {
	combinedSequenceNumber, err := NewSequenceNumber()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x0, 0x0}, combinedSequenceNumber[0:2])
	assert.NotEqual(t, []byte{0x0, 0x0, 0x0}, combinedSequenceNumber[:2])
}
