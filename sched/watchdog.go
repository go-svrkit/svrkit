// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package sched

import (
	"fmt"
	"os"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/svrkit.v1/debug"
	"gopkg.in/svrkit.v1/slog"
)

type WatchDog struct {
	path       string
	wg         sync.WaitGroup
	running    atomic.Bool
	lastUpdate int64
	ttl        int64
	ticker     *time.Ticker
}

func NewWatchDog(path string, ttl int64) *WatchDog {
	return &WatchDog{
		path:       path,
		ttl:        ttl,
		lastUpdate: time.Now().Unix(),
	}
}

func (wd *WatchDog) Go() {
	wd.wg.Add(1)
	go wd.worker()
}

func (wd *WatchDog) KeepAlive() {
	wd.lastUpdate = time.Now().Unix()
}

func (wd *WatchDog) Stop() {
	if !wd.running.CompareAndSwap(true, false) {
		return
	}
	if wd.ticker != nil {
		wd.ticker.Stop()
	}
	wd.wg.Wait()
	wd.ticker = nil
}

func (wd *WatchDog) worker() {
	defer func() {
		if v := recover(); v != nil {
			debug.Backtrace("watchdog panic", v, os.Stderr)
			slog.Errorf("watchdog panic: %v", v)
		}
		wd.wg.Done()
	}()

	wd.running.Store(true)
	wd.ticker = time.NewTicker(time.Second * time.Duration(wd.ttl) / 2)
	defer wd.ticker.Stop()

	for wd.running.Load() {
		select {
		case now := <-wd.ticker.C:
			if now.Unix()-wd.lastUpdate > wd.ttl {
				go wd.dumpCPU()
				wd.dumpGoroutines()
				return
			}
		}
	}
}

func (wd *WatchDog) genPprofFileName(name string) string {
	var now = time.Now()
	var filename = fmt.Sprintf("%s_%d%02d%02d_%02d%02d%02d_%s.pprof", wd.path, now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), name)
	return filename
}

func (wd *WatchDog) dumpGoroutines() {
	var filename = wd.genPprofFileName("goroutine")
	profile := pprof.Lookup("goroutine")
	if f, err := os.Create(filename); err == nil {
		if er := profile.WriteTo(f, 0); er != nil {
			slog.Errorf("write goroutine profile failed: %v", er)
		}
	} else {
		slog.Errorf("create goroutine profile file failed: %v", err)
	}
}

func (wd *WatchDog) dumpCPU() {
	var filename = wd.genPprofFileName("cpu")
	f, err := os.Create(filename)
	if err != nil {
		slog.Errorf("create cpu profile file failed: %v", err)
		return
	}
	var timer = time.NewTimer(time.Minute)
	defer func() {
		pprof.StopCPUProfile()
		timer.Stop()
		if err := f.Close(); err != nil {
			slog.Errorf("close cpu profile file failed: %v", err)
		}
	}()

	if err = pprof.StartCPUProfile(f); err != nil {
		return
	}
	for {
		select {
		case <-timer.C:
			return
		}
	}
}
