// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package timers

import (
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/svrkit.v1/sched"
	"gopkg.in/svrkit.v1/zlog"
)

const (
	DefaultTimeoutCapacity = 1 << 16
)

// reserved action IDs
const (
	ActionBlackHole  = 0
	ActionExecFunc   = 1
	ActionExecRunner = 2
)

// make it easy to mock
var clockNow = time.Now

// TimeoutMsg 定时器到期消息
type TimeoutMsg struct {
	Owner    int64 // 定时器归属
	Deadline int64 // 定时器到期时间
	Action   int   // 定时器触发行为枚举
	Data     any   // 注意这个字段可能会序列化
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

type TimerMgr struct {
	guard       sync.Mutex
	running     atomic.Bool
	nexId       atomic.Int64
	timers      TimerScheduler
	timeouts    map[int64]*TimeoutMsg
	timeoutChan chan *TimeoutMsg
}

func NewTimerMgr() *TimerMgr {
	return &TimerMgr{
		timers:      NewTimerQueue(DefaultTimeoutCapacity),
		timeouts:    make(map[int64]*TimeoutMsg, 1024),
		timeoutChan: make(chan *TimeoutMsg, DefaultTimeoutCapacity),
	}
}

func (tm *TimerMgr) worker(ready chan struct{}) {
	zlog.Infof("timer worker thread started")
	defer zlog.Infof("timer worker thread exit")

	ready <- struct{}{}

	for tm.running.Load() {
		select {
		case tid := <-tm.timers.TimedOutChan():
			tm.guard.Lock()
			var msg = tm.timeouts[tid]
			delete(tm.timeouts, tid)
			tm.guard.Unlock()
			tm.timeoutChan <- msg
		}
	}
}

func (tm *TimerMgr) Shutdown() {
	if !tm.running.CompareAndSwap(true, false) {
		return
	}
	tm.timers.Shutdown()

	tm.guard.Lock()
	defer tm.guard.Unlock()

	tm.timers = nil
	close(tm.timeoutChan)
	tm.timeoutChan = nil
	tm.timeouts = nil
}

func (tm *TimerMgr) IsPending(timerId int64) bool {
	tm.guard.Lock()
	msg, found := tm.timeouts[timerId]
	tm.guard.Unlock()
	return found && msg != nil
}

func (tm *TimerMgr) TimedOutChan() <-chan *TimeoutMsg {
	return tm.timeoutChan
}

func (tm *TimerMgr) Cancel(timerId int64) bool {
	tm.guard.Lock()
	delete(tm.timeouts, timerId)
	tm.guard.Unlock()
	return tm.timers.CancelTimeout(timerId)
}

func (tm *TimerMgr) AddTimerAt(deadline int64, msg *TimeoutMsg) int64 {
	if tm.running.CompareAndSwap(false, true) {
		tm.timers.Start()
		var ready = make(chan struct{})
		go tm.worker(ready)
		<-ready
	}
	tm.guard.Lock()
	var id = tm.nexId.Add(1)
	tm.timers.AddTimeoutAt(id, deadline)
	tm.timeouts[id] = msg
	tm.guard.Unlock()
	return id
}

func (tm *TimerMgr) AddTimer(owner, durationMs int64, action int, data any) int64 {
	var deadline = clockNow().UnixMilli() + durationMs
	var msg = &TimeoutMsg{Owner: owner, Deadline: deadline, Action: action, Data: data}
	return tm.AddTimerAt(deadline, msg)
}

// RunAt 在`deadline`执行`action`
func (tm *TimerMgr) RunAt(deadlineMs int64, action func()) int64 {
	var msg = &TimeoutMsg{Action: ActionExecFunc, Data: action}
	return tm.AddTimerAt(deadlineMs, msg)
}

// RunAfter 在`durationMs`毫秒后执行`action`
func (tm *TimerMgr) RunAfter(durationMs int64, action func()) int64 {
	var deadline = clockNow().UnixMilli() + durationMs
	var msg = &TimeoutMsg{Action: ActionExecFunc, Data: action}
	return tm.AddTimerAt(deadline, msg)
}

// ScheduleAt 在`deadline`执行`runnable`
func (tm *TimerMgr) ScheduleAt(deadline int64, runnable sched.IRunner) int64 {
	var msg = &TimeoutMsg{Action: ActionExecRunner, Data: runnable}
	return tm.AddTimerAt(deadline, msg)
}

// Schedule 在`durationMs`毫秒后执行`runnable`
func (tm *TimerMgr) Schedule(durationMs int64, runnable sched.IRunner) int64 {
	var deadline = clockNow().UnixMilli() + durationMs
	return tm.ScheduleAt(deadline, runnable)
}

// TimerData 持久化定时器
type TimerData struct {
	Id       int64 `json:"id"`
	Deadline int64 `json:"deadline"`
	Owner    int64 `json:"owner"`
	Action   int   `json:"action"`
	Data     any   `json:"data,omitempty"`
}

type RawTimerData struct {
	Timers []TimerData `json:"timers"`
	NextId int64       `json:"next_id"`
}

func (tm *TimerMgr) Dump() *RawTimerData {
	tm.guard.Lock()
	defer tm.guard.Unlock()

	var info = &RawTimerData{
		NextId: tm.nexId.Add(100), // reserve a little space
		Timers: make([]TimerData, 0, len(tm.timeouts)),
	}
	for id, timeout := range tm.timeouts {
		if timeout != nil && timeout.Owner > 0 {
			var ti = TimerData{
				Id:       id,
				Deadline: timeout.Deadline,
				Owner:    timeout.Owner,
				Action:   timeout.Action,
				Data:     timeout.Data,
			}
			info.Timers = append(info.Timers, ti)
		}
	}
	return info
}

func (tm *TimerMgr) Restore(td *RawTimerData) {
	tm.guard.Lock()
	tm.nexId.Store(td.NextId)
	tm.guard.Unlock()

	for _, ti := range td.Timers {
		var msg = &TimeoutMsg{
			Owner:  ti.Owner,
			Action: ti.Action,
			Data:   ti.Data,
		}
		tm.AddTimerAt(ti.Deadline, msg)
	}
}
