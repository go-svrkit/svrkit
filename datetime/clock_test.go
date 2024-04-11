// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package datetime

import (
	"testing"
	"time"
)

func TestClockExample(t *testing.T) {
	clock := NewClock(ClockPrecision)
	clock.Go()
	defer clock.Stop()
	time.Sleep(10 * time.Millisecond)

	now := clock.DateTime()
	t.Logf("now: %v", now)

	clock.Travel(time.Hour * 2) // 往前拨2小时
	t.Logf("t1: %v", clock.DateTime())

	clock.Travel(time.Hour * -3) // 往后拨3小时
	t.Logf("t2: %v", clock.DateTime())
}
