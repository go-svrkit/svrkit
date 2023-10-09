// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package events

import (
	"reflect"

	"go.uber.org/multierr"
)

// Emitter event subscription and publishing
// API modeled after node https://nodejs.org/docs/latest-v16.x/api/events.html#class-eventemitter
type Emitter struct {
	listeners map[string][]IListener
}

func NewEmitter() *Emitter {
	return &Emitter{
		listeners: make(map[string][]IListener),
	}
}

func (e *Emitter) EventNames() []string {
	var names = make([]string, 0, len(e.listeners))
	for name := range e.listeners {
		names = append(names, name)
	}
	return names
}

func (e *Emitter) ListenerCount(eventName string) int {
	return len(e.listeners[eventName])
}

func (e *Emitter) Listeners(eventName string) []IListener {
	return e.listeners[eventName]
}

// AddListener Adds the listener function to the end of the listeners array for the event named eventName.
// No checks are made to see if the listener has already been added. Multiple calls passing the same combination
// of eventName and listener will result in the listener being added, and called, multiple times.
func (e *Emitter) AddListener(eventName string, listener EventListener) {
	var h = NewEventHandler(e, listener)
	e.listeners[eventName] = append(e.listeners[eventName], h)
}

// PrependListener Adds the listener function to the beginning of the listeners array for the event named eventName.
// No checks are made to see if the listener has already been added. Multiple calls passing the same combination of
// eventName and listener will result in the listener being added, and called, multiple times.
func (e *Emitter) PrependListener(eventName string, listener EventListener) {
	var h = NewEventHandler(e, listener)
	var listeners = e.listeners[eventName]
	listeners = append([]IListener{h}, listeners...)
	e.listeners[eventName] = listeners
}

// RemoveListener will remove, at most, one instance of a listener from the listener array.
// If any single listener has been added multiple times to the listener array for the specified eventName,
// then RemoveListener() must be called multiple times to remove each instance.
func (e *Emitter) RemoveListener(eventName string, listener EventListener) bool {
	var fn = reflect.ValueOf(listener).Pointer()
	var listeners = e.listeners[eventName]
	for i, h := range listeners {
		var ptr = reflect.ValueOf(h.Get()).Pointer()
		if fn == ptr {
			copy(listeners[i:], listeners[i+1:])
			listeners[len(listeners)-1] = nil
			listeners = listeners[:len(listeners)-1]
			e.listeners[eventName] = listeners
			return true
		}
	}
	return false
}

// On No checks are made to see if the listener has already been added
func (e *Emitter) On(eventName string, listener EventListener) {
	e.AddListener(eventName, listener)
}

func (e *Emitter) Off(eventName string, listener EventListener) bool {
	return e.RemoveListener(eventName, listener)
}

// AddOnceListener Adds a one-time listener function for the event named eventName.
// The next time eventName is triggered, this listener is removed and then invoked.
func (e *Emitter) AddOnceListener(eventName string, listener EventListener) {
	var h = NewEventOnceHandler(e, listener)
	e.listeners[eventName] = append(e.listeners[eventName], h)
}

func (e *Emitter) PrependOnceListener(eventName string, listener EventListener) {
	var h = NewEventOnceHandler(e, listener)
	var listeners = e.listeners[eventName]
	listeners = append([]IListener{h}, listeners...)
	e.listeners[eventName] = listeners
}

func (e *Emitter) Once(eventName string, listener EventListener) {
	e.AddOnceListener(eventName, listener)
}

func (e *Emitter) RemoveListeners(eventName string) {
	delete(e.listeners, eventName)
}

func (e *Emitter) RemoveAllListeners() {
	e.listeners = make(map[string][]IListener)
}

// Emit synchronously calls each of the listeners registered for the event named eventName,
// in the order they were registered, passing the supplied arguments to each.
func (e *Emitter) Emit(eventName string, args ...any) error {
	var handlers = e.listeners[eventName]
	if len(handlers) == 0 {
		return ErrNoListenerFired
	}
	var event = NewEvent(e, eventName, args)
	var err error
	for _, h := range handlers {
		if er := h.Fire(event); er != nil {
			err = multierr.Append(err, er)
		}
	}
	return err
}
