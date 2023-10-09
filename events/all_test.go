// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package events

import (
	"testing"
)

func TestEmitter_AddListener(t *testing.T) {
	var fired int
	AddListener("error", func(event *Event) error {
		fired++
		return nil
	})
	if err := Emit("error"); err != nil {
		t.Fatalf("%v", err)
	}
	if fired != 1 {
		t.Fatalf("event not fired")
	}
	if err := Emit("error"); err != nil {
		t.Fatalf("%v", err)
	}
	if fired != 2 {
		t.Fatalf("event not fired")
	}
}

func TestEmitter_PrependListener(t *testing.T) {
	var a int
	var b int
	PrependListener("error", func(event *Event) error {
		if b <= a {
			t.Fatalf("prepend failed")
		}
		a++
		return nil
	})
	PrependListener("error", func(event *Event) error {
		b++
		if b <= a {
			t.Fatalf("prepend failed")
		}
		return nil
	})
	if err := Emit("error"); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestEmitter_AddOnceListener(t *testing.T) {
	var fired int
	AddOnceListener("error", func(event *Event) error {
		fired++
		return nil
	})
	if err := Emit("error"); err != nil {
		t.Fatalf("%v", err)
	}
	if fired != 1 {
		t.Fatalf("event not fired")
	}
	if err := Emit("error"); err != ErrNoListenerFired {
		t.Fatalf("%v", err)
	}
	if fired != 1 {
		t.Fatalf("event not fired")
	}
}
