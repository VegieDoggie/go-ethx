package ethx

import (
	"context"
	"golang.org/x/exp/slices"
	"golang.org/x/time/rate"
	"math"
	"math/rand"
	"reflect"
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
	r.posi %= r.Len()
	t := r.list[r.posi]
	r.posi++
	return t
}

func (r *Iterator[T]) Add(item ...T) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	beforeLen := r.Len()
	r.list = append(r.list, item...)
	r.updateLimit(beforeLen)
}

func (r *Iterator[T]) Remove(item ...T) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	beforeLen := r.Len()
	for i := r.Len() - 1; i >= 0; i-- {
		for j := range item {
			if r.isEqual(r.list[i], item[j]) {
				slices.Delete(r.list, i, i+1)
			}
		}
	}
	r.updateLimit(beforeLen)
}

func (r *Iterator[T]) updateLimit(beforeLen int) {
	if curLen := r.Len(); beforeLen != curLen && curLen > 0 {
		if beforeLen > 0 {
			ratio := float64(curLen) / float64(beforeLen)
			r.limit.SetLimit(rate.Limit(float64(r.limit.Limit()) * ratio))
			r.limit.SetBurst(int(math.Round(float64(r.limit.Burst()) * ratio)))
		} else {
			if r.limit == nil {
				r.limit = rate.NewLimiter(rate.Limit(curLen), curLen)
			}
		}
	}
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

func (r *Iterator[T]) isEqual(a, b any) bool {
	if reflect.TypeOf(a).Kind() == reflect.Ptr || reflect.TypeOf(b).Kind() == reflect.Ptr {
		if a == b {
			return true
		}
	}
	return reflect.DeepEqual(a, b)
}
