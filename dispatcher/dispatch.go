// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package handler

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"gopkg.in/svrkit.v1/qnet"
	"gopkg.in/svrkit.v1/slog"
)

type (
	MessageHandlerV1 func(proto.Message) error
	MessageHandlerV2 func(proto.Message) (proto.Message, error)
	MessageHandlerV3 func(context.Context, proto.Message) error
	MessageHandlerV4 func(context.Context, proto.Message) (proto.Message, error)
	MessageHandlerV5 func(context.Context, *qnet.NetMessage) error

	PreHookFunc  func(context.Context, *qnet.NetMessage) bool
	PostHookFunc func(context.Context, *qnet.NetMessage)
)

// 消息派发
var (
	handlers  = make(map[uint32]any)
	preHooks  []PreHookFunc
	postHooks []PostHookFunc
)

func HasRegistered(cmd uint32) bool {
	_, found := handlers[cmd]
	return found
}

// Deregister 取消所有
func Deregister(cmd uint32) any {
	var old = handlers[cmd]
	delete(handlers, cmd)
	return old
}

func RegisterPreHook(prepend bool, h PreHookFunc) {
	if prepend {
		preHooks = append([]PreHookFunc{h}, preHooks...)
	} else {
		preHooks = append(preHooks, h)
	}
}

func RegisterPostHook(prepend bool, h PostHookFunc) {
	if prepend {
		postHooks = append([]PostHookFunc{h}, postHooks...)
	} else {
		postHooks = append(postHooks, h)
	}
}

// RegisterV1 注册消息处理函数
func RegisterV1(cmd uint32, action MessageHandlerV1) {
	if HasRegistered(cmd) {
		slog.Warnf("duplicate handler registration of message %v", cmd)
	}
	handlers[cmd] = action
}

func RegisterV2(cmd uint32, action MessageHandlerV2) {
	if HasRegistered(cmd) {
		slog.Warnf("duplicate handler registration of message %v", cmd)
	}
	handlers[cmd] = action
}

func RegisterV3(cmd uint32, action MessageHandlerV3) {
	if HasRegistered(cmd) {
		slog.Warnf("duplicate handler registration of message %v", cmd)
	}
	handlers[cmd] = action
}

func RegisterV4(cmd uint32, action MessageHandlerV4) {
	if HasRegistered(cmd) {
		slog.Warnf("duplicate handler registration of message %v", cmd)
	}
	handlers[cmd] = action
}

func RegisterV5(cmd uint32, action MessageHandlerV5) {
	if HasRegistered(cmd) {
		slog.Warnf("duplicate handler registration of message %v", cmd)
	}
	handlers[cmd] = action
}

func Handle(ctx context.Context, message *qnet.NetMessage) (proto.Message, error) {
	var cmd = message.Command
	action, found := handlers[cmd]
	if !found {
		return nil, fmt.Errorf("message %v handler not found", cmd)
	}

	if !invokePreHooks(ctx, message) {
		return nil, nil // stop continue
	}
	defer invokePostHooks(ctx, message)

	return dispatch(ctx, action, message)
}

func invokePreHooks(ctx context.Context, msg *qnet.NetMessage) bool {
	for _, h := range preHooks {
		if !h(ctx, msg) {
			return false // stop continue
		}
	}
	return true
}

func invokePostHooks(ctx context.Context, msg *qnet.NetMessage) {
	for _, h := range postHooks {
		h(ctx, msg)
	}
}

func dispatch(ctx context.Context, action any, msg *qnet.NetMessage) (resp proto.Message, err error) {
	switch h := action.(type) {
	case MessageHandlerV1:
		err = h(msg.Body)
	case MessageHandlerV2:
		resp, err = h(msg.Body)
	case MessageHandlerV3:
		err = h(ctx, msg.Body)
	case MessageHandlerV4:
		resp, err = h(ctx, msg.Body)
	case MessageHandlerV5:
		err = h(ctx, msg)
	default:
		err = fmt.Errorf("unexpected handler type %T", h)
	}
	return resp, err
}
