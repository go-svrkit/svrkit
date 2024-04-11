// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package timers

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

type TimerData struct {
	Id       int64 `json:"id"`
	Deadline int64 `json:"deadline"`
	Owner    int64 `json:"owner"`
	Action   int   `json:"action"`
	Arg      int64 `json:"arg,omitempty"`
}

type AllTimersData struct {
	Timers []TimerData `json:"timers"`
	NextId int64       `json:"next_id"`
}

// MarshalTimers use msgpack format (https://msgpack.org/) to encode timer data
func MarshalTimers(info *AllTimersData) ([]byte, error) {
	var buf bytes.Buffer
	var enc = msgpack.NewEncoder(&buf)
	enc.SetOmitEmpty(true)
	enc.SetCustomStructTag("json")
	if err := enc.Encode(info); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func UnmarshalTimers(data []byte, info *AllTimersData) error {
	var dec = msgpack.NewDecoder(bytes.NewReader(data))
	dec.SetCustomStructTag("json")
	if err := dec.Decode(&info); err != nil {
		return err
	}
	return nil
}

func DumpTimers() *AllTimersData {
	gLock.Lock()
	defer gLock.Unlock()

	var info = &AllTimersData{}
	info.NextId = gTid.Add(1)
	info.Timers = make([]TimerData, 0, len(gTimeouts))
	for id, timeout := range gTimeouts {
		if timeout != nil && timeout.Owner > 0 {
			var ti = TimerData{
				Id:       id,
				Deadline: timeout.Deadline,
				Owner:    timeout.Owner,
				Action:   timeout.Action,
				Arg:      timeout.Arg,
			}
			info.Timers = append(info.Timers, ti)
		}
	}
	return info
}

func RestoreTimers(info *AllTimersData) {
	gLock.Lock()
	gTid.Store(info.NextId)
	gLock.Unlock()

	for _, ti := range info.Timers {
		var msg = &TimeoutMsg{
			Owner:  ti.Owner,
			Action: ti.Action,
			Arg:    ti.Arg,
		}
		AddTimerAt(ti.Deadline, msg)
	}
}
