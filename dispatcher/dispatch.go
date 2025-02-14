// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package handler

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/svrkit.v1/debug"
	"gopkg.in/svrkit.v1/qlog"
	"gopkg.in/svrkit.v1/qnet"

	"github.com/golang/protobuf/proto"
)

type (
	MessageHandlerV1 = func(proto.Message) error
	MessageHandlerV2 = func(proto.Message) (proto.Message, error)
	MessageHandlerV3 = func(context.Context, proto.Message) error
	MessageHandlerV4 = func(context.Context, proto.Message) (proto.Message, error)
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
	errCapture  func(*qnet.NetMessage, any)
}

var dp = NewDispatcher()

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers:   make(map[uint32]any),
		once:       make(map[uint32]bool),
		errCapture: onError,
	}
}

func (d *Dispatcher) Clear() {
	clear(d.handlers)
	clear(d.once)
	d.beforeHooks = nil
	d.afterHooks = nil
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

func (d *Dispatcher) SetErrorCapture(f func(*qnet.NetMessage, any)) {
	d.errCapture = f
}

func (d *Dispatcher) Handle(ctx context.Context, message *qnet.NetMessage) (proto.Message, error) {
	defer func() {
		if v := recover(); v != nil {
			d.errCapture(message, v)
		}
	}()
	if !d.invokePreHooks(ctx, message) {
		return nil, nil // stop continue
	}
	defer func() {
		if d.once[message.Command] {
			delete(d.handlers, message.Command)
			delete(d.once, message.Command)
		}
		d.invokePostHooks(ctx, message)
	}()

	var h = d.handlers[message.Command]
	return d.dispatch(ctx, h, message)
}

func (d *Dispatcher) dispatch(ctx context.Context, action any, msg *qnet.NetMessage) (resp proto.Message, err error) {
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

func onError(msg *qnet.NetMessage, err any) {
	var sb strings.Builder
	var title = fmt.Sprintf("dispatch message %d", msg.Command)
	debug.TraceStack(1, title, err, &sb)
	qlog.Error(sb.String())
}

func G() *Dispatcher {
	return dp
}

func Register[F IHandler](cmd uint32, action F) {
	if dp.HasRegistered(cmd) {
		qlog.Warnf("duplicate handler registration of message %d", cmd)
	}
	if action != nil {
		dp.handlers[cmd] = action
	}
}

// HasRegistered 是否注册
func HasRegistered(cmd uint32) bool {
	return dp.HasRegistered(cmd)
}

// Deregister 取消所有
func Deregister(cmd uint32) any {
	return dp.Deregister(cmd)
}

func RegOnce[F IHandler](cmd uint32, action F) {
	Register(cmd, action)
	dp.once[cmd] = true
}

func Handle(ctx context.Context, message *qnet.NetMessage) (proto.Message, error) {
	return dp.Handle(ctx, message)
}
