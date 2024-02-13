package handler

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/qnet"
	"testing"
	"time"
)

func nowNano() int64 {
	return time.Now().UnixNano()
}

func TestRegister(t *testing.T) {
	defer Clear()
	assert.False(t, HasRegistered(1234))
	RegisterV1(1234, func(proto.Message) error { return nil })
	assert.True(t, HasRegistered(1234))
}

func TestDeregister(t *testing.T) {
	defer Clear()
	assert.Nil(t, Deregister(1234))
	RegisterV1(1234, func(proto.Message) error { return nil })
	assert.NotNil(t, Deregister(1234))
}

func TestHandle(t *testing.T) {
	defer Clear()
	var triggerAt = make(map[int]int64)
	RegisterV1(101, func(proto.Message) error {
		triggerAt[1] = nowNano()
		return nil
	})
	RegisterV2(102, func(proto.Message) (proto.Message, error) {
		triggerAt[2] = nowNano()
		return nil, nil
	})
	RegisterV3(103, func(context.Context, proto.Message) error {
		triggerAt[3] = nowNano()
		return nil
	})
	RegisterV4(104, func(context.Context, proto.Message) (proto.Message, error) {
		triggerAt[4] = nowNano()
		return nil, nil
	})
	RegisterV5(105, func(context.Context, *qnet.NetMessage) error {
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

func TestPreHook(t *testing.T) {
	defer Clear()
	var t1, t2, t3 int64
	RegisterV1(101, func(proto.Message) error {
		t1 = nowNano()
		return nil
	})
	RegisterPreHook(false, func(context.Context, *qnet.NetMessage) bool {
		t2 = nowNano()
		assert.Greater(t, t2, t1)
		return true
	})
	RegisterPreHook(true, func(context.Context, *qnet.NetMessage) bool {
		t3 = nowNano()
		assert.Greater(t, t3, t2)
		return true
	})

	var msg = qnet.CreateNetMessage(101, 0, nil)
	Handle(context.Background(), msg)
	assert.Greater(t, t1, int64(0))
}

func TestPostHook(t *testing.T) {
	defer Clear()
	var t1, t2, t3 int64
	RegisterV1(101, func(proto.Message) error {
		t1 = nowNano()
		assert.Greater(t, t1, t2)
		return nil
	})
	RegisterPostHook(false, func(context.Context, *qnet.NetMessage) {
		t3 = nowNano()
	})
	RegisterPostHook(true, func(context.Context, *qnet.NetMessage) {
		t2 = nowNano()
		assert.Greater(t, t2, t3)
	})

	var msg = qnet.CreateNetMessage(101, 0, nil)
	Handle(context.Background(), msg)
	assert.Greater(t, t3, int64(0))
}
