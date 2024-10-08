// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package timers

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	var logger = log.Default()
	logger.SetFlags(logger.Flags() | log.Lmicroseconds | log.Lshortfile)
}

func TestTimerWheel_Start(t *testing.T) {
	var timer = NewTimerWheel(64)
	defer timer.Shutdown()

	assert.True(t, timer.running.Load() == 0)

	timer.Start()

	assert.True(t, timer.running.Load() == 1)
}

func TestTimerWheel_IsPending(t *testing.T) {
	var timer = NewTimerWheel(64)

	defer timer.Shutdown()

	timer.AddTimeout(1, 0)
	timer.AddTimeout(2, 150)
	timer.AddTimeout(3, 500)
	assert.Equal(t, 3, timer.Size())

	assert.True(t, timer.IsPending(3))
	assert.False(t, timer.IsPending(4))

	timer.Start()
	<-time.NewTimer(300 * time.Millisecond).C

	assert.False(t, timer.IsPending(1))
	assert.False(t, timer.IsPending(2))
	assert.True(t, timer.IsPending(3))
	assert.False(t, timer.IsPending(4))
	assert.Equal(t, 1, timer.Size())
}

func TestTimerWheel_AddTimeoutAt(t *testing.T) {
	var timer = NewTimerWheel(64)
	timer.Start()
	defer timer.Shutdown()

	var now = clockNow().UnixMilli()
	timer.AddTimeoutAt(1, now)
	timer.AddTimeoutAt(2, now+int64(150*time.Millisecond))
	timer.AddTimeoutAt(3, now+int64(500*time.Millisecond))
	assert.True(t, timer.IsPending(1))
	assert.True(t, timer.IsPending(2))
	assert.True(t, timer.IsPending(3))
	assert.Equal(t, 3, timer.Size())

	<-time.NewTimer(100 * time.Millisecond).C
	var timeouts = pollExpiredTimeouts(timer)
	assert.Equal(t, 1, len(timeouts))
	assert.Equal(t, int64(1), timeouts[0])

	assert.Equal(t, 2, timer.Size())
	assert.False(t, timer.IsPending(1))
	assert.True(t, timer.IsPending(2))
	assert.True(t, timer.IsPending(3))

	<-time.NewTimer(500 * time.Millisecond).C
	timeouts = pollExpiredTimeouts(timer)
	assert.Equal(t, 2, len(timeouts))
	assert.Equal(t, int64(2), timeouts[0])
	assert.Equal(t, int64(3), timeouts[1])

	assert.Equal(t, 0, timer.Size())
	assert.False(t, timer.IsPending(1))
	assert.False(t, timer.IsPending(2))
	assert.False(t, timer.IsPending(3))
}

func TestTimerWheel_CancelTimeout(t *testing.T) {
	var timer = NewTimerWheel(64)
	timer.Start()
	defer timer.Shutdown()

	timer.AddTimeout(1, 0)
	timer.AddTimeout(2, 150)
	timer.AddTimeout(3, 500)
	assert.True(t, timer.IsPending(1))
	assert.True(t, timer.IsPending(2))
	assert.True(t, timer.IsPending(3))
	assert.Equal(t, 3, timer.Size())

	assert.True(t, timer.CancelTimeout(2))
	var timeouts = pollExpiredTimeouts(timer)
	assert.Equal(t, 2, timer.Size())
	assert.Equal(t, 0, len(timeouts))

	<-time.NewTimer(100 * time.Millisecond).C
	assert.True(t, timer.IsPending(3))
	assert.False(t, timer.CancelTimeout(2))
	assert.True(t, timer.CancelTimeout(3))
}

func TestTimerWheel_Range(t *testing.T) {
	var timer = NewTimerWheel(64)
	defer timer.Shutdown()

	var now = clockNow().UnixMilli()
	var d = map[int64]int64{
		1: now,
		2: now + 150,
		3: now + 500,
	}
	for tid, delay := range d {
		timer.AddTimeoutAt(tid, delay)
	}

	assert.Equal(t, 3, timer.Size())

	var d2 = map[int64]int64{}
	timer.Range(func(id, deadline int64) {
		d2[id] = deadline
	})
	assert.Equal(t, len(d2), len(d))

	for k := range d {
		assert.Equal(t, d[k], d2[k])
	}
}

// wheel timer should be First-In First-Out
func TestTimerWheel_FIFO(t *testing.T) {
	var timer = NewTimerWheel(64)
	timer.Start()
	defer timer.Shutdown()

	for i := 1; i <= 10; i++ {
		timer.AddTimeout(int64(i), 50)
	}

	<-time.NewTimer(250 * time.Millisecond).C

	var timeouts = pollExpiredTimeouts(timer)
	assert.Equal(t, 10, len(timeouts))
	for i := 1; i <= 10; i++ {
		assert.Equal(t, int64(i), timeouts[i-1])
	}
}
