// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package timers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimerWheel_Start(t *testing.T) {
	var timer = NewTimerWheel()
	defer timer.Shutdown()

	assert.True(t, timer.running.Load() == 0)

	timer.Start()

	assert.True(t, timer.running.Load() == 1)
}

func TestTimerWheel_IsPending(t *testing.T) {
	var timer = NewTimerWheel()
	defer timer.Shutdown()

	timer.AddTimeout(1, 0)
	timer.AddTimeout(2, 50)
	timer.AddTimeout(3, 500)
	assert.Equal(t, 3, timer.Size())

	assert.True(t, timer.IsPending(3))
	assert.False(t, timer.IsPending(4))

	timer.Start()

	<-time.NewTimer(400 * time.Millisecond).C

	assert.False(t, timer.IsPending(1))
	assert.False(t, timer.IsPending(2))
	assert.True(t, timer.IsPending(3))
	assert.False(t, timer.IsPending(4))
	assert.Equal(t, 1, timer.Size())
}

func TestTimerWheel_AddTimeoutAt(t *testing.T) {
	var timer = NewTimerWheel()
	timer.Start()
	defer timer.Shutdown()

	var now = currentUnixNano()
	timer.AddTimeoutAt(1, now)
	timer.AddTimeoutAt(2, now+int64(50*time.Millisecond))
	timer.AddTimeoutAt(3, now+int64(500*time.Millisecond))
	assert.True(t, timer.IsPending(1))
	assert.True(t, timer.IsPending(2))
	assert.True(t, timer.IsPending(3))
	assert.Equal(t, 3, timer.Size())

	<-time.NewTimer(10 * time.Millisecond).C
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
	var timer = NewTimerWheel()
	timer.Start()
	defer timer.Shutdown()

	timer.AddTimeout(1, 0)
	timer.AddTimeout(2, 50)
	timer.AddTimeout(3, 500)
	assert.True(t, timer.IsPending(1))
	assert.True(t, timer.IsPending(2))
	assert.True(t, timer.IsPending(3))
	assert.Equal(t, 3, timer.Size())

	assert.True(t, timer.CancelTimeout(2))
	var timeouts = pollExpiredTimeouts(timer)
	assert.Equal(t, 2, timer.Size())
	assert.Equal(t, 0, len(timeouts))

	<-time.NewTimer(20 * time.Millisecond).C
	assert.True(t, timer.IsPending(3))
	assert.False(t, timer.CancelTimeout(2))
	assert.True(t, timer.CancelTimeout(3))
}

func TestTimerWheel_Range(t *testing.T) {
	var timer = NewTimerWheel()

	var now = currentUnixNano()
	var d = map[int64]int64{
		1: now,
		2: now + 50,
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
