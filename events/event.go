// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package events

import (
	"errors"
)

var (
	ErrDuplicateOnce   = errors.New("duplicate once call")
	ErrNoListenerFired = errors.New("no listener fired")
)

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
