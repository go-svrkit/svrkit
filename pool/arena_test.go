// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package pool

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArena_Simple(t *testing.T) {
	var a = NewArenaPool[bytes.Buffer](100)
	for i := 0; i < 102; i++ {
		var buf = a.Alloc()
		assert.NotNil(t, buf)
	}
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var a = NewArenaPool[bytes.Buffer](100)
	for i := 0; i < 10; i++ {
		go arenaWorker(ctx, t, a)
	}
	<-ctx.Done()
}
