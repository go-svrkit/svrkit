// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmitter_AddListener(t *testing.T) {
	var emitter = NewEmitter()
	var firstTriggered int64
	var secondTriggered int64
	emitter.AddListener("onclick", func(event *Event) error {
		firstTriggered = time.Now().UnixNano()
		time.Sleep(time.Millisecond) // simulate a slow listener
		return nil
	})
	emitter.AddListener("onclick", func(event *Event) error {
		secondTriggered = time.Now().UnixNano()
		return nil
	})
	assert.Equal(t, int64(0), firstTriggered)
	assert.Equal(t, int64(0), secondTriggered)
	assert.Equal(t, 2, emitter.ListenerCount("onclick"))
	emitter.Emit("onclick", 1245, 5678)
	assert.True(t, firstTriggered > 0)
	assert.True(t, secondTriggered > 0)
	assert.True(t, firstTriggered < secondTriggered)
}

func TestEmitter_PrependListener(t *testing.T) {
	// TODO: Add test cases.
}

func TestEmitter_RemoveListener(t *testing.T) {
	// TODO: Add test cases.
}
