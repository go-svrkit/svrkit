package datetime

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	for i, tt := range []struct {
		dur      string
		expected time.Duration
	}{
		{"1h", time.Hour},
		{"1m", time.Minute},
		{"1s", time.Second},
		{"1ms", time.Millisecond},
		{"1µs", time.Microsecond},
		{"1us", time.Microsecond},
		{"1ns", time.Nanosecond},
		{"4.000000001s", 4*time.Second + time.Nanosecond},
		{"1h0m4.000000001s", time.Hour + 4*time.Second + time.Nanosecond},
		{"1h1m0.01s", 61*time.Minute + 10*time.Millisecond},
		{"1h1m0.123456789s", 61*time.Minute + 123456789*time.Nanosecond},
		{"1.00002ms", time.Millisecond + 20*time.Nanosecond},
		{"1.00000002s", time.Second + 20*time.Nanosecond},
		{"693ns", 693 * time.Nanosecond},
		{"10s1us693ns", 10*time.Second + time.Microsecond + 693*time.Nanosecond},

		{"1ms1ns", time.Millisecond + 1*time.Nanosecond},
		{"1s20ns", time.Second + 20*time.Nanosecond},
		{"60h8ms", 60*time.Hour + 8*time.Millisecond},
		{"96h63s", 96*time.Hour + 63*time.Second},

		{"2d3s96ns", 48*time.Hour + 3*time.Second + 96*time.Nanosecond},
		{"9d3s96ns", 168*time.Hour + 48*time.Hour + 3*time.Second + 96*time.Nanosecond},
		{"9d3s3µs96ns", 168*time.Hour + 48*time.Hour + 3*time.Second + 3*time.Microsecond + 96*time.Nanosecond},
	} {
		d, err := ParseDuration(tt.dur)
		if err != nil {
			t.Logf("index %d -> in: %s returned: %s\tnot equal to %s", i, tt.dur, err.Error(), tt.expected.String())

		} else if tt.expected != d {
			t.Errorf("index %d -> in: %s returned: %s\tnot equal to %s", i, tt.dur, d.String(), tt.expected.String())
		}
	}
}

func TestPrettyDuration(t *testing.T) {
	for i, tt := range []struct {
		dur      time.Duration
		expected string
	}{
		{0, "0s"},
		{time.Hour, "1h"},
		{time.Minute, "1m"},
		{time.Second, "1s"},
		{time.Millisecond, "1ms"},
		{time.Microsecond, "1us"},
		{time.Nanosecond, "1ns"},
		{4*time.Second + time.Nanosecond, "4s1ns"},
		{time.Hour + 4*time.Second + time.Nanosecond, "1h4s1ns"},
		{61*time.Minute + 10*time.Millisecond, "1h1m10ms"},
		{61*time.Minute + 123456789*time.Nanosecond, "1h1m123ms456us789ns"},
		{time.Millisecond + 20*time.Nanosecond, "1ms20ns"},
		{time.Second + 20*time.Nanosecond, "1s20ns"},
		{693 * time.Nanosecond, "693ns"},
		{10*time.Second + time.Microsecond + 693*time.Nanosecond, "10s1us693ns"},
		{time.Millisecond + 1*time.Nanosecond, "1ms1ns"},
		{time.Second + 20*time.Nanosecond, "1s20ns"},
		{60*time.Hour + 8*time.Millisecond, "2d12h8ms"},
		{96*time.Hour + 63*time.Second, "4d1m3s"},
		{48*time.Hour + 3*time.Second + 96*time.Nanosecond, "2d3s96ns"},
		{168*time.Hour + 48*time.Hour + 3*time.Second + 96*time.Nanosecond, "9d3s96ns"},
		{168*time.Hour + 48*time.Hour + 3*time.Second + 3*time.Microsecond + 96*time.Nanosecond, "9d3s3us96ns"},

		{2540400*time.Hour + 10*time.Minute + 10*time.Second, "105850d10m10s"},

		{9_223_372_036_854_775_807 * time.Nanosecond, "106751d23h47m16s854ms775us807ns"},

		{-9_223_372_036_854_775_807 * time.Nanosecond, "-106751d23h47m16s854ms775us807ns"},

		{-time.Hour, "-1h"},
		{-time.Minute, "-1m"},
		{-time.Second, "-1s"},
		{-time.Millisecond, "-1ms"},
		{-time.Microsecond, "-1us"},
		{-time.Nanosecond, "-1ns"},
		{-4*time.Second - time.Nanosecond, "-4s1ns"},
	} {
		s := PrettyDuration(tt.dur)
		if tt.expected != s {
			t.Errorf("index %d -> in: %s returned: %s\tnot equal to %s", i, tt.dur, s, tt.expected)
		}

		d, _ := ParseDuration(s)
		if d != tt.dur {
			t.Errorf("error converting string to duration: index %d -> in: %s returned: %s", i, tt.dur, d)
		}
	}
}
