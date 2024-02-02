// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package timers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func randTimerData() TimerData {
	var randVal = rand.Int63()
	return TimerData{
		Id:       randVal % 1000,
		Owner:    randVal % 10000,
		Action:   int32(randVal % 100),
		Deadline: rand.Int63(),
		Arg:      randVal % 0xFFFFF,
	}
}

func TestMarshalTimers1(t *testing.T) {
	var info AllTimersData
	info.NextId = 1
	for i := 0; i < 100; i++ {
		info.Timers = append(info.Timers, randTimerData())
	}
	data, err := MarshalTimers(&info)
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	var info2 AllTimersData
	err = UnmarshalTimers(data, &info2)
	assert.Nil(t, err)

	assert.True(t, reflect.DeepEqual(info, info2))
}

func TestDumpTimers(t *testing.T) {
	td := []struct {
		owner    int64
		duration int64
		action   int32
		arg      int64
		data     any
	}{
		{1, 100, 1, 1, 1},
		{2, 200, 2, 2, 2},
		{3, 300, 3, 3, 3},
	}
	for _, d := range td {
		AddTimer(d.owner, d.duration, d.action, d.arg, d.data)
	}

	var allData = DumpTimers()
	assert.NotNil(t, allData)
	assert.Equal(t, len(allData.Timers), len(td))
}
