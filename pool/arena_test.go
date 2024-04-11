// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package pool

import (
	"bytes"
	"context"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestArena_Simple(t *testing.T) {
	var a = NewArenaPool[bytes.Buffer](100)
	var b = a.Alloc()
	assert.NotNil(t, b)

	a = NewArenaPool[bytes.Buffer](1)
	b = a.Alloc()
	assert.NotNil(t, b) // exhaust the block
	var idx = a.idx
	a.Free(b)
	assert.Equal(t, idx, a.idx)
	var ptr = unsafe.Pointer(&a.block[0])

	b = a.Alloc() // should allocate a new block
	assert.NotNil(t, b)
	assert.NotEqual(t, ptr, unsafe.Pointer(&a.block[0]))
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

	var a = NewArenaPool[bytes.Buffer](100)
	for i := 0; i < 10; i++ {
		go arenaWorker(ctx, t, a)
	}
	<-ctx.Done()
}
