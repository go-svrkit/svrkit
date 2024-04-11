// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package handler

import (
	"context"
	"fmt"

	"gopkg.in/svrkit.v1/codec"
	"gopkg.in/svrkit.v1/qnet"
	"gopkg.in/svrkit.v1/zlog"
)

type (
	MessageHandlerV1 = func(codec.Message) error
	MessageHandlerV2 = func(codec.Message) (codec.Message, error)
	MessageHandlerV3 = func(context.Context, codec.Message) error
	MessageHandlerV4 = func(context.Context, codec.Message) (codec.Message, error)
	MessageHandlerV5 = func(context.Context, *qnet.NetMessage) error

	BeforeHookFunc func(context.Context, *qnet.NetMessage) bool
	AfterHookFunc  func(context.Context, *qnet.NetMessage)
)

type IHandler interface {
	MessageHandlerV1 | MessageHandlerV2 | MessageHandlerV3 | MessageHandlerV4 | MessageHandlerV5
}

type Dispatcher struct {
	handlers    map[uint32]any
	once        map[uint32]bool
	beforeHooks []BeforeHookFunc
	afterHooks  []AfterHookFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[uint32]any),
		once:     make(map[uint32]bool),
	}
}

func (d *Dispatcher) Clear() {
	clear(d.handlers)
	clear(d.once)
	d.beforeHooks = nil
	d.afterHooks = nil
}

func (d *Dispatcher) Register(cmd uint32, action any) bool {
	if d.HasRegistered(cmd) {
		zlog.Warnf("duplicate handler registration of cmd %d", cmd)
	}
	switch action.(type) {
	case MessageHandlerV1, MessageHandlerV2, MessageHandlerV3, MessageHandlerV4, MessageHandlerV5:
		d.handlers[cmd] = action
		return true
	}
	return false
}

func (d *Dispatcher) RegisterOnce(cmd uint32, action any) {
	d.once[cmd] = true
	d.Register(cmd, action)
}

// HasRegistered 是否注册
func (d *Dispatcher) HasRegistered(cmd uint32) bool {
	_, found := d.handlers[cmd]
	return found
}

// Deregister 取消所有
func (d *Dispatcher) Deregister(cmd uint32) any {
	var old = d.handlers[cmd]
	delete(d.handlers, cmd)
	return old
}

func (d *Dispatcher) RegisterBeforeHook(prepend bool, h BeforeHookFunc) {
	if prepend {
		d.beforeHooks = append([]BeforeHookFunc{h}, d.beforeHooks...)
	} else {
		d.beforeHooks = append(d.beforeHooks, h)
	}
}

func (d *Dispatcher) RegisterAfterHook(prepend bool, h AfterHookFunc) {
	if prepend {
		d.afterHooks = append([]AfterHookFunc{h}, d.afterHooks...)
	} else {
		d.afterHooks = append(d.afterHooks, h)
	}
}

func (d *Dispatcher) invokePreHooks(ctx context.Context, msg *qnet.NetMessage) bool {
	for _, h := range d.beforeHooks {
		if !h(ctx, msg) {
			return false // stop continue
		}
	}
	return true
}

func (d *Dispatcher) invokePostHooks(ctx context.Context, msg *qnet.NetMessage) {
	for _, h := range d.afterHooks {
		h(ctx, msg)
	}
}

func (d *Dispatcher) Handle(ctx context.Context, message *qnet.NetMessage) (codec.Message, error) {
	if !d.invokePreHooks(ctx, message) {
		return nil, nil // stop continue
	}
	defer d.invokePostHooks(ctx, message)

	var h = d.handlers[message.Command]
	return d.dispatch(ctx, h, message)
}

func (d *Dispatcher) dispatch(ctx context.Context, action any, msg *qnet.NetMessage) (resp codec.Message, err error) {
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
