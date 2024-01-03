// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package timers

import (
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/svrkit.v1/datetime"
	"gopkg.in/svrkit.v1/sched"
	"gopkg.in/svrkit.v1/slog"
)

const (
	DefaultTimeoutCapacity = 1 << 15
)

// reserved action IDs
const (
	ActionExecFunc   = 1
	ActionExecRunner = 2
)

var (
	guard       sync.Mutex
	running     atomic.Bool
	timeouts    = make(map[int64]*TimeoutMsg, DefaultQueueCapacity)
	timeoutChan = make(chan *TimeoutMsg, DefaultTimeoutCapacity)
	defTimer    = NewDefaultTimerQueue()
)

// TimerScheduler 定时器
type TimerScheduler interface {
	Start() error
	Shutdown()

	Size() int
	IsPending(timerId int64) bool

	TimedOutChan() <-chan int64

	AddTimer(delayMs int64) int64
	AddTimerAt(deadline int64) int64
	CancelTimer(timerId int64) bool
}

type TimeoutMsg struct {
	Owner  int64 // 定时器归属
	Action int32 //
	Param  int32 //
	Data   any   //
}

func currentUnixNano() int64 {
	return datetime.NowNano()
}

func AddTimerAt(deadline int64, msg *TimeoutMsg) int64 {
	if running.CompareAndSwap(false, true) {
		defTimer.Start()
		go timeoutWorker()
	}
	var tid = defTimer.AddTimerAt(deadline)
	guard.Lock()
	timeouts[tid] = msg
	guard.Unlock()
	return tid
}

func AddTimer(owner, duration int64, action, param int32, data any) int64 {
	var deadline = currentUnixNano() + duration*int64(time.Millisecond)
	var msg = &TimeoutMsg{Owner: owner, Action: action, Param: param, Data: data}
	return AddTimerAt(deadline, msg)
}

// RunAt 在`deadline`纳秒时间戳执行`action`
func RunAt(deadline int64, action func()) int64 {
	var msg = &TimeoutMsg{Action: int32(ActionExecFunc), Data: action}
	return AddTimerAt(deadline, msg)
}

// RunAfter 在`duration`毫秒后执行`action`
func RunAfter(duration int64, action func()) int64 {
	var deadline = currentUnixNano() + duration*int64(time.Millisecond)
	var msg = &TimeoutMsg{Action: int32(ActionExecFunc), Data: action}
	return AddTimerAt(deadline, msg)
}

// ScheduleAt 在`deadline`纳秒时间戳执行`runnable`
func ScheduleAt(deadline int64, runnable sched.IRunner) int64 {
	var msg = &TimeoutMsg{Action: int32(ActionExecRunner), Data: runnable}
	return AddTimerAt(deadline, msg)
}

// Schedule 在`duration`毫秒后执行`runnable`
func Schedule(duration int64, runnable sched.IRunner) int64 {
	var deadline = currentUnixNano() + duration*int64(time.Millisecond)
	return ScheduleAt(deadline, runnable)
}

func IsPending(timerId int64) bool {
	guard.Lock()
	msg, found := timeouts[timerId]
	guard.Unlock()
	return found && msg != nil
}

func Cancel(timerId int64) {
	guard.Lock()
	delete(timeouts, timerId)
	guard.Unlock()
	defTimer.CancelTimer(timerId)
}

func TimedOutChan() <-chan *TimeoutMsg {
	return timeoutChan
}

func Shutdown() {
	defTimer.Shutdown()
	defTimer = nil
}

func timeoutWorker() {
	defer slog.Infof("timer worker exit")

	for running.Load() {
		select {
		case id := <-defTimer.TimedOutChan():
			guard.Lock()
			var msg = timeouts[id]
			delete(timeouts, id)
			guard.Unlock()
			timeoutChan <- msg
		}
	}
}

func Preprocess(msg *TimeoutMsg) bool {
	switch msg.Action {
	case ActionExecRunner:
		var runner = msg.Data.(sched.IRunner)
		if runner != nil {
			if err := runner.Run(); err != nil {
				slog.Errorf("run timeout msg %d %T: %v", msg.Action, runner, err)
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
