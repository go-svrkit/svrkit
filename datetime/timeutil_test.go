// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package datetime

import (
	"testing"
	"time"
)

type testGedSuit1 struct {
	s1  string
	s2  string
	day int
}

func Test_GetNextMondayMs(t *testing.T) {
	const nextMonday = "2023-02-06T00:00:00Z"
	var tests = []struct {
		input    string
		expected string
	}{
		{"2023-01-30T00:00:01Z", nextMonday}, // 周一
		{"2023-01-31T00:00:00Z", nextMonday}, // 周二
		{"2023-02-01T00:00:00Z", nextMonday}, // 周三
		{"2023-02-02T00:00:00Z", nextMonday}, // 周四
		{"2023-02-03T12:00:00Z", nextMonday}, // 周五
		{"2023-02-04T12:00:00Z", nextMonday}, // 周六
		{"2023-02-05T23:59:59Z", nextMonday},
		{nextMonday, "2023-02-13T00:00:00Z"},
	}
	for i, tc := range tests {
		t1, err := time.Parse(time.RFC3339, tc.input)
		if err != nil {
			t.Fatalf("%v", err)
		}
		t2, err := time.Parse(time.RFC3339, tc.expected)
		if err != nil {
			t.Fatalf("%v", err)
		}
		var out = GetNextMonday(t1) // 下周一
		if t2.Unix() != out.Unix() {
			var str = out.Format(time.RFC3339)
			t.Fatalf("case #%d(%s) next monday error, %s != %s", i+1, tc.input, str, tc.expected)
		}
	}
}

func TestElapsedDaysBetween(t *testing.T) {
	var gedCases = []testGedSuit1{
		{"2018-01-01T00:00:00", "2018-01-01T23:59:59", 0},
		{"2018-01-01T00:00:00", "2018-02-01T00:00:00", 31},
		{"2016-02-28T00:00:00", "2016-03-01T00:00:00", 2},
		{"2018-01-01T00:00:00", "2017-12-30T00:00:00", -2},
	}
	for _, item := range gedCases {
		var t1, _ = time.Parse(DateLayout, item.s1)
		var t2, _ = time.Parse(DateLayout, item.s2)
		var d = ElapsedDaysBetween(t1, t2)
		if d != item.day {
			t.Fatalf("%s, %s, %d != %d", item.s1, item.s2, d, item.day)
		}
	}
}

func TestMidnightTimeOf(t *testing.T) {
	var cases = []testGedSuit1{
		{"2018-01-01T00:00:00", "2018-01-01T00:00:00", 0},
		{"2018-02-01T12:34:56", "2018-02-01T00:00:00", 0},
	}
	for _, item := range cases {
		var t1, _ = time.Parse(DateLayout, item.s1)
		var t2, _ = time.Parse(DateLayout, item.s2)
		var tm = MidnightTimeOf(t1)
		if tm != t2 {
			t.Fatalf("%s, %s, %v != %v", item.s1, item.s2, tm, t2)
		}
	}
}

func TestThisMomentAfterDays(t *testing.T) {
	var cases = []testGedSuit1{
		{"2018-01-01T00:00:00", "2018-01-01T00:00:00", 0},
		{"2018-02-28T12:34:56", "2018-03-01T12:34:56", 1},
		{"2017-02-28T12:34:56", "2018-02-28T12:34:56", 365},
		{"2000-02-28T12:34:56", "2000-03-01T12:34:56", 2},
		{"2016-03-01T12:34:56", "2016-02-28T12:34:56", -2},
	}
	for _, item := range cases {
		var t1, _ = time.Parse(DateLayout, item.s1)
		var t2, _ = time.Parse(DateLayout, item.s2)
		var tm = ThisMomentAfterDays(t1, item.day)
		if tm != t2 {
			t.Fatalf("%s, %s, %v != %v", item.s1, item.s2, tm, t2)
		}
	}
}

func TestEndOfWeek(t *testing.T) {
	var cases = []testGedSuit1{
		{"2021-02-02T12:23:34", "2021-02-06T23:59:59", 0}, // Tue
	}
	for _, item := range cases {
		var t1, _ = time.Parse(DateLayout, item.s1)
		var t2, _ = time.Parse(DateLayout, item.s2)
		var tm = EndOfWeek(t1)
		if tm != t2 {
			t.Fatalf("end week of %s alias %s, not %s", item.s1, t2.Format(time.RFC3339), tm.Format(time.RFC3339))
		}
	}
}
