package ethx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Iterator(t *testing.T) {
	humans := []string{"A1", "B2", "C3", "D4", "E5"}
	iterator := NewIterator[string](humans)
	assert.Equal(t, len(humans), iterator.Len())
	iterator.Add("F6", "G7")
	assert.Equal(t, len(humans)+2, iterator.Len())
	for i := 0; i < 20; i++ {
		iterator.UnwaitNext()
		iterator.WaitNext()
	}
	iterator.Shuffle()
}
