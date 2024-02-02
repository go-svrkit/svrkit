// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package timers

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	WHEEL_SIZE    = 512
	TICK_DURATION = 10
	TIME_UNIT     = 10
)

// A hashed wheel timer inspired by [Netty HashedWheelTimer]
// see https://github.com/netty/netty/blob/netty-4.1.106.Final/common/src/main/java/io/netty/util/HashedWheelTimer.java

type TimerWheel struct {
	done      chan struct{}
	wg        sync.WaitGroup //
	running   atomic.Int32
	guard     sync.Mutex              // 多线程
	ref       map[int64]*WheelTimeout //
	wheels    []WheelBucket           //
	C         chan int64              // 到期的定时器
	startedAt int64
	lastTime  int64
	ticks     int64
}

var _ TimerScheduler = (*TimerWheel)(nil)

func NewTimerWheel() *TimerWheel {
	return new(TimerWheel).Init(DefaultTimeoutCapacity)
}

func (w *TimerWheel) Init(capacity int) *TimerWheel {
	w.done = make(chan struct{})
	w.wheels = make([]WheelBucket, WHEEL_SIZE)
	w.ref = make(map[int64]*WheelTimeout, 1024)
	w.C = make(chan int64, capacity)

	var curMilliSec = currentUnixNano() / int64(time.Millisecond)
	w.startedAt = curMilliSec
	w.lastTime = curMilliSec
	return w
}

func (w *TimerWheel) Size() int {
	w.guard.Lock()
	var n = len(w.ref)
	w.guard.Unlock()
	return n
}

func (w *TimerWheel) IsPending(tid int64) bool {
	w.guard.Lock()
	var node = w.ref[tid]
	w.guard.Unlock()
	return node != nil
}

func (w *TimerWheel) TimedOutChan() <-chan int64 {
	return w.C
}

func (w *TimerWheel) Range(action func(id, deadline int64)) {
	w.guard.Lock()
	defer w.guard.Unlock()

	for _, node := range w.ref {
		action(node.id, node.deadline)
	}
}

func (w *TimerWheel) Start() error {
	if !w.running.CompareAndSwap(0, 1) {
		return nil
	}
	w.wg.Add(1)
	go w.worker()
	return nil
}

func (w *TimerWheel) AddTimeoutAt(tid int64, deadline int64) {
	w.guard.Lock()
	defer w.guard.Unlock()

	var timeout = NewWheelTimeout(tid, deadline)
	var calculated = (timeout.deadline - w.startedAt) / int64(TICK_DURATION)
	timeout.remainRounds = int(calculated-w.ticks) / WHEEL_SIZE
	var idx = w.ticks & (WHEEL_SIZE - 1)
	w.wheels[idx].AddTimeout(timeout)
	w.ref[tid] = timeout
}

func (w *TimerWheel) AddTimeout(tid int64, delayMs int64) {
	if delayMs < 0 {
		delayMs = 0
	}
	var deadline = currentUnixNano()/int64(time.Millisecond) + delayMs*int64(time.Millisecond)
	w.AddTimeoutAt(tid, deadline)
}

func (w *TimerWheel) CancelTimeout(tid int64) bool {
	w.guard.Lock()
	defer w.guard.Unlock()

	if timeout, found := w.ref[tid]; found {
		if timeout != nil {
			timeout.bucket.RemoveTimeout(timeout)
			timeout.bucket = nil
			timeout.prev = nil
			timeout.next = nil
		}
		delete(w.ref, tid)
		return true
	}
	return false
}

func (w *TimerWheel) Clear() {
	w.guard.Lock()
	defer w.guard.Unlock()

	for _, bucket := range w.wheels {
		var node = bucket.head
		for node != nil {
			var next = node.next
			delete(w.ref, node.id)
			node.bucket = nil
			node.prev = nil
			node.next = nil
			node = next
		}
		bucket.head = nil
		bucket.tail = nil
	}
}

func (w *TimerWheel) Shutdown() {
	if !w.running.CompareAndSwap(1, 0) {
		return
	}
	close(w.done)
	w.wg.Wait()
	w.Clear()
}

func (w *TimerWheel) tick() {
	var deadline = w.startedAt + int64(TICK_DURATION)*(w.ticks+1)
	var idx = w.ticks % (WHEEL_SIZE - 1)
	var bucket = w.wheels[idx]
	var expired = bucket.ExpireTimeouts(deadline)
	for _, timeout := range expired {
		delete(w.ref, timeout.id)
		w.C <- timeout.id
	}
	w.ticks++
}

func (w *TimerWheel) worker() {
	defer w.wg.Done()

	var ticker = time.NewTicker(TICK_DURATION)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			var ticks = (now.UnixNano() - w.lastTime) / int64(TIME_UNIT)
			if ticks > 0 {
				w.lastTime = now.UnixNano()
				for i := int64(0); i < ticks; i++ {
					w.tick()
				}
			}

		case <-w.done:
			return
		}
	}
}

type WheelTimeout struct {
	prev, next   *WheelTimeout
	bucket       *WheelBucket
	id           int64
	deadline     int64
	remainRounds int
}

func NewWheelTimeout(id int64, deadline int64) *WheelTimeout {
	return &WheelTimeout{id: id, deadline: deadline}
}

type WheelBucket struct {
	head, tail *WheelTimeout
}

// AddTimeout add `timeout` to this bucket
func (b *WheelBucket) AddTimeout(timeout *WheelTimeout) {
	timeout.bucket = b
	if b.head == nil {
		b.head = timeout
		b.tail = timeout
	} else {
		b.tail.next = timeout
		timeout.prev = b.tail
		b.tail = timeout
	}
}

// RemoveTimeout remove `timeout` from linked list and return next linked one
func (b *WheelBucket) RemoveTimeout(timeout *WheelTimeout) *WheelTimeout {
	var next = timeout.next
	if timeout.prev != nil {
		timeout.prev.next = next
	}
	if timeout.next != nil {
		timeout.next.prev = timeout.prev
	}
	if timeout == b.head {
		// if timeout is also the tail we need to adjust the entry too
		if timeout == b.tail {
			b.head = nil
			b.tail = nil
		} else {
			b.head = next
		}
	} else if timeout == b.tail {
		// if the timeout is the tail modify the tail to be the prev node.
		b.tail = timeout.prev
	}
	return next
}

func (b *WheelBucket) ExpireTimeouts(deadline int64) []*WheelTimeout {
	var expired []*WheelTimeout
	var timeout = b.head
	for timeout != nil {
		var next = timeout.next
		if timeout.remainRounds <= 0 {
			next = b.RemoveTimeout(timeout)
			if timeout.deadline > deadline {
				// The timeout was placed into a wrong slot. This should never happen.
				log.Printf("timeout.deadline > now %d/%d\n", timeout.deadline, deadline)
			}
		} else {
			timeout.remainRounds--
		}
		timeout = next
	}
	return expired
}
