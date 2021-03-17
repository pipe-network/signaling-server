package values

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCombinedSequenceNumber_Increment(t *testing.T) {
	combinedSequenceNumber := CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(12),
		OverflowNumber: OverflowNumberFromInt(0),
	}
	actual, err := combinedSequenceNumber.Increment()
	assert.NoError(t, err)
	assert.Equal(t, CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(13),
		OverflowNumber: OverflowNumberFromInt(0),
	}, *actual)
}

func TestCombinedSequenceNumber_Increment_Overflow(t *testing.T) {
	combinedSequenceNumber := CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(4294967295),
		OverflowNumber: OverflowNumberFromInt(0),
	}
	actual, err := combinedSequenceNumber.Increment()
	assert.NoError(t, err)
	assert.Equal(t, CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(0),
		OverflowNumber: OverflowNumberFromInt(1),
	}, *actual)
}

func TestCombinedSequenceNumber_Increment_Overflow_Reached(t *testing.T) {
	combinedSequenceNumber := CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(4294967295),
		OverflowNumber: OverflowNumberFromInt(65535),
	}
	_, err := combinedSequenceNumber.Increment()
	assert.Error(t, err)
	assert.Equal(t, err, OverflowReached)
}

func TestCombinedSequenceNumber_Empty(t *testing.T) {
	combinedSequenceNumber := CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(4294967295),
		OverflowNumber: OverflowNumberFromInt(65535),
	}
	assert.False(t, combinedSequenceNumber.Empty())

	combinedSequenceNumber = CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(0),
		OverflowNumber: OverflowNumberFromInt(0),
	}
	assert.True(t, combinedSequenceNumber.Empty())
}

func TestCombinedSequenceNumber_Equal(t *testing.T) {
	assert.True(t, CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(4294967295),
		OverflowNumber: OverflowNumberFromInt(65535),
	}.Equal(CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(4294967295),
		OverflowNumber: OverflowNumberFromInt(65535),
	}))
	assert.False(t, CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(4294967295),
		OverflowNumber: OverflowNumberFromInt(65535),
	}.Equal(CombinedSequenceNumber{
		SequenceNumber: SequenceNumberFromInt(4294967295),
		OverflowNumber: OverflowNumberFromInt(0),
	}))
}
