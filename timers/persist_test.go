// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package timers

import (
	"testing"
)

func TestDumpTimers(t *testing.T) {
	for i := 0; i < 100; i++ {
		AddTimer(1000, 10000, int32(i), int32(i*10), nil)
	}
	data, err := DumpTimers()
	if err != nil {
		t.Fatalf("DumpTimers failed: %v", err)
	}
	Shutdown()
	defTimer = NewDefaultTimerQueue()
	if err := LoadTimers(data); err != nil {
		t.Fatalf("LoadTimers failed: %v", err)
	}
}
