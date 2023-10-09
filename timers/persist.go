// Copyright © 2023 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package timers

import (
	"github.com/fxamacker/cbor/v2"
)

type TimerInfo struct {
	Id       int64 `json:"id"`
	Deadline int64 `json:"deadline"`
	Owner    int64 `json:"owner"`
	Action   int32 `json:"action"`
	Param    int32 `json:"param,omitempty"`
}

type AllTimersInfo struct {
	Timers []TimerInfo `json:"timers"`
	NextId int64       `json:"next_id"`
}

func DumpTimers() ([]byte, error) {
	guard.Lock()
	defer guard.Unlock()

	var info = &AllTimersInfo{
		NextId: defTimer.NextID(),
	}
	info.Timers = make([]TimerInfo, 0, defTimer.Size())
	defTimer.RangeTimers(func(node *timerNode) {
		var msg = timeouts[node.id]
		if msg != nil && msg.Owner > 0 {
			var ti = TimerInfo{
				Id:       node.id,
				Deadline: node.deadline,
				Owner:    msg.Owner,
				Action:   msg.Action,
				Param:    msg.Param,
			}
			info.Timers = append(info.Timers, ti)
		}
	})
	// use RFC8949 Concise Binary Object Representation (CBOR) format to encode timer data
	return cbor.Marshal(info)
}

func LoadTimers(data []byte) error {
	var info AllTimersInfo
	if err := cbor.Unmarshal(data, &info); err != nil {
		return err
	}

	guard.Lock()
	defTimer.nextId.Store(info.NextId)
	guard.Unlock()

	for _, ti := range info.Timers {
		var msg = &TimeoutMsg{
			Owner:  ti.Owner,
			Action: ti.Action,
			Param:  ti.Param,
		}
		AddTimerAt(ti.Deadline, msg)
	}
	return nil
}
