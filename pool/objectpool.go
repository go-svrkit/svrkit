// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package pool

import (
	"sync"
)

// ObjectPool is a generic wrapper around [sync.Pool] to provide strongly-typed object pooling.
// all internal pool use must take care to only store pointer types.
type ObjectPool[T any] struct {
	pool sync.Pool
}

func NewObjectPool[T any]() *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: sync.Pool{
			New: func() any {
				return new(T)
			},
		},
	}
}

func NewObjectPoolWith[T any](creator func() *T) *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: sync.Pool{
			New: func() any {
				return creator()
			},
		},
	}
}

func (a *ObjectPool[T]) Put(p *T) {
	if p != nil {
		a.pool.Put(p)
	}
}

func (a *ObjectPool[T]) Get() *T {
	var p = a.pool.Get().(*T)
	if p == nil {
		return new(T)
	}
	var zero T
	*p = zero
	return p
}
