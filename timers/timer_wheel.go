// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package timers

import (
	"log"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

const (
	WHEEL_SIZE    = 512                   // the size of the wheel
	TICK_DURATION = 50 * time.Millisecond // the duration between tick
)

// TimerWheel A hashed wheel timer inspired by [Netty HashedWheelTimer]
// see https://github.com/netty/netty/blob/netty-4.1.106.Final/common/src/main/java/io/netty/util/HashedWheelTimer.java
type TimerWheel struct {
	done         chan struct{}
	wg           sync.WaitGroup //
	running      atomic.Int32
	tickDuration time.Duration

	guard  sync.Mutex              // 多线程
	ref    map[int64]*WheelTimeout //
	wheels []*WheelBucket          //
	C      chan int64              // 到期的定时器

	startedAt int64
	lastTime  int64
	ticks     int64
}

var _ TimerScheduler = (*TimerWheel)(nil)

func NewTimerWheel(capacity int) *TimerWheel {
	return new(TimerWheel).Init(capacity, TICK_DURATION/2)
}

func (w *TimerWheel) Init(capacity int, tickDuration time.Duration) *TimerWheel {
	if tickDuration < time.Millisecond {
		tickDuration = time.Millisecond
	}
	w.tickDuration = tickDuration
	w.done = make(chan struct{})
	w.C = make(chan int64, capacity)
	w.ref = make(map[int64]*WheelTimeout, 1024)
	w.wheels = make([]*WheelBucket, WHEEL_SIZE)
	for i := 0; i < len(w.wheels); i++ {
		w.wheels[i] = &WheelBucket{timers: w, level: i + 1}
	}
	var nowNano = clockNow().UnixNano()
	w.startedAt = nowNano
	w.lastTime = nowNano
	//log.Printf("timer wheel start at %s\n", datetime.FormatNanoTime(nowNano))
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
	var ready = make(chan struct{})
	go w.worker(ready)
	<-ready
	return nil
}

func (w *TimerWheel) AddTimeoutAt(tid int64, deadline int64) {
	w.guard.Lock()
	defer w.guard.Unlock()

	var timeout = NewWheelTimeout(tid, deadline)
	var calculated = (timeout.deadline - w.startedAt) / int64(TICK_DURATION)
	timeout.remainRounds = (calculated - w.ticks) / WHEEL_SIZE
	var ticks = calculated
	if ticks < w.ticks {
		ticks = w.ticks
	}
	var idx = ticks & (WHEEL_SIZE - 1)
	var bucket = w.wheels[idx]
	bucket.AddTimeout(timeout)
	w.ref[tid] = timeout
}

func (w *TimerWheel) AddTimeout(tid int64, delayMs int64) {
	var deadline = clockNow().UnixMilli() + delayMs
	if delayMs > 0 && deadline < 0 {
		deadline = math.MaxInt64 // guard against overflow
	}
	w.AddTimeoutAt(tid, deadline)
}

func (w *TimerWheel) CancelTimeout(tid int64) bool {
	w.guard.Lock()
	defer w.guard.Unlock()

	if timeout, found := w.ref[tid]; found {
		if timeout != nil {
			timeout.bucket.Remove(timeout)
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
	close(w.C)
	clear(w.ref)
	w.ref = nil
	w.wheels = nil
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
	w.guard.Lock()
	defer w.guard.Unlock()

	var deadline = w.startedAt + int64(TICK_DURATION)*(w.ticks+1)
	var idx = w.ticks % (WHEEL_SIZE - 1)
	var bucket = w.wheels[idx]
	//log.Printf("tick %d update bucket=%d size=%d deadline=%s", w.ticks, idx, bucket.SlowSize(), datetime.FormatNanoTime(deadline))
	bucket.ExpireTimeouts(deadline)
	w.ticks++
}

func (w *TimerWheel) update() {
	var nowNano = clockNow().UnixNano()
	for w.lastTime+int64(TICK_DURATION) <= nowNano {
		w.lastTime += int64(TICK_DURATION)
		w.tick()
	}
}

func (w *TimerWheel) worker(ready chan struct{}) {
	defer w.wg.Done()

	var ticker = time.NewTicker(w.tickDuration)
	defer ticker.Stop()

	ready <- struct{}{}

	for w.running.Load() > 0 {
		select {
		case <-ticker.C:
			w.update()

		case <-w.done:
			return
		}
	}
}

type WheelTimeout struct {
	prev, next   *WheelTimeout // This will be used to chain timeouts in WheelBucket via a double-linked-list.
	bucket       *WheelBucket  // The bucket to which the timeout was added
	id           int64
	deadline     int64 // expired time in nanoseconds
	remainRounds int64
}

func NewWheelTimeout(id int64, deadline int64) *WheelTimeout {
	return &WheelTimeout{id: id, deadline: deadline}
}

type WheelBucket struct {
	head, tail *WheelTimeout
	timers     *TimerWheel
	level      int
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

func (b *WheelBucket) SlowSize() int {
	var count = 0
	var t = b.head
	for t != nil {
		count++
		t = t.next
	}
	return count
}

// Remove removes `timeout` from linked list and return next linked one
func (b *WheelBucket) Remove(timeout *WheelTimeout) *WheelTimeout {
	var next = timeout.next
	// remove timeout that was either processed or cancelled by updating the linked-list
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
	// unchain from this bucket to allow for GC.
	timeout.prev = nil
	timeout.next = nil
	timeout.bucket = nil
	return next
}

// ExpireTimeouts expire all timeouts for the given deadline.
func (b *WheelBucket) ExpireTimeouts(deadline int64) int {
	var count = 0
	var timeout = b.head

	// process all timeouts
	for timeout != nil {
		var next = timeout.next
		if timeout.remainRounds <= 0 {
			delete(b.timers.ref, timeout.id)
			next = b.Remove(timeout)
			if timeout.deadline <= deadline {
				count++
				b.timers.C <- timeout.id
				//log.Printf("timeout %d expired deadline=%s\n", timeout.id, datetime.FormatNanoTime(timeout.deadline))
			} else {
				// The timeout was placed into a wrong slot. This should never happen.
				log.Printf("timeout %d deadline greater than now %d > %d\n", timeout.id, timeout.deadline, deadline)
			}
		} else {
			timeout.remainRounds--
		}
		timeout = next
	}

	return count
}
