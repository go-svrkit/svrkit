// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package pool

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArena_Simple(t *testing.T) {
	var a = NewArenaPool[int]()
	assert.Equal(t, a.Len(), 0)
	assert.Equal(t, a.Cap(), 1024)
	assert.Equal(t, a.Size(), 1024)

	var b = a.Alloc()
	assert.NotNil(t, b)
	assert.Equal(t, a.Len(), 1)
	assert.Equal(t, a.Cap(), 1023)
	assert.Equal(t, a.Size(), 1024)

	var c = a.AllocN(23)
	assert.NotNil(t, c)
	assert.Equal(t, len(c), 23)
	assert.Equal(t, a.Len(), 24)
	assert.Equal(t, a.Cap(), 1000)
	assert.Equal(t, a.Size(), 1024)

	var oldPtr = a.Ptr()
	var capacity = a.Cap()
	for i := 0; i < capacity; i++ {
		a.Alloc()
	}
	a.Alloc() // alloc new block
	assert.NotEqual(t, oldPtr, a.Ptr())
}

func arenaWorker(ctx context.Context, t *testing.T, pool *ArenaPool[bytes.Buffer]) {
	var count = 0
	var ticker = time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for count < 1000 {
		select {
		case <-ticker.C:
			count++
			var buf = pool.Alloc()
			assert.NotNil(t, buf)

		case <-ctx.Done():
			return
		}
	}
}

func TestArena_Concurrent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	var a = NewArenaPool[bytes.Buffer]()
	for i := 0; i < 10; i++ {
		go arenaWorker(ctx, t, a)
	}
	<-ctx.Done()
}
