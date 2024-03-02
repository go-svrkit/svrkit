// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package maps

import (
	"fmt"
	"strings"
	"testing"
)

func createTreeMap2() *TreeMap[int, string] {
	var m = NewOrderedTreeMap[int, string]()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	// │   ┌── 7
	// └── 6
	//     │   ┌── 5
	//     └── 4
	//         │   ┌── 3
	//         └── 2
	//             └── 1
	return m
}

func TestTreeMapIterator(t *testing.T) {
	var m = createTreeMap2()

	var count = 0
	var sb1 strings.Builder
	var sb2 strings.Builder

	var iter = m.Iterator()
	for iter.HasNext() {
		count++
		var entry = iter.Next()
		fmt.Fprintf(&sb1, "%v", entry.GetKey())
		fmt.Fprintf(&sb2, "%v", entry.GetValue())
	}
	if actualValue, expectedValue := sb1.String(), "1234567"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := sb2.String(), "abcdefg"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := count, m.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestTreeMapDescendingIterator(t *testing.T) {
	var m = createTreeMap2()

	var count = 0
	var sb1 strings.Builder
	var sb2 strings.Builder

	var iter = m.DescendingIterator()
	for iter.HasNext() {
		count++
		var entry = iter.Next()
		fmt.Fprintf(&sb1, "%v", entry.GetKey())
		fmt.Fprintf(&sb2, "%v", entry.GetValue())
	}
	if actualValue, expectedValue := sb1.String(), "7654321"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := sb2.String(), "gfedcba"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := count, m.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}
