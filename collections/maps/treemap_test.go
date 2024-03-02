// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package maps

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/collections/util"
	"gopkg.in/svrkit.v1/strutil"
)

func createTreeMap(text string) *TreeMap[int, string] {
	keys, values := strutil.ParseKeyValues[int, string](text, "=", ",")
	var m = NewTreeMap[int, string](util.OrderedCmp[int])
	for i, k := range keys {
		m.Put(k, values[i])
	}
	return m
}

func formatTreeMap(m *TreeMap[int, string]) string {
	var sb strings.Builder
	var entry = m.getFirstEntry()
	for entry != nil {
		fmt.Fprintf(&sb, "%v=%v", entry.key, entry.value)
		entry = successor(entry)
		if entry != nil {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

func TestTreeMapPut(t *testing.T) {
	var s = "1=a,2=b,3=c,4=d,5=e,6=f,7=g"

	tests := []struct {
		input string
		want  string
	}{
		{"1=a,2=b,3=c,4=d,5=e,6=f,7=g", s},
		{"1=a,4=d,5=e,6=f,2=b,3=c,7=g", s},
		{"1=a,2=b,,6=f,3=c,4=d,5=e7=g", s},
		{"4=d,5=e,1=a,6=f,7=g,2=b,3=c", s},
	}
	for i, tc := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var m = createTreeMap(s)
			var out = formatTreeMap(m)
			assert.Equal(t, tc.want, out)
		})
	}
}

func TestTreeMapGet(t *testing.T) {
	var s = "1=a,2=b,3=c,4=d,5=e,6=f,7=g"
	var m = createTreeMap(s)
	assert.Equal(t, formatTreeMap(m), s)

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
	var m = createTreeMap("")
	for i := 5; i <= 8; i++ {
		m.Remove(i)
	}
	m.Remove(5) // remove again

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
}

func TestTreeMapFirstLast(t *testing.T) {
	var m = NewTreeMap[int, string](util.OrderedCmp[int])
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
	var m = NewTreeMap[int, string](util.OrderedCmp[int])

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
