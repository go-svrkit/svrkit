// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package timers

import (
	"container/heap"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TickInterval = 10 * time.Millisecond
)

// TimerQueue 最小堆实现的定时器
// 注意：对相同过期时间的多个Timer，TimerQueue不保证FIFO触发
type TimerQueue struct {
	done         chan struct{}
	wg           sync.WaitGroup //
	running      atomic.Int32
	tickInterval time.Duration

	guard  sync.Mutex          // 多线程
	refer  map[int64]*heapNode // O(1)查找
	timers timerHeap           //
	C      chan int64          // 到期的定时器
}

var _ TimerScheduler = (*TimerQueue)(nil)

func NewTimerQueue(capacity int) *TimerQueue {
	return new(TimerQueue).Init(capacity, TickInterval)
}

func (s *TimerQueue) Init(capacity int, tickInterval time.Duration) *TimerQueue {
	if tickInterval < time.Millisecond {
		tickInterval = time.Millisecond
	}
	s.tickInterval = tickInterval
	s.done = make(chan struct{})
	s.timers = make(timerHeap, 0, 1024)
	s.refer = make(map[int64]*heapNode, 1024)
	s.C = make(chan int64, capacity)
	return s
}

func (s *TimerQueue) Size() int {
	s.guard.Lock()
	var n = len(s.refer)
	s.guard.Unlock()
	return n
}

func (s *TimerQueue) TimedOutChan() <-chan int64 {
	return s.C
}

// IsPending 判断定时器是否在等待触发
func (s *TimerQueue) IsPending(id int64) bool {
	s.guard.Lock()
	var node = s.refer[id]
	s.guard.Unlock()
	return node != nil
}

// Start starts the background thread explicitly
func (s *TimerQueue) Start() error {
	if !s.running.CompareAndSwap(0, 1) {
		return nil
	}
	s.wg.Add(1)
	var ready = make(chan struct{})
	go s.worker(ready)
	<-ready
	return nil
}

func (s *TimerQueue) Shutdown() {
	if !s.running.CompareAndSwap(1, 0) {
		return
	}
	close(s.done)
	s.wg.Wait()

	s.guard.Lock()
	defer s.guard.Unlock()

	close(s.C)
	s.C = nil
	s.refer = nil
	s.timers = nil
}

// AddTimeoutAt 创建在指定时机触发的定时器 `deadline`使用纳秒
func (s *TimerQueue) AddTimeoutAt(tid int64, deadline int64) {
	s.guard.Lock()
	defer s.guard.Unlock()

	var node = newTimerNode(tid, deadline)
	s.refer[tid] = node
	heap.Push(&s.timers, node)
}

// AddTimeout 创建一个定时器，在`delayMs`毫秒后过期
func (s *TimerQueue) AddTimeout(tid int64, delayMs int64) {
	var deadline = currentUnixNano() + delayMs*int64(time.Millisecond)
	if delayMs > 0 && deadline < 0 {
		deadline = math.MaxInt64 // guard against overflow
	}
	s.AddTimeoutAt(tid, deadline)
}

// CancelTimeout 取消一个timer
func (s *TimerQueue) CancelTimeout(id int64) bool {
	s.guard.Lock()
	defer s.guard.Unlock()

	if node, found := s.refer[id]; found {
		heap.Remove(&s.timers, node.index)
		delete(s.refer, id)
		return true
	}
	return false
}

func (s *TimerQueue) Range(action func(id, deadline int64)) {
	s.guard.Lock()
	defer s.guard.Unlock()

	for _, node := range s.refer {
		action(node.id, node.deadline)
	}
}

func (s *TimerQueue) worker(ready chan struct{}) {
	defer s.wg.Done()

	var ticker = time.NewTicker(s.tickInterval)
	defer ticker.Stop()

	ready <- struct{}{}

	for {
		select {
		case <-ticker.C:
			s.update(currentUnixNano())

		case <-s.done:
			return
		}
	}
}

// 返回触发的timer列表
func (s *TimerQueue) update(nowNano int64) {
	s.guard.Lock()
	defer s.guard.Unlock()

	for len(s.timers) > 0 {
		var node = s.timers[0] // peek first item of heap
		if nowNano < node.deadline {
			break // no new timer expired
		}

		heap.Pop(&s.timers)
		delete(s.refer, node.id)
		s.C <- node.id
		//log.Printf("timer %d expired deadline=%s\n", node.id, time.Unix(0, node.deadline).Format(datetime.TimestampLayout))
	}
}

// 二叉堆节点
type heapNode struct {
	id       int64 // unique id
	deadline int64 // unix nano
	index    int   // heap index
}

func newTimerNode(id, deadline int64) *heapNode {
	return &heapNode{
		id:       id,
		deadline: deadline,
	}
}

type timerHeap []*heapNode

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
	v := x.(*heapNode)
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
