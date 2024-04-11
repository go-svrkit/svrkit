// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package sched

import (
	"sync/atomic"
)

const (
	StateInit       = 0
	StateStarting   = 1
	StateRunning    = 2
	StateClosing    = 3
	StateTerminated = 4
)

// State of service
type State int32

func (s *State) Get() int32 {
	return atomic.LoadInt32((*int32)(s))
}

func (s *State) Set(n int32) {
	atomic.StoreInt32((*int32)(s), n)
}

func (s *State) CAS(old, new int32) bool {
	return atomic.CompareAndSwapInt32((*int32)(s), old, new)
}

func (s *State) SetStarting() {
	s.Set(StateStarting)
}

func (s *State) IsStarting() bool {
	return s.Get() == StateStarting
}

func (s *State) SetRunning() {
	s.Set(StateRunning)
}

func (s *State) IsRunning() bool {
	return s.Get() == StateRunning
}

func (s *State) SetClosing() {
	s.Set(StateClosing)
}

func (s *State) IsClosing() bool {
	return s.Get() == StateClosing
}

func (s *State) SetTerminated() {
	s.Set(StateTerminated)
}

func (s *State) IsTerminated() bool {
	return s.Get() == StateTerminated
}
