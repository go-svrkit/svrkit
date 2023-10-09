// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package timers

import (
	"container/heap"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TickInterval         = 5 * time.Millisecond // 5ms
	DefaultQueueCapacity = 1 << 10
)

// TimerQueue 最小堆实现的定时器
// 标准库(src/runtime/time.go)用四叉堆实现的`time.Timer`已经可以满足大多数定时需求；
// 这个实现主要是做用户层的timer，增加自定义功能支持、减少runtime的压力；
type TimerQueue struct {
	done    chan struct{}
	wg      sync.WaitGroup //
	running atomic.Int32

	guard  sync.Mutex           // 多线程
	nextId atomic.Int64         // id生成
	refer  map[int64]*timerNode // O(1)查找
	timers timerHeap            //
	C      chan int64           // 到期的定时器
}

func NewDefaultTimerQueue() *TimerQueue {
	var timer = NewTimerQueue(DefaultTimeoutCapacity)
	if err := timer.Start(); err != nil {
		log.Printf("start default timer %v", err)
	}
	return timer
}

func NewTimerQueue(capacity int) *TimerQueue {
	if capacity <= 0 {
		capacity = DefaultTimeoutCapacity
	}
	return new(TimerQueue).Init(capacity)
}

func (s *TimerQueue) Init(capacity int) *TimerQueue {
	s.done = make(chan struct{})
	s.timers = make(timerHeap, 0, DefaultQueueCapacity)
	s.refer = make(map[int64]*timerNode, DefaultQueueCapacity)
	s.C = make(chan int64, capacity)
	return s
}

func (s *TimerQueue) Size() int {
	s.guard.Lock()
	var n = len(s.refer)
	s.guard.Unlock()
	return n
}

func (s *TimerQueue) NextID() int64 {
	return s.nextId.Add(1)
}

func (s *TimerQueue) TimedOutChan() <-chan int64 {
	return s.C
}

func (s *TimerQueue) IsPending(id int64) bool {
	s.guard.Lock()
	var node = s.refer[id]
	s.guard.Unlock()
	return node != nil
}

// Start starts the background thread explicitly
func (s *TimerQueue) Start() error {
	if !s.running.CompareAndSwap(0, 1) {
		var ready = make(chan struct{}, 1)
		s.wg.Add(1)
		go s.worker(ready)
		<-ready
		return nil
	}
	return nil
}

func (s *TimerQueue) Shutdown() {
	if !s.running.CompareAndSwap(1, 0) {
		return
	}

	close(s.done)
	s.wg.Wait()

	s.guard.Lock()
	close(s.C)
	s.C = nil
	s.refer = nil
	s.timers = nil
	s.guard.Unlock()
}

// AddTimerAt 创建在指定时机触发的定时器 `deadline`使用纳秒
func (s *TimerQueue) AddTimerAt(deadline int64) int64 {
	s.guard.Lock()
	defer s.guard.Unlock()

	var tid = s.NextID()
	var node = newTimerNode(tid, deadline)
	s.refer[tid] = node
	heap.Push(&s.timers, node)

	return tid
}

// AddTimer 创建一个定时器，在`delayMs`毫秒后过期
func (s *TimerQueue) AddTimer(delayMs int64) int64 {
	if delayMs < 0 {
		delayMs = 0
	}
	var deadline = currentUnixNano() + delayMs*int64(time.Millisecond)
	return s.AddTimerAt(deadline)
}

// CancelTimer 取消一个timer
func (s *TimerQueue) CancelTimer(id int64) bool {
	s.guard.Lock()
	defer s.guard.Unlock()

	if node, found := s.refer[id]; found {
		heap.Remove(&s.timers, node.index)
		delete(s.refer, id)
		return true
	}
	return false
}

func (s *TimerQueue) RangeTimers(action func(node *timerNode)) {
	s.guard.Lock()
	defer s.guard.Unlock()

	for _, node := range s.refer {
		action(node)
	}
}

func (s *TimerQueue) worker(ready chan struct{}) {
	defer s.wg.Done()

	var ticker = time.NewTicker(TickInterval)
	defer ticker.Stop()

	ready <- struct{}{}

	for {
		select {
		case now := <-ticker.C:
			s.trigger(now.UnixNano())

		case <-s.done:
			return
		}
	}
}

// 返回触发的timer列表
func (s *TimerQueue) trigger(now int64) {
	s.guard.Lock()
	defer s.guard.Unlock()

	for len(s.timers) > 0 {
		var node = s.timers[0] // peek first item of heap
		if now < node.deadline {
			break // no new timer expired
		}

		heap.Pop(&s.timers)
		delete(s.refer, node.id)
		s.C <- node.id
	}
}

// 二叉堆节点
type timerNode struct {
	id       int64 // unique id
	deadline int64 // unix nano
	index    int   // heap index
}

func newTimerNode(id, deadline int64) *timerNode {
	return &timerNode{
		id:       id,
		deadline: deadline,
	}
}

type timerHeap []*timerNode

func (q timerHeap) Len() int {
	return len(q)
}

func (q timerHeap) Less(i, j int) bool {
	if q[i].deadline == q[j].deadline {
		return q[i].id > q[j].id
	}
	return q[i].deadline < q[j].deadline
}

func (q timerHeap) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *timerHeap) Push(x any) {
	v := x.(*timerNode)
	v.index = len(*q)
	*q = append(*q, v)
}

func (q *timerHeap) Pop() any {
	old := *q
	n := len(old)
	if n > 0 {
		v := old[n-1]
		v.index = -1 // for safety
		*q = old[:n-1]
		return v
	}
	return nil
}
