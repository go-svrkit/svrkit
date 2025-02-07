// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package timers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gopkg.in/svrkit.v1/qlog"
	"gopkg.in/svrkit.v1/sched"
)

func preprocess(msg *TimeoutMsg) bool {
	switch msg.Action {
	case ActionExecRunner:
		var runner = msg.Data.(sched.IRunner)
		if runner != nil {
			if err := runner.Run(); err != nil {
				qlog.Errorf("run timeout msg %d %T: %v", msg.Action, runner, err)
			}
		}
		return true

	case ActionExecFunc:
		var action = msg.Data.(func())
		if action != nil {
			action()
		}
		return true
	}
	return false
}

func pollExpiredTimeoutMsg(mgr *TimerMgr) (expired []*TimeoutMsg) {
	for {
		select {
		case msg := <-mgr.timeoutChan:
			expired = append(expired, msg)
		default:
			return
		}
	}
}

func pollTimeoutRunners(ctx context.Context, mgr *TimerMgr) {
	for {
		select {
		case msg := <-mgr.timeoutChan:
			if msg != nil {
				preprocess(msg)
			}
		case <-ctx.Done():
			return
		}
	}
}

func TestIsPending(t *testing.T) {
	var mgr = NewTimerMgr()
	defer mgr.Shutdown()

	var tids []int64
	for i := 0; i < 10; i++ {
		var tid = mgr.AddTimer(int64(i), 100, ActionBlackHole, 0)
		tids = append(tids, tid)
	}
	for _, tid := range tids {
		assert.True(t, mgr.IsPending(tid))
	}

	<-time.NewTimer(150 * time.Millisecond).C

	for _, tid := range tids {
		assert.False(t, mgr.IsPending(tid))
	}
}

func TestCancel(t *testing.T) {
	var mgr = NewTimerMgr()
	defer mgr.Shutdown()

	var tids []int64
	for i := 0; i < 10; i++ {
		var tid = mgr.AddTimer(int64(i), 100, ActionBlackHole, 0)
		tids = append(tids, tid)
	}
	for _, tid := range tids {
		assert.True(t, mgr.IsPending(tid))
	}

	for _, tid := range tids {
		assert.True(t, mgr.Cancel(tid))
	}

	for _, tid := range tids {
		assert.False(t, mgr.IsPending(tid))
	}
}

func TestAddTimer(t *testing.T) {
	var mgr = NewTimerMgr()
	defer mgr.Shutdown()

	var tids = map[int]int64{}
	for i := 0; i < 10; i++ {
		var tid = mgr.AddTimer(int64(i), 100, i, 0)
		tids[i] = tid
	}

	<-time.NewTimer(120 * time.Millisecond).C

	for _, tid := range tids {
		assert.False(t, mgr.IsPending(tid))
	}
	var timeouts = pollExpiredTimeoutMsg(mgr)
	assert.Equal(t, len(tids), len(timeouts))
	for _, msg := range timeouts {
		tid, found := tids[msg.Action]
		assert.True(t, found)
		assert.False(t, mgr.IsPending(tid))
	}
}

func TestRunAt(t *testing.T) {
	var mgr = NewTimerMgr()
	defer mgr.Shutdown()

	var firedAt int64
	var startAt = clockNow().UnixMilli()
	var tid = mgr.RunAt(startAt+100, func() {
		firedAt = clockNow().UnixMilli()
	})

	assert.True(t, mgr.IsPending(tid))

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	go pollTimeoutRunners(ctx, mgr)

	<-ctx.Done()
	assert.False(t, mgr.IsPending(tid))

	var endAt = clockNow().UnixMilli()
	assert.True(t, firedAt > 0)
	assert.True(t, firedAt > startAt)
	assert.True(t, firedAt < endAt)
}

func TestRunAfter(t *testing.T) {
	var mgr = NewTimerMgr()
	defer mgr.Shutdown()

	var firedAt int64
	var startAt = clockNow().UnixMilli()
	var tid = mgr.RunAfter(100, func() {
		firedAt = clockNow().UnixMilli()
	})

	assert.True(t, mgr.IsPending(tid))

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	go pollTimeoutRunners(ctx, mgr)

	<-ctx.Done()
	assert.False(t, mgr.IsPending(tid))

	var endAt = clockNow().UnixMilli()
	assert.True(t, firedAt > 0)
	assert.True(t, firedAt > startAt)
	assert.True(t, firedAt < endAt)
}

func TestScheduleAt(t *testing.T) {
	var mgr = NewTimerMgr()
	defer mgr.Shutdown()

	var firedAt int64
	var run = sched.NewRunnable(func() error {
		firedAt = clockNow().UnixMilli()
		return nil
	})

	var startAt = clockNow().UnixMilli()
	var tid = mgr.ScheduleAt(startAt+100, run)

	assert.True(t, mgr.IsPending(tid))

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	go pollTimeoutRunners(ctx, mgr)

	<-ctx.Done()
	assert.False(t, mgr.IsPending(tid))

	var endAt = clockNow().UnixMilli()
	assert.True(t, firedAt > 0)
	assert.True(t, firedAt > startAt)
	assert.True(t, firedAt < endAt)
}

func TestSchedule(t *testing.T) {
	var mgr = NewTimerMgr()
	defer mgr.Shutdown()

	var firedAt int64
	var run = sched.NewRunnable(func() error {
		firedAt = clockNow().UnixMilli()
		return nil
	})

	var startAt = clockNow().UnixMilli()
	var tid = mgr.Schedule(50, run)

	assert.True(t, mgr.IsPending(tid))

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	go pollTimeoutRunners(ctx, mgr)

	<-ctx.Done()
	assert.False(t, mgr.IsPending(tid))

	var endAt = clockNow().UnixMilli()
	assert.True(t, firedAt > 0)
	assert.True(t, firedAt > startAt)
	assert.True(t, firedAt < endAt)
}
