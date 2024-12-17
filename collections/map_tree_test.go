// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package collections

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
	"unsafe"
)

func randText(n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b = make([]byte, n)
	for i := range b {
		idx := rand.Int() % len(alphabet)
		b[i] = alphabet[idx]
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func TestTreeMap_Get(t *testing.T) {
	tree := NewTreeMapCmp[int, string]()

	if actualValue := tree.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	if actualValue := tree.GetEntry(2).Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	tree.Put(1, "x") // 1->x
	tree.Put(2, "b") // 1->x, 2->b (in order)
	tree.Put(1, "a") // 1->a, 2->b (in order, replacement)
	tree.Put(3, "c") // 1->a, 2->b, 3->c (in order)
	tree.Put(4, "d") // 1->a, 2->b, 3->c, 4->d (in order)
	tree.Put(5, "e") // 1->a, 2->b, 3->c, 4->d, 5->e (in order)
	tree.Put(6, "f") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f (in order)

	fmt.Println(tree)
	//
	//  RBTree
	//  │           ┌── 6
	//  │       ┌── 5
	//  │   ┌── 4
	//  │   │   └── 3
	//  └── 2
	//       └── 1

	if actualValue := tree.Size(); actualValue != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}

	if actualValue := tree.GetEntry(4).Size(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}

	if actualValue := tree.GetEntry(2).Size(); actualValue != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}

	if actualValue := tree.GetEntry(8).Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestTreeMap_Put(t *testing.T) {
	tree := NewTreeMapCmp[int, string]()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite

	if actualValue := tree.Size(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}
	if actualValue, expectedValue := tree.Keys(), []int{1, 2, 3, 4, 5, 6, 7}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := tree.Values(), []string{"a", "b", "c", "d", "e", "f", "g"}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tests1 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue, actualFound := tree.Get(test[0].(int))
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}
}

func TestTreeMap_Remove(t *testing.T) {
	tree := NewTreeMapCmp[int, string]()
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite

	tree.Remove(5)
	tree.Remove(6)
	tree.Remove(7)
	tree.Remove(8)
	tree.Remove(5)

	if actualValue, expectedValue := tree.Keys(), []int{1, 2, 3, 4}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := tree.Values(), []string{"a", "b", "c", "d"}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := tree.Size(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	tests2 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "", false},
		{6, "", false},
		{7, "", false},
		{8, "", false},
	}

	for _, test := range tests2 {
		actualValue, actualFound := tree.Get(test[0].(int))
		if actualValue != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}

	tree.Remove(1)
	tree.Remove(4)
	tree.Remove(2)
	tree.Remove(3)
	tree.Remove(2)
	tree.Remove(2)

	if actualValue, expectedValue := tree.Keys(), []int{}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := tree.Values(), []string{}; !slices.Equal(actualValue, expectedValue) {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if empty, size := tree.IsEmpty(), tree.Size(); empty != true || size != -0 {
		t.Errorf("Got %v expected %v", empty, true)
	}
}

func TestTreeMap_FirstAndLast(t *testing.T) {
	tree := NewTreeMapCmp[int, string]()
	if actualValue := tree.FirstEntry(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := tree.LastEntry(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	tree.Put(1, "a")
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x") // overwrite
	tree.Put(2, "b")

	if actualValue, expectedValue := tree.FirstEntry().Key, 1; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := tree.FirstEntry().Value, "x"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := tree.LastEntry().Key, 7; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := tree.LastEntry().Value, "g"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestTreeMap_CeilingAndFloor(t *testing.T) {
	tree := NewTreeMapCmp[int, string]()
	if node := tree.FloorEntry(0); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
	if node := tree.CeilingEntry(0); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")

	if node := tree.FloorEntry(4); node.Key != 4 {
		t.Errorf("Got %v expected %v", node.Key, 4)
	}
	if node := tree.FloorEntry(0); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	if node := tree.CeilingEntry(4); node.Key != 4 {
		t.Errorf("Got %v expected %v", node.Key, 4)
	}
	if node := tree.CeilingEntry(8); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
}
