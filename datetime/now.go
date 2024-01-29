// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package datetime

import (
	"time"
)

// global clock offset
var clockOffset time.Duration

// Now current time
func Now() time.Time {
	return time.Now().UTC().Add(clockOffset)
}

func NowNano() int64 {
	return time.Now().Add(clockOffset).UnixNano()
}

// NowMs current time in millisecond
func NowMs() int64 {
	return time.Now().Add(clockOffset).UnixNano() / int64(time.Millisecond)
}

func NowTime() string {
	return Now().Format(DateLayout)
}

func ResetClock() {
	clockOffset = 0
}

func ClockOffset() time.Duration {
	return clockOffset
}

func MoveClock(d time.Duration) {
	clockOffset += d
}
