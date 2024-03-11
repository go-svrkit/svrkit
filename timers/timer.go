// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package timers

import (
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/svrkit.v1/datetime"
	"gopkg.in/svrkit.v1/sched"
	"gopkg.in/svrkit.v1/zlog"
)

const (
	DefaultTimeoutCapacity = 1 << 16
)

// reserved action IDs
const (
	ActionExecFunc   = 1
	ActionExecRunner = 2
)

type TimeoutMsg struct {
	Owner    int64 // 定时器归属
	Deadline int64 //
	Arg      int64 // persisted data
	Action   int   //
	Data     any   // not persisted data
}

// TimerScheduler 定时器
type TimerScheduler interface {
	Size() int
	IsPending(tid int64) bool

	TimedOutChan() <-chan int64
	Range(func(id, deadline int64))

	AddTimeout(tid int64, delayMs int64)
	AddTimeoutAt(tid int64, deadline int64)
	CancelTimeout(tid int64) bool

	Start() error
	Shutdown()
}

var (
	gLock        sync.Mutex
	gRunning     atomic.Bool
	gTid         atomic.Int64
	gTimer       TimerScheduler
	gTimeouts    = make(map[int64]*TimeoutMsg, 1024)
	gTimeoutChan = make(chan *TimeoutMsg, DefaultTimeoutCapacity)
)

func init() {
	gTimer = NewTimerQueue(DefaultTimeoutCapacity)
}

func currentUnixNano() int64 {
	return datetime.NowNano()
}

func Init() {
	gLock.Lock()
	defer gLock.Unlock()

	if gTimer == nil {
		gTimer = NewTimerQueue(DefaultTimeoutCapacity)
	}
	if gTimeouts == nil {
		gTimeouts = make(map[int64]*TimeoutMsg, 1024)
	}
	if gTimeoutChan == nil {
		gTimeoutChan = make(chan *TimeoutMsg, DefaultTimeoutCapacity)
	}
}

func SetDefault(timer TimerScheduler) {
	gLock.Lock()
	gTimer = timer
	gLock.Unlock()
}

func IsPending(timerId int64) bool {
	gLock.Lock()
	msg, found := gTimeouts[timerId]
	gLock.Unlock()
	return found && msg != nil
}

func TimedOutChan() <-chan *TimeoutMsg {
	return gTimeoutChan
}

func Cancel(timerId int64) bool {
	gLock.Lock()
	delete(gTimeouts, timerId)
	gLock.Unlock()
	return gTimer.CancelTimeout(timerId)
}

func AddTimerAt(deadline int64, msg *TimeoutMsg) int64 {
	if gRunning.CompareAndSwap(false, true) {
		gTimer.Start()
		var ready = make(chan struct{})
		go gTimerWorker(ready)
		<-ready
	}
	gLock.Lock()
	var id = gTid.Add(1)
	gTimer.AddTimeoutAt(id, deadline)
	gTimeouts[id] = msg
	gLock.Unlock()
	return id
}

func AddTimer(owner, duration int64, action int, arg int64, data any) int64 {
	var deadline = currentUnixNano() + duration*int64(time.Millisecond)
	var msg = &TimeoutMsg{Owner: owner, Deadline: deadline, Action: action, Arg: arg, Data: data}
	return AddTimerAt(deadline, msg)
}

// RunAt 在`deadline`执行`action`
func RunAt(deadline int64, action func()) int64 {
	var msg = &TimeoutMsg{Action: ActionExecFunc, Data: action}
	return AddTimerAt(deadline, msg)
}

// RunAfter 在`durationMs`毫秒后执行`action`
func RunAfter(durationMs int64, action func()) int64 {
	var deadline = currentUnixNano() + durationMs*int64(time.Millisecond)
	var msg = &TimeoutMsg{Action: ActionExecFunc, Data: action}
	return AddTimerAt(deadline, msg)
}

// ScheduleAt 在`deadline`执行`runnable`
func ScheduleAt(deadline int64, runnable sched.IRunner) int64 {
	var msg = &TimeoutMsg{Action: ActionExecRunner, Data: runnable}
	return AddTimerAt(deadline, msg)
}

// Schedule 在`durationMs`毫秒后执行`runnable`
func Schedule(durationMs int64, runnable sched.IRunner) int64 {
	var deadline = currentUnixNano() + durationMs*int64(time.Millisecond)
	return ScheduleAt(deadline, runnable)
}

func gTimerWorker(ready chan struct{}) {
	zlog.Infof("timeout worker thread started")
	defer zlog.Infof("timeout worker thread exit")

	ready <- struct{}{}

	for gRunning.Load() {
		select {
		case tid := <-gTimer.TimedOutChan():
			gLock.Lock()
			var msg = gTimeouts[tid]
			delete(gTimeouts, tid)
			gLock.Unlock()
			gTimeoutChan <- msg
		}
	}
}

func Shutdown() {
	if !gRunning.CompareAndSwap(true, false) {
		return
	}
	gTimer.Shutdown()

	gLock.Lock()
	defer gLock.Unlock()

	gTimer = nil
	close(gTimeoutChan)
	gTimeoutChan = nil
}

func Preprocess(msg *TimeoutMsg) bool {
	switch msg.Action {
	case ActionExecRunner:
		var runner = msg.Data.(sched.IRunner)
		if runner != nil {
			if err := runner.Run(); err != nil {
				zlog.Errorf("run timeout msg %d %T: %v", msg.Action, runner, err)
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
