// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		key  string
		want string
	}{
		{"", ""},
		{"foo", "bar"},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if tt.key != "" {
				os.Setenv(tt.key, tt.want)
				defer os.Unsetenv(tt.key)
			}

			if got := GetEnv(tt.key, tt.want); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		key  string
		want int
	}{
		{"", 0},
		{"foo", 123},
		{"bar", math.MaxInt32},
		{"xyz", math.MinInt32},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if tt.key != "" {
				var val = strconv.Itoa(tt.want)
				os.Setenv(tt.key, val)
				defer os.Unsetenv(tt.key)
			}
			if got := GetEnvInt(tt.key, tt.want); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvInt64(t *testing.T) {
	tests := []struct {
		key  string
		want int64
	}{
		{"", 0},
		{"foo", 123},
		{"bar", math.MaxInt64},
		{"xyz", math.MinInt64},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if tt.key != "" {
				var val = strconv.FormatInt(tt.want, 10)
				os.Setenv(tt.key, val)
				defer os.Unsetenv(tt.key)
			}
			if got := GetEnvInt64(tt.key, tt.want); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvFloat(t *testing.T) {
	tests := []struct {
		key  string
		want float64
	}{
		{"", 0},
		{"foo", 123},
		{"bar", 3.14159},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if tt.key != "" {
				var val = strconv.FormatFloat(tt.want, 'f', 5, 64)
				os.Setenv(tt.key, val)
				defer os.Unsetenv(tt.key)
			}
			if got := GetEnvFloat(tt.key, tt.want); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"", false},
		{"false", false},
		{"true", true},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if tt.key != "" {
				var val = strconv.FormatBool(tt.want)
				os.Setenv(tt.key, val)
				defer os.Unsetenv(tt.key)
			}
			if got := GetEnvBool(tt.key); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
