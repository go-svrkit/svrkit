// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package datetime

import (
	"time"
)

const (
	DateLayout      = "2006-01-02T15:04:05"
	TimestampLayout = "2006-01-02T15:04:05.999"

	MsPerSecond = 1000
	MsPerMinute = 60 * MsPerSecond
	MsPerHour   = 60 * MsPerMinute
	MsPerDay    = 24 * MsPerHour
	MsPerWeek   = 7 * MsPerDay

	SecondsPerMin  = 60
	SecondsPerHour = 60 * SecondsPerMin
	SecondsPerDay  = 24 * SecondsPerHour
)

var (
	FirstDayIsMonday = false
	DefaultLoc       = time.UTC
)

func FormatTime(t time.Time) string {
	return t.Format(TimestampLayout)
}

func FormatMsTime(ms int64) string {
	return FormatTime(Ms2Time(ms))
}

func Ms2Time(ms int64) time.Time {
	var t time.Time
	if ms == 0 {
		return t
	}
	var sec = ms / MsPerSecond
	var nsec = (ms % MsPerSecond) * int64(time.Millisecond)
	return time.Unix(sec, nsec).UTC()
}

// 是否闰年
func IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

// 当日零点
func MidnightTimeOf(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 下一个凌晨
func NextMidnight(ts int64) int64 {
	var midTime = MidnightTimeOf(Ms2Time(ts))
	return midTime.UnixNano()/int64(time.Millisecond) + MsPerDay
}

// N天后的这个时候
func ThisMomentAfterDays(this time.Time, days int) time.Time {
	if days == 0 {
		return this
	}
	return this.Add(time.Duration(days) * time.Hour * 24)
}

// 本周的起点
func StartingOfWeek(t time.Time) time.Time {
	t2 := MidnightTimeOf(t)
	weekday := int(t2.Weekday())
	if FirstDayIsMonday {
		if weekday == 0 {
			weekday = 7
		}
		weekday = weekday - 1
	}
	d := time.Duration(-weekday) * 24 * time.Hour
	return t2.Add(d)
}

// 本周的最后一天
func EndOfWeek(t time.Time) time.Time {
	begin := StartingOfWeek(t)
	end := ThisMomentAfterDays(begin, 7)
	return end.Add(-time.Second) // 23:59:59
}

// 年度第一天
func FirstDayOfYear(year int) time.Time {
	return time.Date(year, 1, 1, 0, 0, 0, 0, DefaultLoc)
}

// 年度最后一天
func LastDayOfYear(year int) time.Time {
	return time.Date(year, 12, 31, 0, 0, 0, 0, DefaultLoc)
}

// 获取两个时间中经过的天数
func ElapsedDaysBetween(start, end time.Time) int {
	var negative = false
	if start.After(end) {
		start, end = end, start
		negative = true
	}
	var days = 0
	if start.Year() != end.Year() {
		t := LastDayOfYear(start.Year())
		days = t.YearDay() - start.YearDay() // start年份的天数
		for i := start.Year() + 1; i < end.Year(); i++ {
			var t1 = LastDayOfYear(i)
			days += t1.YearDay() // start-end中间每年的天数
		}
		days += end.YearDay() // end年份的天数
	} else {
		days = end.YearDay() - start.YearDay()
	}
	if negative {
		days = -days
	}
	return days
}

// 一个月的天数
func DaysCountOfMonth(year, month int) int {
	switch time.Month(month) {
	case time.January:
		return 31
	case time.February:
		if year > 0 && IsLeapYear(year) {
			return 29
		}
		return 28
	case time.March:
		return 31
	case time.April:
		return 30
	case time.May:
		return 31
	case time.June:
		return 30
	case time.July:
		return 31
	case time.August:
		return 31
	case time.September:
		return 30
	case time.October:
		return 31
	case time.November:
		return 30
	case time.December:
		return 31
	}
	return 0
}
