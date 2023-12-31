// Copyright © 2018 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package strutil

import (
	"testing"
)

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
