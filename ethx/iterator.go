package ethx

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/time/rate"
	"log"
	"math/rand"
	"sync"
	"time"
)

// NewClientIterator creates a clientIterator.
func NewClientIterator(rpcList []string, limiter ...*rate.Limiter) (clientIterator *Iterator[*ethclient.Client]) {
	return NewClientIteratorWithWeight(rpcList, make([]int, len(rpcList)), limiter...)
}

// NewClientIteratorWithWeight creates a clientIterator with wights.
func NewClientIteratorWithWeight(rpcList []string, weightList []int, limiter ...*rate.Limiter) *Iterator[*ethclient.Client] {
	if len(rpcList) != len(weightList) {
		log.Panicln("len(rpcList) != len(weightList)")
	}
	var reliableClients []*ethclient.Client
	for i, rpc := range rpcList {
		client, err := ethclient.Dial(rpc)
		if err == nil {
			blockNumber, err := client.BlockNumber(context.TODO())
			if err == nil && blockNumber > 0 {
				for j := 1; j < weightList[i]; j++ {
					reliableClients = append(reliableClients, client)
				}
				reliableClients = append(reliableClients, client)
			}
		}
	}
	if len(reliableClients) == 0 {
		log.Panicln("No reliable rpc nodes connection!")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(reliableClients), func(i, j int) {
		reliableClients[i], reliableClients[j] = reliableClients[j], reliableClients[i]
	})
	return NewIterator[*ethclient.Client](reliableClients, limiter...)
}

type Iterator[T any] struct {
	list    []T
	posi    int
	ctx     context.Context
	mutex   sync.Mutex
	limiter *rate.Limiter
}

func NewIterator[T any](list []T, limiter ...*rate.Limiter) *Iterator[T] {
	var defaultLimiter *rate.Limiter
	if len(limiter) > 0 {
		defaultLimiter = limiter[0]
	} else {
		defaultLimiter = rate.NewLimiter(rate.Every(time.Duration(len(list))*time.Second), len(list))
	}
	return &Iterator[T]{
		ctx:     context.Background(),
		list:    list,
		limiter: defaultLimiter,
		mutex:   sync.Mutex{},
	}
}

func (i *Iterator[T]) Next() T {
	i.mutex.Lock()
	t := i.list[i.posi]
	i.posi = (i.posi + 1) % len(i.list)
	i.mutex.Unlock()
	return t
}

func (i *Iterator[T]) WaitNext() T {
	_ = i.limiter.Wait(i.ctx)
	return i.Next()
}

func (i *Iterator[T]) Len() int {
	return len(i.list)
}

func (i *Iterator[T]) Limiter() *rate.Limiter {
	return i.limiter
}

func (i *Iterator[T]) All() []T {
	return i.list
}

func (i *Iterator[T]) Call(c func(t T) bool, isWait ...bool) {
	if len(isWait) > 0 && isWait[0] {
		for !c(i.WaitNext()) {
		}
	} else {
		for !c(i.Next()) {
		}
	}
}
