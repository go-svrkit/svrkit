// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package handler

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"gopkg.in/svrkit.v1/logger"
	"gopkg.in/svrkit.v1/qnet"
)

type (
	MessageHandlerV1 func(proto.Message) error
	MessageHandlerV2 func(proto.Message) proto.Message
	MessageHandlerV3 func(proto.Message) (proto.Message, error)
	MessageHandlerV4 func(context.Context, proto.Message) error
	MessageHandlerV5 func(context.Context, proto.Message) proto.Message
	MessageHandlerV6 func(context.Context, proto.Message) (proto.Message, error)
	MessageHandlerV7 func(*qnet.NetMessage) error
	MessageHandlerV8 func(context.Context, *qnet.NetMessage) error
)

// 消息派发
var handlers = make(map[uint32]any)

func HasRegistered(msgId uint32) bool {
	_, found := handlers[msgId]
	return found
}

// Deregister 取消所有
func Deregister(msgId uint32) any {
	var old = handlers[msgId]
	delete(handlers, msgId)
	return old
}

// RegisterV1 注册消息处理函数
func RegisterV1(msgId uint32, action MessageHandlerV1) {
	if HasRegistered(msgId) {
		logger.Warnf("duplicate handler registration of message %v", msgId)
	}
	handlers[msgId] = action
}

func RegisterV2(msgId uint32, action MessageHandlerV2) {
	if HasRegistered(msgId) {
		logger.Warnf("duplicate handler registration of message %v", msgId)
	}
	handlers[msgId] = action
}

func RegisterV3(msgId uint32, action MessageHandlerV3) {
	if HasRegistered(msgId) {
		logger.Warnf("duplicate handler registration of message %v", msgId)
	}
	handlers[msgId] = action
}

func RegisterV4(msgId uint32, action MessageHandlerV4) {
	if HasRegistered(msgId) {
		logger.Warnf("duplicate handler registration of message %v", msgId)
	}
	handlers[msgId] = action
}

func RegisterV5(msgId uint32, action MessageHandlerV5) {
	if HasRegistered(msgId) {
		logger.Warnf("duplicate handler registration of message %v", msgId)
	}
	handlers[msgId] = action
}

func RegisterV6(msgId uint32, action MessageHandlerV6) {
	if HasRegistered(msgId) {
		logger.Warnf("duplicate handler registration of message %v", msgId)
	}
	handlers[msgId] = action
}

func RegisterV7(msgId uint32, action MessageHandlerV7) {
	if HasRegistered(msgId) {
		logger.Warnf("duplicate handler registration of message %v", msgId)
	}
	handlers[msgId] = action
}

func RegisterV8(msgId uint32, action MessageHandlerV8) {
	if HasRegistered(msgId) {
		logger.Warnf("duplicate handler registration of message %v", msgId)
	}
	handlers[msgId] = action
}

func Handle(ctx context.Context, message *qnet.NetMessage) (proto.Message, error) {
	var msgId = message.MsgID
	action, found := handlers[msgId]
	if !found {
		return nil, fmt.Errorf("message %v handler not found", msgId)
	}
	return dispatch(ctx, action, message)
}

func dispatch(ctx context.Context, action any, msg *qnet.NetMessage) (resp proto.Message, err error) {
	switch h := action.(type) {
	case MessageHandlerV1:
		err = h(msg.Body)
	case MessageHandlerV2:
		resp = h(msg.Body)
	case MessageHandlerV3:
		resp, err = h(msg.Body)
	case MessageHandlerV4:
		err = h(ctx, msg.Body)
	case MessageHandlerV5:
		resp = h(ctx, msg.Body)
	case MessageHandlerV6:
		resp, err = h(ctx, msg.Body)
	case MessageHandlerV7:
		err = h(msg)
	case MessageHandlerV8:
		err = h(ctx, msg)
	default:
		err = fmt.Errorf("unexpected handler type %T", h)
	}
	return resp, err
}
