// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package timers

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"
)

type timerContext struct {
	fireCount    int
	startTime    time.Time
	lastFireTime time.Time
}

func newTimerContext() *timerContext {
	return &timerContext{
		startTime: time.Now(),
	}
}

func (r *timerContext) Run() error {
	r.lastFireTime = time.Now()
	r.fireCount++
	return nil
}

func testTimerCancel(t *testing.T, sched TimerScheduler) {
	const interval = 1000 // 1s
	var timerCtx = newTimerContext()

	var timerId = sched.AddTimer(interval)
	time.Sleep(time.Millisecond) // wait for worker
	if n := sched.Size(); n != 1 {
		t.Fatalf("timer size unexpected %d", n)
	}
	sched.CancelTimer(timerId)
	if n := sched.Size(); n != 0 {
		t.Fatalf("timer size unexpected %d", n)
	}
	if timerCtx.fireCount > 0 {
		t.Fatalf("timeout %d unexpectly triggered", timerId)
	}
}

func testTimerRunAfter(t *testing.T, sched TimerScheduler, interval int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	var timerCtx = newTimerContext()
	log.Printf("timer start at %s, interval is %d", timerCtx.startTime.Format(time.RFC3339), interval)
	sched.AddTimer(interval)

	for timerCtx.fireCount == 0 {
		select {
		case <-sched.TimedOutChan():
			//msg.Runner.Run()
			duration := timerCtx.lastFireTime.Sub(timerCtx.startTime)
			log.Printf("timer fired after %v at %s", duration, timerCtx.lastFireTime.Format(time.RFC3339))
			// 允许1个timeUnit(毫秒）的误差
			if d := duration.Milliseconds(); d < int64(interval) && (int64(interval)-d) > 1 {
				t.Fatalf("fired too early %v != %v", duration, interval)
			}
			timerCtx.startTime = time.Now()

		case <-ctx.Done():
			t.Fatalf("test deadline exceeded")
			return
		}
	}
}

func TestTimerQueue_RunAfter(t *testing.T) {
	var timer = NewDefaultTimerQueue()
	timer.Start()
	defer timer.Shutdown()

	testTimerCancel(t, timer)
	for i := 100; i <= 1000; i += 100 {
		testTimerRunAfter(t, timer, int64(i))
	}
}

func timerExpireWorker(t *testing.T, ctx context.Context, wg *sync.WaitGroup, sched TimerScheduler) {
	var fired = make([]int64, 0, 100)
	defer func() {
		t.Logf("timer trigger order is: %v", fired)
		wg.Done()
	}()

	for {
		select {
		case id := <-sched.TimedOutChan():
			var firedAt = time.Now().Format(time.RFC3339)
			t.Logf("timer %d fired at %s", id, firedAt)

		case <-ctx.Done():
			return
		}
	}
}

func testTimerRunFIFO(t *testing.T, sched TimerScheduler) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go timerExpireWorker(t, ctx, &wg, sched)

	var deadline = time.Now().Add(5 * time.Second) // 让所有timer都在5s后过期
	t.Logf("all timers expect to be fired at %v", deadline.Format(time.RFC3339))

	for i := 0; i < 100; i++ {
		sched.AddTimerAt(deadline.UnixNano())
	}
	wg.Wait()
}

func TestTimerQueue_RunFIFO(t *testing.T) {
	var timer = NewDefaultTimerQueue()
	timer.Start()
	defer timer.Shutdown()

	testTimerRunFIFO(t, timer)
}
