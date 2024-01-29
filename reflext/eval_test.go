// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package reflext

import (
	"image"
	"reflect"
	"testing"
)

type TA struct {
	B        bool
	Str      string
	N        int64
	F        float64
	Pos      image.Point
	Rect     *image.Rectangle
	Triangle []image.Point
	Dict     map[int]string

	Next *TA
}

func (t *TA) GetArgs() map[int]string {
	return t.Dict
}

func (t *TA) Get(key int) string {
	return t.Dict[key]
}

func (t *TA) GetRect() *image.Rectangle {
	return t.Rect
}

func (t *TA) SetPos(x, y int) {
	t.Pos.X = x
	t.Pos.Y = y
}

func createTA() *TA {
	return &TA{
		B:   true,
		Str: "hello",
		N:   12345,
		F:   3.14159,
		Pos: image.Point{X: 10, Y: 20},
		Rect: &image.Rectangle{
			Min: image.Point{X: 500, Y: 150},
			Max: image.Point{X: 100, Y: 200},
		},
		Triangle: []image.Point{
			{X: 10, Y: 20},
			{X: 30, Y: 40},
			{X: 50, Y: 60},
		},
		Dict: map[int]string{
			123: "hello",
			456: "world",
		},
	}
}

func TestEvalGet(t *testing.T) {
	var ta = createTA()
	tests := []struct {
		expr         string
		shouldHasErr bool
		result       interface{}
	}{
		{"NotExist", true, nil},
		{"B", false, ta.B},
		{"Str", false, ta.Str},
		{"Pos.X", false, ta.Pos.X},
		{"Rect.Min.X", false, ta.Rect.Min.X},
		{"Rect.NotExist.X", true, nil},
		{"Rect.0.X", true, nil},
		{"Triangle[0].X", false, ta.Triangle[0].X},
		{"Triangle[3].X", true, nil},
		{"Triangle[`3`].X", true, nil},
		{"Triangle[kk].X", true, nil},
		{"Dict[123]", false, ta.Dict[123]},
		{"Dict[`123`]", true, nil},
		{"Dict[kkk]", true, nil},
		{"Dict[`kkk`]", true, nil},
	}
	for _, tc := range tests {
		var evalCtx = NewEvalContext(ta)
		evalCtx.ReadOnly = true
		v, err := evalCtx.Eval(tc.expr)
		if tc.shouldHasErr {
			if err == nil {
				t.Fatalf("%s: %v", tc.expr, err)
			}
			continue
		}
		t.Logf("%s: %v, %v", tc.expr, v, err)
		//evalCtx.PrintNodes(os.Stdout)
		if tc.result == nil && !IsInterfaceNil(v) {
			t.Fatalf("%s: %v", tc.expr, err)
		} else if !reflect.DeepEqual(v, tc.result) {
			t.Fatalf("%s: %v", tc.expr, err)
		}
	}
}

func TestEvalEval(t *testing.T) {
	var ta = createTA()
	tests := []struct {
		expr         string
		shouldHasErr bool
		result       interface{}
	}{
		{"GetXXX().A", true, false},
		{"GetRect().Min.X", false, ta.GetRect().Min.X},
		{"GetRect().Dx()", false, ta.GetRect().Dx()},
		{"GetRect().DDD()", true, ta.GetRect().Dx()},
		{"Get(123)", false, ta.Get(123)},
		{"Get(`123`)", false, ta.Get(123)},
		{"Get(`hello`)", true, nil},
	}
	for _, tc := range tests {
		var evalCtx = NewEvalContext(ta)
		v, err := evalCtx.Eval(tc.expr)
		if tc.shouldHasErr {
			if err == nil {
				t.Fatalf("%s: %v", tc.expr, err)
			}
			continue
		}
		t.Logf("%s: %v, %v", tc.expr, v, err)
		//evalCtx.PrintNodes(os.Stdout)
		if tc.result == nil && !IsInterfaceNil(v) {
			t.Fatalf("%s: %v", tc.expr, err)
		} else if !reflect.DeepEqual(v, tc.result) {
			t.Fatalf("%s: %v", tc.expr, err)
		}
	}
}

func TestEvalSet(t *testing.T) {
	tests := []struct {
		expr         string
		val          interface{}
		shouldHasErr bool
	}{
		{"B", false, false},
		{"Str", "woohaha", false},
		{"N", 98765, false},
		{"F", 1.618, false},
		{"Pos", `{"X": 11, "Y":22}`, false},
		{"Rect", `{"Min": {"X": 1, "Y":2}, "Max":{"X":3,"Y":4}}`, false},
		{"Triangle[0].X", 12345, false},
		{"Dict[123]", "123", false},
		{"GetRect().Min.X", 12345, false},
	}
	for _, tc := range tests {
		var obj = createTA()
		var ctx = NewEvalContext(obj)
		err := ctx.Set(tc.expr, tc.val)
		t.Logf("%s: %v", tc.expr, err)
		if tc.shouldHasErr {
			if err == nil {
				t.Fatalf("%s: %v", tc.expr, err)
			}
			continue
		}
		v, err := ctx.Eval(tc.expr)
		if err != nil {
			t.Fatalf("%s: %v", tc.expr, err)
		}
		if !reflect.DeepEqual(v, tc.val) {
			t.Fatalf("%s: %v", tc.expr, err)
		}
	}
}

func TestEvalRemove(t *testing.T) {

	tests := []struct {
		expr   string
		hasErr bool
	}{
		{"B", true},
		{"Str", true},
		{"Pos", true},
		{"Triangle[1]", false},
		{"Dict[`123`]", false},
	}
	for _, tc := range tests {
		var obj = createTA()
		var ctx = NewEvalContext(obj)
		err := ctx.Delete(tc.expr)
		t.Logf("%s: %v", tc.expr, err)
		if tc.hasErr {
			if err == nil {
				t.Fatalf("%s: %v", tc.expr, err)
			}
		} else {
			if err != nil {
				t.Fatalf("%s: %v", tc.expr, err)
			}
		}
	}
}
