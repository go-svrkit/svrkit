package datetime

import (
	"testing"
	"time"
)

type DurationPairs struct {
	Text     string
	Duration time.Duration
}

var testDurationCases = []DurationPairs{
	{"1h", time.Hour},
	{"1m", time.Minute},
	{"1s", time.Second},
	{"1ms", time.Millisecond},
	{"1µs", time.Microsecond},
	{"1ns", time.Nanosecond},
	{"4.000000001s", 4*time.Second + time.Nanosecond},
	{"1h4.000000001s", time.Hour + 4*time.Second + time.Nanosecond},
	{"1h1m0.01s", 61*time.Minute + 10*time.Millisecond},
	{"1h1m0.123456789s", 61*time.Minute + 123456789*time.Nanosecond},
	{"1.00002ms", time.Millisecond + 20*time.Nanosecond},
	{"1.00000002s", time.Second + 20*time.Nanosecond},
	{"1h1m0.123456789s", 61*time.Minute + 123456789*time.Nanosecond},
	{"693ns", 693 * time.Nanosecond},
	{"10.000001693s", 10*time.Second + time.Microsecond + 693*time.Nanosecond},

	{"1.000001ms", time.Millisecond + 1*time.Nanosecond},
	{"1.00000002s", time.Second + 20*time.Nanosecond},
	{"2d12h0.008s", 60*time.Hour + 8*time.Millisecond},
	{"4d1m3s", 96*time.Hour + 63*time.Second},

	{"2d3.000000096s", 48*time.Hour + 3*time.Second + 96*time.Nanosecond},
	{"9d3.000000096s", 168*time.Hour + 48*time.Hour + 3*time.Second + 96*time.Nanosecond},
	{"9d3.000003096s", 168*time.Hour + 48*time.Hour + 3*time.Second + 3*time.Microsecond + 96*time.Nanosecond},
	{"106751d23h47m16.854775807s", 106751*Day + 23*time.Hour + 47*time.Minute + 16*time.Second + 854*time.Millisecond + 775*time.Microsecond + 807*time.Nanosecond},

	{"-106751d23h47m16.854775807s", -9_223_372_036_854_775_807 * time.Nanosecond},

	{"-1h", -time.Hour},
	{"-1m", -time.Minute},
	{"-1s", -time.Second},
	{"-1ms", -time.Millisecond},
	{"-1µs", -time.Microsecond},
	{"-1ns", -time.Nanosecond},
	{"-4.000000001s", -4*time.Second - time.Nanosecond},
}

func TestParseDuration(t *testing.T) {
	for i, tt := range testDurationCases {
		d, err := ParseDuration(tt.Text)
		if err != nil {
			t.Logf("index %d -> in: %s returned: %v\tnot equal to %s", i, tt.Text, err, tt.Duration.String())

		} else if tt.Duration != d {
			t.Errorf("index %d -> in: %s returned: %s\tnot equal to %s", i, tt.Text, d.String(), tt.Duration.String())
		}
	}
}

func TestPrettyDuration(t *testing.T) {
	for i, tt := range testDurationCases {
		s := PrettyDuration(tt.Duration)
		if tt.Text != s {
			t.Errorf("index %d -> in: %s returned: %s\tnot equal to %s", i, tt.Duration, s, tt.Text)
		}
	}
}
