// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package events

import (
	"gopkg.in/svrkit.v1/gutil"
)

type Event struct {
	Name   string
	Target IEventTarget
	Args   []any
}

func NewEvent(target IEventTarget, name string, args ...any) *Event {
	return &Event{
		Target: target,
		Name:   name,
		Args:   args,
	}
}

type EventListener func(event *Event) error

type IListener interface {
	Get() EventListener
	Fire(*Event) error
}

type IEventTarget interface {
	AddListener(eventName string, listener EventListener)
	RemoveListener(eventName string, listener EventListener) bool
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
	defer gutil.CatchPanic(event.Name)
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
	defer gutil.CatchPanic(event.Name)

	if !h.fired {
		h.fired = true
		h.target.RemoveListener(event.Name, h.callback)
		return h.callback(event)
	}
	return ErrDuplicateOnce
}
