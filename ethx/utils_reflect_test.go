package ethx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type human1 string

func (h human1) Name() string {
	return string(h)
}

type human2 struct {
	n string
}

func (h human2) Name() string {
	return h.n
}

func name(n string) string {
	return n
}

func Test_callFunc(t *testing.T) {
	assert.Equal(t, callFunc(name, "A1")[0], "A1")
	assert.Equal(t, callFunc(name, "A2")[0], "A2")
	assert.Equal(t, callFunc(name, "")[0], "")
}

func Test_callStructFunc(t *testing.T) {
	h1 := human1("A1")
	assert.Equal(t, callStructMethod(h1, human1.Name)[0], "A1")
	assert.Equal(t, callStructMethod(&h1, human1.Name)[0], "A1")
	assert.Equal(t, callStructMethod(&h1, new(human1).Name)[0], "A1")
	assert.Equal(t, callStructMethod(&h1, "Name")[0], "A1")

	h2 := human2{"A2"}
	assert.Equal(t, callStructMethod(h2, human2.Name)[0], "A2")
	assert.Equal(t, callStructMethod(&h2, human2.Name)[0], "A2")
	assert.Equal(t, callStructMethod(&h2, new(human2).Name)[0], "A2")
	assert.Equal(t, callStructMethod(&h2, "Name")[0], "A2")
}
