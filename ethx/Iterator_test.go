package ethx

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func Test_IsEqual(t *testing.T) {
	type Z struct {
		a int
		b uint
	}
	x, y, z1, z2 := 1, 1, Z{1, 2}, Z{1, 2}
	data := []struct {
		a      any
		b      any
		except bool
	}{
		{x, y, true},
		{&x, y, false},
		{&x, &y, true},
		{x, &y, false},
		{x, x, true},
		{&x, &x, true},
		{int64(x), int64(y), true},
		{uint64(x), uint64(y), true},
		{float32(x), float32(y), true},
		{strconv.Itoa(x), strconv.Itoa(x), true},
		{strconv.Itoa(x), strconv.Itoa(y), true},
		{true, true, true},
		{false, false, true},
		{[]int{x, x, x}, []int{y, y, y}, true},
		{[]int{x, x, x}, x, false},
		{&x, []int{x, x, x}, false},
		//{Test_IsEqual, Test_IsEqual, true},
		{Z{1, 2}, Z{1, 2}, true},
		{z1, Z{1, 2}, true},
		{z1, z2, true},
		{&z1, &z2, true},
	}
	for i, d := range data {
		assert.Equal(t, d.except, new(Iterator[any]).isEqual(d.a, d.b), i, d)
	}
}

func Test_Iterator1(t *testing.T) {
	humans := []string{"A1", "B2", "C3", "D4", "E5"}
	iterator := NewIterator[string](humans)
	assert.Equal(t, len(humans), iterator.Len())
	iterator.Add("F6", "G7")
	assert.Equal(t, len(humans)+2, iterator.Len())
	for i := 0; i < 20; i++ {
		iterator.UnwaitNext()
	}
	iterator.Shuffle()
}

func Test_Iterator(t *testing.T) {
	it := new(Iterator[any])
	for i := 0; i < 10; i++ {
		for j := 0; j < 3; j++ {
			it.Remove(strconv.Itoa(rand.Int() % 256))
		}
		n1 := rand.Int() % 256
		for j := 0; j < n1; j++ {
			it.Add(strconv.Itoa(rand.Int() % 256))
		}
		n2 := rand.Int() % 256
		for j := 0; j < n2; j++ {
			it.Remove(strconv.Itoa(rand.Int() % 256))
		}
		n3 := rand.Int() % 256
		for j := 0; j < n3; j++ {
			it.UnwaitNext()
		}
		n4 := rand.Int() % 32
		for j := 0; j < n4; j++ {
			it.Shuffle()
		}
		for j := 0; j < n2; j++ {
			it.Remove(strconv.Itoa(rand.Int() % 256))
			it.UnwaitNext()
		}
		for j := 0; j < n2; j++ {
			it.Remove(strconv.Itoa(rand.Int() % 256))
			it.UnwaitNext()
		}
	}
}
