package ethx

import (
	"context"
	"golang.org/x/time/rate"
	"math/rand"
	"sync"
	"time"
)

type Iterator[T any] struct {
	list  []T
	posi  int
	ctx   context.Context
	mutex sync.Mutex
	Limit *rate.Limiter
}

func NewIterator[T any](list []T, limiter ...*rate.Limiter) *Iterator[T] {
	var defaultLimiter *rate.Limiter
	if len(limiter) > 0 {
		defaultLimiter = limiter[0]
	} else {
		defaultLimiter = rate.NewLimiter(rate.Limit(len(list)), len(list))
	}
	return &Iterator[T]{
		ctx:   context.Background(),
		list:  list,
		Limit: defaultLimiter,
		mutex: sync.Mutex{},
	}
}

func (r *Iterator[T]) Next() T {
	r.mutex.Lock()
	t := r.list[r.posi]
	r.posi = (r.posi + 1) % len(r.list)
	r.mutex.Unlock()
	return t
}

func (r *Iterator[T]) WaitNext() T {
	_ = r.Limit.Wait(r.ctx)
	return r.Next()
}

func (r *Iterator[T]) Len() int {
	return len(r.list)
}

func (r *Iterator[T]) Limiter() *rate.Limiter {
	return r.Limit
}

func (r *Iterator[T]) All() []T {
	return r.list
}

func (r *Iterator[T]) Call(c func(t T) bool, isWait ...bool) {
	if len(isWait) > 0 && isWait[0] {
		for !c(r.WaitNext()) {
		}
	} else {
		for !c(r.Next()) {
		}
	}
}

func (r *Iterator[T]) Shuffle() *Iterator[T] {
	r.mutex.Lock()
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rd.Shuffle(r.Len(), func(a, b int) {
		r.list[a], r.list[b] = r.list[b], r.list[a]
	})
	r.mutex.Unlock()
	return r
}
