package stream

import "sync"

type Pool[T any] struct{ pool sync.Pool }

func NewPool[T any](fn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return fn()
			},
		},
	}
}

func (p *Pool[T]) Put(v T) { p.pool.Put(v) }
func (p *Pool[T]) Get() T  { return p.pool.Get().(T) }
