// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package timers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func pollExpiredTimeouts(timer TimerScheduler) (expired []int64) {
	for {
		select {
		case tid := <-timer.TimedOutChan():
			expired = append(expired, tid)
		default:
			return
		}
	}
}

func TestTimerQueue_Start(t *testing.T) {
	var timer = NewTimerQueue(64)
	defer timer.Shutdown()

	assert.True(t, timer.running.Load() == 0)

	timer.Start()

	assert.True(t, timer.running.Load() == 1)
}

func TestTimerQueue_IsPending(t *testing.T) {
	var timer = NewTimerQueue(64)
	defer timer.Shutdown()

	timer.AddTimeout(1, 0)
	timer.AddTimeout(2, 150)
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

func TestTimerQueue_AddTimeoutAt(t *testing.T) {
	var timer = NewTimerQueue(64)
	timer.Start()
	defer timer.Shutdown()

	var now = currentUnixNano()
	timer.AddTimeoutAt(1, now)
	timer.AddTimeoutAt(2, now+int64(150*time.Millisecond))
	timer.AddTimeoutAt(3, now+int64(500*time.Millisecond))
	assert.True(t, timer.IsPending(1))
	assert.True(t, timer.IsPending(2))
	assert.True(t, timer.IsPending(3))
	assert.Equal(t, 3, timer.Size())

	<-time.NewTimer(50 * time.Millisecond).C
	var timedOut = pollExpiredTimeouts(timer)
	assert.Equal(t, 1, len(timedOut))
	assert.Equal(t, 1, len(timedOut))
	assert.Equal(t, int64(1), timedOut[0])

	assert.Equal(t, 2, timer.Size())
	assert.False(t, timer.IsPending(1))
	assert.True(t, timer.IsPending(2))
	assert.True(t, timer.IsPending(3))

	<-time.NewTimer(500 * time.Millisecond).C
	timedOut = pollExpiredTimeouts(timer)
	assert.Equal(t, 2, len(timedOut))
	assert.Equal(t, int64(2), timedOut[0])
	assert.Equal(t, int64(3), timedOut[1])

	assert.Equal(t, 0, timer.Size())
	assert.False(t, timer.IsPending(1))
	assert.False(t, timer.IsPending(2))
	assert.False(t, timer.IsPending(3))
}

func TestTimerQueue_CancelTimeout(t *testing.T) {
	var timer = NewTimerQueue(64)
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

	<-time.NewTimer(50 * time.Millisecond).C
	assert.True(t, timer.IsPending(3))
	assert.False(t, timer.CancelTimeout(2))
	assert.True(t, timer.CancelTimeout(3))
	assert.Equal(t, 0, timer.Size())
}

func TestTimerQueue_Range(t *testing.T) {
	var timer = NewTimerQueue(64)

	var now = currentUnixNano()
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
