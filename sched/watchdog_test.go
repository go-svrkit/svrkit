package sched

import (
	"testing"
	"time"
)

func TestWatchDog_NewWatchDog(t *testing.T) {
	var dog = NewWatchDog("test", 1)
	dog.Stop()
	dog.Go()
	<-time.NewTimer(time.Millisecond * 1500).C
	dog.Stop()
}
