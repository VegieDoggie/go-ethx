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
	limit *rate.Limiter
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
		limit: defaultLimiter,
		mutex: sync.Mutex{},
	}
}

func (r *Iterator[T]) WaitNext() T {
	_ = r.limit.Wait(r.ctx)
	return r.UnwaitNext()
}

func (r *Iterator[T]) UnwaitNext() T {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	t := r.list[r.posi]
	r.posi = (r.posi + 1) % len(r.list)
	return t
}

func (r *Iterator[T]) Add(item ...T) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.list = append(r.list, item...)
}

func (r *Iterator[T]) Len() int {
	return len(r.list)
}

func (r *Iterator[T]) Limiter() *rate.Limiter {
	return r.limit
}

func (r *Iterator[T]) All() []T {
	return r.list
}

func (r *Iterator[T]) NewCall(maxTry ...int) func(f any, args ...any) []any {
	n := 1
	if len(maxTry) > 0 && maxTry[0] > n {
		n = maxTry[0]
	}
	return func(f any, args ...any) (rets []any) {
		for i := 0; i < n; i++ {
			rets = callStructFunc(r.WaitNext(), f, args...)
			if len(rets) > 0 {
				if _, ok := rets[len(rets)-1].(error); ok {
					continue
				}
			}
			break
		}
		return rets
	}
}

func (r *Iterator[T]) Shuffle() *Iterator[T] {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rd.Shuffle(r.Len(), func(a, b int) {
		r.list[a], r.list[b] = r.list[b], r.list[a]
	})
	return r
}
