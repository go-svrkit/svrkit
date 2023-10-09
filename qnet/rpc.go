// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"

	"backend/protos"
	"backend/svrkit/logger"
)

type RPCHandler func(protos.ErrCode, protos.Message) error

// RPC上下文
type RPCContext struct {
	index int // 在最小堆里的索引
	resp  protos.Message

	Deadline int64      // 超时时间
	Seq      uint32     //
	Errno    int32      //
	Action   RPCHandler // 异步回调
}

func NewRPCContext(action RPCHandler, deadline int64) *RPCContext {
	return &RPCContext{
		Action:   action,
		Deadline: deadline,
	}
}

func (r *RPCContext) Run() error {
	if r.Action != nil {
		return r.Action(protos.ErrCode(r.Errno), r.resp)
	}
	return nil
}

type NodeHeap []*RPCContext

func (q NodeHeap) Len() int {
	return len(q)
}

func (q NodeHeap) Less(i, j int) bool {
	return q[i].Deadline < q[j].Deadline
}

func (q NodeHeap) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *NodeHeap) Push(x any) {
	v := x.(*RPCContext)
	v.index = len(*q)
	*q = append(*q, v)
}

func (q *NodeHeap) Pop() any {
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

// RPC client stub
type RPCClient struct {
	done       chan struct{}
	wg         sync.WaitGroup         //
	guard      sync.Mutex             // guard heap and pending
	heap       NodeHeap               // 使用最小堆减少主动检测超时节点
	pending    map[uint32]*RPCContext // 待响应的RPC
	expireChan chan *RPCContext       //
	counter    atomic.Uint32          // 序列号生成
	sender     func(*NetMessage) error
}

func (c *RPCClient) Init(sender func(*NetMessage) error) *RPCClient {
	c.done = make(chan struct{})
	c.pending = make(map[uint32]*RPCContext)
	c.expireChan = make(chan *RPCContext, 1024)
	c.sender = sender
	return c
}

func (c *RPCClient) Start() {
	var ready = make(chan struct{}, 1)
	c.wg.Add(1)
	go c.worker(ready)
	<-ready
}

func (c *RPCClient) Shutdown() {
	if c.done != nil {
		close(c.done)
		c.wg.Wait()
		c.done = nil
	}
}

// 异步调用
func (c *RPCClient) AsyncCall(req protos.Message, callback RPCHandler) {
	c.guard.Lock()
	defer c.guard.Unlock()

	var ttl = time.Now().Add(time.Minute).UnixNano() // 1分钟ttl
	var rpc = NewRPCContext(callback, ttl)
	var seq = c.counter.Add(1)
	rpc.Seq = seq
	c.pending[seq] = rpc
	heap.Push(&c.heap, rpc)

	var msgId = protos.GetMessageIDOf(req)
	var message = NewNetMessage(protos.MsgID(msgId), seq, req)
	// here may block
	if err := c.sender(message); err != nil {
		logger.Errorf("send message %d failed: %v", message.MsgID, err)
	}
}

func (c *RPCClient) getAndDelete(seq uint32) *RPCContext {
	c.guard.Lock()
	defer c.guard.Unlock()

	ctx, found := c.pending[seq]
	if found {
		delete(c.pending, seq)
		heap.Remove(&c.heap, ctx.index)
		ctx.index = -1
	}
	return ctx
}

func (c *RPCClient) Expired() <-chan *RPCContext {
	return c.expireChan
}

// 在主线程运行
func (c *RPCClient) Dispatch(ack *NetMessage) (bool, error) {
	var rpc = c.getAndDelete(ack.Seq)
	if rpc == nil {
		return false, nil
	}
	rpc.resp = ack.Body
	var err = rpc.Run()
	return true, err
}

func (c *RPCClient) clearTimeout(now int64) {
	for {
		c.guard.Lock()
		if len(c.heap) == 0 {
			c.guard.Unlock()
			return
		}
		var rpc = c.heap[0] // peek first item of heap
		if now < rpc.Deadline {
			c.guard.Unlock()
			break // no new context expired
		}
		heap.Pop(&c.heap)
		delete(c.pending, rpc.Seq)
		c.guard.Unlock()

		rpc.index = -1
		rpc.Errno = int32(protos.ErrRequestTimeout)
		c.expireChan <- rpc
	}
}

// 处理超时
func (c *RPCClient) worker(ready chan struct{}) {
	defer c.wg.Done()

	var ticker = time.NewTicker(time.Second)
	defer ticker.Stop()
	ready <- struct{}{}

	for {
		select {
		case now := <-ticker.C:
			c.clearTimeout(now.UnixNano() / int64(time.Millisecond))

		case <-c.done:
			return
		}
	}
}
