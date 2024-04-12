// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package handler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/codec"
	"gopkg.in/svrkit.v1/qnet"
)

func nowNano() int64 {
	return time.Now().UnixNano()
}

func TestRegister(t *testing.T) {
	defer dp.Clear()

	assert.False(t, dp.HasRegistered(1234))
	Register(1234, func(codec.Message) error { return nil })
	assert.True(t, dp.HasRegistered(1234))
}

func TestDeregister(t *testing.T) {
	defer dp.Clear()

	assert.Nil(t, dp.Deregister(1234))
	Register(1234, func(codec.Message) error { return nil })
	assert.NotNil(t, dp.Deregister(1234))
}

func TestHandle(t *testing.T) {
	defer dp.Clear()

	var triggerAt = make(map[int]int64)
	Register(101, func(codec.Message) error {
		triggerAt[1] = nowNano()
		return nil
	})
	Register(102, func(codec.Message) (codec.Message, error) {
		triggerAt[2] = nowNano()
		return nil, nil
	})
	Register(103, func(context.Context, codec.Message) error {
		triggerAt[3] = nowNano()
		return nil
	})
	Register(104, func(context.Context, codec.Message) (codec.Message, error) {
		triggerAt[4] = nowNano()
		return nil, nil
	})
	Register(105, func(context.Context, *qnet.NetMessage) error {
		triggerAt[5] = nowNano()
		return nil
	})

	for i, cmd := range []uint32{101, 102, 103, 104, 105} {
		var msg = qnet.CreateNetMessage(cmd, 0, nil)
		resp, err := Handle(context.Background(), msg)
		assert.Nil(t, resp)
		assert.Nil(t, err)
		assert.Greater(t, triggerAt[i+1], int64(0))
	}
}

func TestBeforeHook(t *testing.T) {
	defer dp.Clear()

	var t1, t2, t3 int64
	Register(101, func(codec.Message) error {
		t1 = nowNano()
		return nil
	})
	dp.RegisterBeforeHook(false, func(context.Context, *qnet.NetMessage) bool {
		t2 = nowNano()
		assert.Greater(t, t2, t1)
		return true
	})
	dp.RegisterBeforeHook(true, func(context.Context, *qnet.NetMessage) bool {
		t3 = nowNano()
		assert.Greater(t, t3, t2)
		return true
	})

	var msg = qnet.CreateNetMessage(101, 0, nil)
	Handle(context.Background(), msg)
	assert.Greater(t, t1, int64(0))
}

func TestAfterHook(t *testing.T) {
	defer dp.Clear()

	var t1, t2, t3 int64
	Register(101, func(codec.Message) error {
		t1 = nowNano()
		assert.Greater(t, t1, t2)
		return nil
	})
	dp.RegisterAfterHook(false, func(context.Context, *qnet.NetMessage) {
		t3 = nowNano()
	})
	dp.RegisterAfterHook(true, func(context.Context, *qnet.NetMessage) {
		t2 = nowNano()
		assert.Greater(t, t2, t3)
	})

	var msg = qnet.CreateNetMessage(101, 0, nil)
	Handle(context.Background(), msg)
	assert.Greater(t, t3, int64(0))
}
