package strutil

import (
	"reflect"
	"testing"
)

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
	const sep1, sep2 = ',', '='
	for i, test := range tests {
		r := ParseKVPairs(test.input, sep1, sep2)
		if len(r) == 0 && len(test.expected) == 0 {
			continue
		}
		if !reflect.DeepEqual(r, test.expected) {
			t.Fatalf("Test case %d failed, expect %v, got %v", i, test.expected, r)
		}
	}
}
