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

var gem = NewEmitter()

func On(eventName string, listener EventListener) {
	gem.On(eventName, listener)
}

func AddListener(eventName string, listener EventListener) {
	gem.AddListener(eventName, listener)
}

func Off(eventName string, listener EventListener) bool {
	return gem.RemoveListener(eventName, listener)
}

func RemoveListener(eventName string, listener EventListener) bool {
	return gem.RemoveListener(eventName, listener)
}

func PrependListener(eventName string, listener EventListener) {
	gem.PrependListener(eventName, listener)
}

func Once(eventName string, listener EventListener) {
	gem.Once(eventName, listener)
}

func AddOnceListener(eventName string, listener EventListener) {
	gem.AddOnceListener(eventName, listener)
}

func PrependOnceListener(eventName string, listener EventListener) {
	gem.PrependListener(eventName, listener)
}

func RemoveListeners(eventName string) {
	gem.RemoveListeners(eventName)
}

func RemoveAllListeners() {
	gem.RemoveAllListeners()
}

func Emit(eventName string, args ...any) error {
	return gem.Emit(eventName, args...)
}
