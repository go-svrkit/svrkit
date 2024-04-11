// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package timers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/sched"
)

func pollExpiredTimeoutMsg() (expired []*TimeoutMsg) {
	for {
		select {
		case msg := <-gTimeoutChan:
			expired = append(expired, msg)
		default:
			return
		}
	}
}

func pollTimeoutRunners(ctx context.Context) {
	for {
		select {
		case msg := <-gTimeoutChan:
			if msg != nil {
				Preprocess(msg)
			}
		case <-ctx.Done():
			return
		}
	}
}

func TestIsPending(t *testing.T) {
	Init()
	defer Shutdown()

	var tids []int64
	for i := 0; i < 10; i++ {
		var tid = AddTimer(int64(i), 100, 0, 0, 0)
		tids = append(tids, tid)
	}
	for _, tid := range tids {
		assert.True(t, IsPending(tid))
	}

	<-time.NewTimer(150 * time.Millisecond).C

	for _, tid := range tids {
		assert.False(t, IsPending(tid))
	}
}

func TestCancel(t *testing.T) {
	Init()
	defer Shutdown()

	var tids []int64
	for i := 0; i < 10; i++ {
		var tid = AddTimer(int64(i), 100, 0, 0, 0)
		tids = append(tids, tid)
	}
	for _, tid := range tids {
		assert.True(t, IsPending(tid))
	}

	for _, tid := range tids {
		assert.True(t, Cancel(tid))
	}

	for _, tid := range tids {
		assert.False(t, IsPending(tid))
	}
}

func TestAddTimer(t *testing.T) {
	Init()
	defer Shutdown()

	var tids = map[int]int64{}
	for i := 0; i < 10; i++ {
		var tid = AddTimer(int64(i), 100, i, 0, 0)
		tids[i] = tid
	}

	<-time.NewTimer(120 * time.Millisecond).C

	for _, tid := range tids {
		assert.False(t, IsPending(tid))
	}
	var timeouts = pollExpiredTimeoutMsg()
	assert.Equal(t, len(tids), len(timeouts))
	for _, msg := range timeouts {
		tid, found := tids[msg.Action]
		assert.True(t, found)
		assert.False(t, IsPending(tid))
	}
}

func TestRunAt(t *testing.T) {
	Init()
	defer Shutdown()

	var firedAt int64
	var startAt = currentUnixNano()
	var deadline = startAt + 100*int64(time.Millisecond)
	var tid = RunAt(deadline, func() {
		firedAt = currentUnixNano()
	})

	assert.True(t, IsPending(tid))

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	go pollTimeoutRunners(ctx)

	<-ctx.Done()
	assert.False(t, IsPending(tid))

	var endAt = currentUnixNano()
	assert.True(t, firedAt > 0)
	assert.True(t, firedAt > startAt)
	assert.True(t, firedAt < endAt)
}

func TestRunAfter(t *testing.T) {
	Init()
	defer Shutdown()

	var firedAt int64
	var startAt = currentUnixNano()
	var tid = RunAfter(100, func() {
		firedAt = currentUnixNano()
	})

	assert.True(t, IsPending(tid))

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	go pollTimeoutRunners(ctx)

	<-ctx.Done()
	assert.False(t, IsPending(tid))

	var endAt = currentUnixNano()
	assert.True(t, firedAt > 0)
	assert.True(t, firedAt > startAt)
	assert.True(t, firedAt < endAt)
}

func TestScheduleAt(t *testing.T) {
	Init()
	defer Shutdown()

	var firedAt int64
	var run = sched.NewRunnable(func() error {
		firedAt = currentUnixNano()
		return nil
	})

	var startAt = currentUnixNano()
	var deadline = startAt + 100*int64(time.Millisecond)
	var tid = ScheduleAt(deadline, run)

	assert.True(t, IsPending(tid))

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	go pollTimeoutRunners(ctx)

	<-ctx.Done()
	assert.False(t, IsPending(tid))

	var endAt = currentUnixNano()
	assert.True(t, firedAt > 0)
	assert.True(t, firedAt > startAt)
	assert.True(t, firedAt < endAt)
}

func TestSchedule(t *testing.T) {
	Init()
	defer Shutdown()

	var firedAt int64
	var run = sched.NewRunnable(func() error {
		firedAt = currentUnixNano()
		return nil
	})

	var startAt = currentUnixNano()
	var tid = Schedule(50, run)

	assert.True(t, IsPending(tid))

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	go pollTimeoutRunners(ctx)

	<-ctx.Done()
	assert.False(t, IsPending(tid))

	var endAt = currentUnixNano()
	assert.True(t, firedAt > 0)
	assert.True(t, firedAt > startAt)
	assert.True(t, firedAt < endAt)
}
