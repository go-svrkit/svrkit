// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package sched

const (
	TaskStateInit = 0
	TaskScheduled = 1 // task is scheduled for execution
	TaskExecuted  = 2 // a non-repeating task has already executed (or is currently executing) and has not been cancelled.
	TaskCancelled = 3 // task has been cancelled (with a call to Cancel).
)

// IRunner r代表一个可执行对象
type IRunner interface {
	Run() error
}

type Runnable struct {
	F func() error
}

func NewRunnable(f func() error) IRunner {
	return &Runnable{F: f}
}

func (r *Runnable) Run() error {
	return r.F()
}
