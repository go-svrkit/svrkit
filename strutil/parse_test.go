// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package strutil

import (
	"reflect"
	"testing"
)

func TestParseToMapN(t *testing.T) {
	tests := []struct {
		input    string
		expected map[int32]int32
	}{
		{"", map[int32]int32{}},
		{"1=2", map[int32]int32{1: 2}},
		{"1=2|3=4", map[int32]int32{1: 2, 3: 4}},
		{"1=2|3=4|5=6", map[int32]int32{1: 2, 3: 4, 5: 6}},
	}
	for i, tc := range tests {
		d, err := ParseToMapN[int32, int32](tc.input)
		if err != nil {
			t.Fatalf("Test case %d failed: %v", i, err)
		}
		if len(d) == 0 && len(tc.expected) == 0 {
			continue
		}
		if !reflect.DeepEqual(d, tc.expected) {
			t.Fatalf("Test case %d failed, expect %v, got %v", i, tc.expected, d)
		}
	}
}

func TestParseKVPairs(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{"", map[string]string{}},
		{"a=", map[string]string{"a": ""}},
		{"a=1,b=2,c=3", map[string]string{"a": "1", "b": "2", "c": "3"}},
		{"a=1,b,c=3", map[string]string{"a": "1", "b": "", "c": "3"}},
		{"a=1,b=", map[string]string{"a": "1", "b": ""}},
		{"a=,b=2,c,d=3,e", map[string]string{"a": "", "b": "2", "c": "", "d": "3", "e": ""}},
		{"a='1,2,3',b=456", map[string]string{"a": "1,2,3", "b": "456"}},
		{"a=123, b=456", map[string]string{"a": "123", "b": "456"}},
		{"a = 123, b = 456", map[string]string{"a": "123", "b": "456"}},
		{"a=123 , b='4,5,6' , c = 789", map[string]string{"a": "123", "b": "4,5,6", "c": "789"}},
	}
	for i, test := range tests {
		r, err := ParseKVPairs[string, string](test.input)
		if err != nil {
			t.Fatalf("Test case %d failed: %v", i, err)
		}
		if len(r) == 0 && len(test.expected) == 0 {
			continue
		}
		if !reflect.DeepEqual(r, test.expected) {
			t.Fatalf("Test case %d failed, expect %v, got %v", i, test.expected, r)
		}
	}
}
