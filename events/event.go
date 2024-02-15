// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package events

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/svrkit.v1/debug"
)

var (
	ErrDuplicateOnce   = errors.New("duplicate once call")
	ErrNoListenerFired = errors.New("no listener fired")
)

type EventListener func(event *Event) error

type IListener interface {
	Get() EventListener
	Fire(*Event) error
}

type IEventTarget interface {
	AddListener(eventName string, listener EventListener)
	RemoveListener(eventName string, listener EventListener) bool
}

type Event struct {
	Name   string
	Target IEventTarget
	Args   []any
}

func NewEvent(target IEventTarget, name string, args []any) *Event {
	return &Event{
		Target: target,
		Name:   name,
		Args:   args,
	}
}

type EventHandler struct {
	callback EventListener
	target   IEventTarget
}

func NewEventHandler(target IEventTarget, cb EventListener) IListener {
	return &EventHandler{
		callback: cb,
		target:   target,
	}
}

func (h *EventHandler) Get() EventListener {
	return h.callback
}

func (h *EventHandler) Fire(event *Event) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
			debug.Backtrace("handle event "+event.Name, v, os.Stderr)
		}
	}()
	return h.callback(event)
}

type EventOnceHandler struct {
	callback EventListener
	target   *Emitter
	fired    bool
}

func NewEventOnceHandler(target *Emitter, cb EventListener) IListener {
	return &EventOnceHandler{
		callback: cb,
		target:   target,
	}
}

func (h *EventOnceHandler) Get() EventListener {
	return h.callback
}

func (h *EventOnceHandler) Fire(event *Event) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
			debug.Backtrace("handle event "+event.Name, v, os.Stderr)
		}
	}()
	if !h.fired {
		h.fired = true
		h.target.RemoveListener(event.Name, h.callback)
		return h.callback(event)
	}
	return ErrDuplicateOnce
}

var em = NewEmitter()

func On(eventName string, listener EventListener) {
	em.On(eventName, listener)
}

func AddListener(eventName string, listener EventListener) {
	em.AddListener(eventName, listener)
}

func Off(eventName string, listener EventListener) bool {
	return em.RemoveListener(eventName, listener)
}

func RemoveListener(eventName string, listener EventListener) bool {
	return em.RemoveListener(eventName, listener)
}

func PrependListener(eventName string, listener EventListener) {
	em.PrependListener(eventName, listener)
}

func Once(eventName string, listener EventListener) {
	em.Once(eventName, listener)
}

func AddOnceListener(eventName string, listener EventListener) {
	em.AddOnceListener(eventName, listener)
}

func PrependOnceListener(eventName string, listener EventListener) {
	em.PrependListener(eventName, listener)
}

func RemoveListeners(eventName string) {
	em.RemoveListeners(eventName)
}

func RemoveAllListeners() {
	em.RemoveAllListeners()
}

func Emit(eventName string, args ...any) error {
	return em.Emit(eventName, args...)
}
