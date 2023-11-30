package ethx

import (
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"testing"
	"time"
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
	assert.Equal(t, callFunc(h1, human1.Name)[0], "A1")
	assert.Equal(t, callFunc(&h1, human1.Name)[0], "A1")
	assert.Equal(t, callFunc(&h1, new(human1).Name)[0], "A1")
	assert.Equal(t, callFunc(&h1, "Name")[0], "A1")

	h2 := human2{"A2"}
	assert.Equal(t, callFunc(h2, human2.Name)[0], "A2")
	assert.Equal(t, callFunc(&h2, human2.Name)[0], "A2")
	assert.Equal(t, callFunc(&h2, new(human2).Name)[0], "A2")
	assert.Equal(t, callFunc(&h2, "Name")[0], "A2")
}

func Test_Ptr(t *testing.T) {
	x := get()
	log.Println(x)
	time.Sleep(2 * time.Second)
	log.Println(x)
}

func get() (y *big.Int) {
	y = BigInt(100)
	go func() {
		time.Sleep(1 * time.Second)
		*y = *BigInt(200)
	}()
	return y
}
