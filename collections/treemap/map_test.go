// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package treemap

import (
	"fmt"
	"strings"
	"testing"

	"gopkg.in/svrkit.v1/collections/util"
)

func createTreeMap() *Map[int, string] {
	var m = New[int, string](util.OrderedCmp[int])
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite
	return m
}

func mapKeysText(m *Map[int, string]) string {
	var sb strings.Builder
	for _, key := range m.Keys() {
		fmt.Fprintf(&sb, "%v", key)
	}
	return sb.String()
}

func mapValuesText(m *Map[int, string]) string {
	var sb strings.Builder
	for _, val := range m.Values() {
		fmt.Fprintf(&sb, "%v", val)
	}
	return sb.String()
}

func checkMapKeyValue(t *testing.T, m *Map[int, string], keyS, valueS string, size int) {
	if actualValue := m.Size(); actualValue != size {
		t.Errorf("Got %v expected %v", actualValue, size)
	}
	if actualValue, expectedValue := mapKeysText(m), keyS; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := mapValuesText(m), valueS; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestTreeMapPut(t *testing.T) {
	var m = createTreeMap()

	checkMapKeyValue(t, m, "1234567", "abcdefg", 7)

	tests := []struct {
		key   int
		value string
		found bool
	}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}
	for _, tc := range tests {
		// retrievals
		actualValue, actualFound := m.Get(tc.key)
		if actualValue != tc.value || actualFound != tc.found {
			t.Errorf("key %v got %v expected %v", tc.key, actualValue, tc.value)
		}
	}
}

func TestTreeMapRemove(t *testing.T) {
	var m = createTreeMap()
	for i := 5; i <= 8; i++ {
		m.Remove(i)
	}
	m.Remove(5) // remove again

	checkMapKeyValue(t, m, "1234", "abcd", 4)

	tests := []struct {
		key   int
		value string
		found bool
	}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "", false},
		{6, "", false},
		{7, "", false},
		{8, "", false},
	}
	for _, tc := range tests {
		// retrievals
		actualValue, actualFound := m.Get(tc.key)
		if actualValue != tc.value || actualFound != tc.found {
			t.Errorf("key %v got %v expected %v", tc.key, actualValue, tc.value)
		}
	}

	m.Clear()
	checkMapKeyValue(t, m, "", "", 0)
}

func TestTreeMapFirstLast(t *testing.T) {
	var m = New[int, string](util.OrderedCmp[int])
	if actualValue, found := m.FirstKey(); found {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue, found := m.LastKey(); found {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	m.Put(1, "a")
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x") // overwrite
	m.Put(2, "b")

	firstKey, _ := m.FirstKey()
	lastKey, _ := m.LastKey()
	firstVal, _ := m.Get(firstKey)
	lastVal, _ := m.Get(lastKey)

	if actualValue, expectedValue := fmt.Sprintf("%v", firstKey), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%v", firstVal), "x"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := fmt.Sprintf("%v", lastKey), "7"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%v", lastVal), "g"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestTreeMapCeilingAndFloor(t *testing.T) {
	var m = New[int, string](util.OrderedCmp[int])

	if entry := m.FloorEntry(0); entry != nil {
		t.Errorf("Got %v expected %v", entry, "<nil>")
	}
	if entry := m.CeilingEntry(0); entry != nil {
		t.Errorf("Got %v expected %v", entry, "<nil>")
	}

	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")

	if node := m.FloorEntry(4); node.GetKey() != 4 {
		t.Errorf("Got %v expected %v", node.GetKey(), 4)
	}
	if node := m.FloorEntry(0); node != nil {
		t.Errorf("Got %v expected %v", node.GetKey(), "<nil>")
	}

	if node := m.CeilingEntry(4); node.GetKey() != 4 {
		t.Errorf("Got %v expected %v", node.GetKey(), 4)
	}
	if node := m.CeilingEntry(8); node != nil {
		t.Errorf("Got %v expected %v", node.GetKey(), "<nil>")
	}
}

func createTreeMap2() *Map[int, string] {
	var m = New[int, string](util.OrderedCmp[int])
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
