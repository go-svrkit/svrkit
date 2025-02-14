// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmitter_AddListener(t *testing.T) {
	var em = NewEmitter()
	var firstTriggered int64
	var secondTriggered int64
	em.AddListener("onclick", func(event *Event) error {
		firstTriggered = time.Now().UnixNano()
		time.Sleep(time.Millisecond) // simulate a slow listener
		return nil
	})
	em.AddListener("onclick", func(event *Event) error {
		secondTriggered = time.Now().UnixNano()
		return nil
	})
	assert.Equal(t, int64(0), firstTriggered)
	assert.Equal(t, int64(0), secondTriggered)
	assert.Equal(t, 2, em.ListenerCount("onclick"))

	em.Emit("onclick", 1245, 5678)

	assert.True(t, firstTriggered > 0)
	assert.True(t, secondTriggered > 0)
	assert.True(t, firstTriggered < secondTriggered)
}

func TestEmitter_PrependListener(t *testing.T) {
	var em = NewEmitter()

	var firstTriggered int64
	var secondTriggered int64
	em.PrependListener("onclick", func(event *Event) error {
		firstTriggered = time.Now().UnixNano()
		return nil
	})
	em.PrependListener("onclick", func(event *Event) error {
		secondTriggered = time.Now().UnixNano()
		return nil
	})

	assert.Equal(t, int64(0), firstTriggered)
	assert.Equal(t, int64(0), secondTriggered)
	assert.Equal(t, 2, em.ListenerCount("onclick"))

	em.Emit("onclick", 1245, 5678)

	assert.True(t, firstTriggered > 0)
	assert.True(t, secondTriggered > 0)
	assert.True(t, firstTriggered > secondTriggered)
}

func TestEmitter_RemoveListener(t *testing.T) {
	var em = NewEmitter()

	var fn = func(event *Event) error {
		return nil
	}
	em.AddListener("onclick", fn)
	assert.Equal(t, 1, em.ListenerCount("onclick"))

	em.RemoveListener("onclick", fn)
	assert.Equal(t, 0, em.ListenerCount("onclick"))

	em.PrependListener("onclick", fn)
	assert.Equal(t, 1, em.ListenerCount("onclick"))

	em.RemoveListeners("onclick")
	assert.Equal(t, 0, em.ListenerCount("onclick"))
}

func TestEmitter_Once(t *testing.T) {
	var em = NewEmitter()
	var triggered int64
	em.Once("onclick", func(event *Event) error {
		triggered = time.Now().UnixNano()
		return nil
	})

	assert.Equal(t, int64(0), triggered)
	assert.Equal(t, 1, em.ListenerCount("onclick"))

	em.Emit("onclick", 1245, 5678)

	assert.Greater(t, triggered, int64(0))
	assert.Equal(t, 0, em.ListenerCount("onclick"))
}
