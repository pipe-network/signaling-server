package values

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCombinedSequenceNumber(t *testing.T) {
	sequenceNumber, err := NewSequenceNumber()
	assert.NoError(t, err)
	assert.NotEqual(t, SequenceNumber{}, sequenceNumber[0:2])
}
