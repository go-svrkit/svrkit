// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package strutil

import (
	"testing"
)

func TestWordCount(t *testing.T) {
	var cases = map[string]int{
		"one word: λ":             3,
		"中文":                      0,
		"你好，sekai！":               1,
		"oh, it's super-fancy!!a": 4,
		"":                        0,
		"-":                       0,
		"it's-'s":                 1,
	}
	for str, cnt := range cases {
		var n = WordCount(str)
		if n != cnt {
			t.Fatalf("%s is not %d length", str, n)
		}
	}
}

func TestRuneWidth(t *testing.T) {
	var cases = map[string]int{
		"a":    1,
		"中":    2,
		"\x11": 0,
	}
	for r, cnt := range cases {
		var n = RuneWidth([]rune(r)[0])
		if n != cnt {
			t.Fatalf("%s is not %d length", r, n)
		}
	}
}
