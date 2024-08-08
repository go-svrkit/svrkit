// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package datetime

import (
	"sync"
	"time"
)

const (
	ClockPrecision = time.Millisecond * 100
	ISO8601Format  = "2006-01-02T15:04:05-0700"
)

// Clock 提供一些对壁钟时间的操作，不适用于高精度的计时场景
// 设计初衷是为精度至少为秒的上层业务服务, 支持时钟的往前/后调拨
type Clock struct {
	done     chan struct{}
	wg       sync.WaitGroup
	guard    sync.RWMutex
	traveled time.Duration // 旅行时间，提供对时钟的往前/后拨动
	now      time.Time
	ticker   *time.Ticker //
}

func NewClock(interval time.Duration) *Clock {
	if interval <= 0 {
		interval = ClockPrecision
	}
	c := &Clock{
		done:   make(chan struct{}),
		now:    time.Now(),
		ticker: time.NewTicker(interval),
	}
	return c
}

func (c *Clock) Go() {
	c.wg.Add(1)
	go c.serve()
}

func (c *Clock) serve() {
	defer c.wg.Done()
	for {
		select {
		case t, ok := <-c.ticker.C:
			if !ok {
				return
			}
			c.guard.Lock()
			c.now = t
			c.guard.Unlock()

		case <-c.done:
			return
		}
	}
}

func (c *Clock) Stop() {
	close(c.done)
	c.wg.Wait()
	c.ticker.Stop()
	c.ticker = nil
}

func (c *Clock) Now() time.Time {
	c.guard.RLock()
	var t = c.now
	c.guard.RUnlock()
	if c.traveled != 0 {
		return t.Add(c.traveled)
	}
	return t
}

func (c *Clock) NowMs() int64 {
	return c.Now().UnixNano() / int64(time.Millisecond)
}

func (c *Clock) DateTime() string {
	now := c.Now()
	return now.Format(ISO8601Format)
}

// Reset 恢复时钟
func (c *Clock) Reset() {
	c.traveled = 0
}

// Travel 拨动时钟（时间旅行）
func (c *Clock) Travel(d time.Duration) {
	c.traveled += d
}
