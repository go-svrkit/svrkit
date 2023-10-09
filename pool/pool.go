// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package pool

import (
	"sync"
)

func createPool[T any]() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return new(T)
		},
	}
}

type ObjectPool[T any] struct {
	pool *sync.Pool
}

func NewObjectPool[T any]() *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: createPool[T](),
	}
}

func (a *ObjectPool[T]) Alloc() *T {
	return a.pool.Get().(*T)
}

func (a *ObjectPool[T]) Free(p *T) {
	a.pool.Put(p)
}
