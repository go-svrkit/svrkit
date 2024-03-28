// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package strutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	var s = RandString(12)
	assert.Equal(t, 12, len(s))
}

func TestFindFirstDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", -1},
		{"abc", -1},
		{"123", 0},
		{"abc123", 3},
	}
	for i, test := range tests {
		output := FindFirstDigit(test.input)
		if test.expected != output {
			t.Fatalf("Test case %d failed, expect %d, got %d", i, test.expected, output)
		}
	}
}

func TestFindFirstNonDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", -1},
		{"123", -1},
		{"abc", 0},
		{"123abc", 3},
	}
	for i, test := range tests {
		output := FindFirstNonDigit(test.input)
		if test.expected != output {
			t.Fatalf("Test case %d failed, expect %d, got %d", i, test.expected, output)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"abc", "cba"},
		{"a", "a"},
		{"çınar", "ranıç"},
		{"    yağmur", "rumğay    "},
		{"επαγγελματίες", "ςείταμλεγγαπε"},
	}

	for i, test := range tests {
		output := Reverse(test.input)
		if test.expected != output {
			t.Fatalf("Test case %d failed, expect %s, got %s", i, test.expected, output)
		}
	}
}

func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		input1   string
		input2   string
		expected string
	}{
		{"", "a", ""},
		{"ab", "cd", ""},
		{"ab123", "abc456", "ab"},
	}
	for i, test := range tests {
		output := LongestCommonPrefix(test.input1, test.input2)
		if test.expected != output {
			t.Fatalf("Test case %d failed, expect %s, got %s", i, test.expected, output)
		}
	}
}

func TestPrettyBytes(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "0B"},
		{KiB, "1KiB"},
		{-KiB, "-1KiB"},
		{KiB + 100, "1.1KiB"},
		{MiB, "1MiB"},
		{MiB + 10*KiB, "1.01MiB"},
		{GiB, "1GiB"},
		{GiB + 100*MiB, "1.098GiB"},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := PrettyBytes(tt.input); got != tt.want {
				t.Errorf("PrettyBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseByteCount(t *testing.T) {
	tests := []struct {
		input  string
		want   int64
		wantOK bool
	}{
		{"", 0, true},
		{"0", 0, true},
		{"0B", 0, false},
		{"64B", 64, false},
		{"1KiB", KiB, true},
		{"1MiB", MiB, true},
		{"1GiB", GiB, true},
		{"1TiB", TiB, true},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i)
		t.Run(name, func(t *testing.T) {
			got, ok := ParseByteCount(tt.input)
			assert.Equal(t, tt.wantOK, ok)
			if ok {
				assert.Equalf(t, tt.want, got, "ParseByteCount(%v)", tt.input)
			}
		})
	}
}
