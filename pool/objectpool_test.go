// Copyright (c) 2023 Uber Technologies, Inc.

package pool

import (
	"runtime/debug"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type pooledValue[T any] struct {
	value T
}

func TestObjectPool_PutAndGet(t *testing.T) {
	p := NewObjectPoolWith(func() *int {
		return new(int)
	})

	value := 10
	p.Put(&value)

	retrievedValue := p.Get()
	require.Equal(t, value, *retrievedValue)
}

func TestObjectPool_MultiplePutAndGet(t *testing.T) {
	p := NewObjectPoolWith(func() *int {
		return new(int)
	})

	values := []int{10, 20, 30, 40, 50}
	for _, v := range values {
		p.Put(&v)
	}

	for range values {
		retrievedValue := p.Get()
		require.Contains(t, values, *retrievedValue)
	}
}

func TestObjectPool_ConcurrentPutAndGet(t *testing.T) {
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	p := NewObjectPoolWith(func() *int {
		return new(int)
	})

	var wg sync.WaitGroup
	values := []int{10, 20, 30, 40, 50}
	for _, v := range values {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			p.Put(&val)
		}(v)
	}

	wg.Wait()

	retrievedValue := p.Get()
	require.Contains(t, values, *retrievedValue)
}
