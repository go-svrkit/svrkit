// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package handler

import (
	"context"
	"fmt"

	"backend/protos"
	"backend/svrkit/logger"
	"backend/svrkit/qnet"
)

type (
	MessageHandlerV1 func(protos.Message) error
	MessageHandlerV2 func(protos.Message) protos.Message
	MessageHandlerV3 func(context.Context, protos.Message) error
	MessageHandlerV4 func(context.Context, protos.Message) protos.Message
	MessageHandlerV5 func(message *qnet.NetMessage) error
)

// 消息派发
var handlers = make(map[protos.MsgID]any)

func HasRegistered(command protos.MsgID) bool {
	_, found := handlers[command]
	return found
}

// 取消所有
func Deregister(command protos.MsgID) any {
	var old = handlers[command]
	delete(handlers, command)
	return old
}

// 注册一个
func RegisterV1(command protos.MsgID, action MessageHandlerV1) {
	if HasRegistered(command) {
		logger.Warnf("duplicate handler registration of message %v", command)
	}
	handlers[command] = action
}

func RegisterV2(command protos.MsgID, action MessageHandlerV2) {
	if HasRegistered(command) {
		logger.Warnf("duplicate handler registration of message %v", command)
	}
	handlers[command] = action
}

func RegisterV3(command protos.MsgID, action MessageHandlerV3) {
	if HasRegistered(command) {
		logger.Warnf("duplicate handler registration of message %v", command)
	}
	handlers[command] = action
}

func RegisterV4(command protos.MsgID, action MessageHandlerV4) {
	if HasRegistered(command) {
		logger.Warnf("duplicate handler registration of message %v", command)
	}
	handlers[command] = action
}

func RegisterV5(command protos.MsgID, action MessageHandlerV5) {
	if HasRegistered(command) {
		logger.Warnf("duplicate handler registration of message %v", command)
	}
	handlers[command] = action
}

func Handle(ctx context.Context, message *qnet.NetMessage) (protos.Message, error) {
	var msgId = protos.MsgID(message.MsgID)
	action, found := handlers[msgId]
	if !found {
		return nil, fmt.Errorf("message %v handler not found", msgId)
	}
	return dispatch(ctx, action, message)
}

func dispatch(ctx context.Context, action any, message *qnet.NetMessage) (resp protos.Message, err error) {
	switch h := action.(type) {
	case MessageHandlerV1:
		err = h(message.Body)
	case MessageHandlerV2:
		resp = h(message.Body)
	case MessageHandlerV3:
		err = h(ctx, message.Body)
	case MessageHandlerV4:
		resp = h(ctx, message.Body)
	case MessageHandlerV5:
		err = h(message)
	default:
		err = fmt.Errorf("unexpected handler type %T", h)
	}
	return resp, err
}
