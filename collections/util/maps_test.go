package util

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKeys(t *testing.T) {
	tests := []struct {
		m    map[int]int
		want []int
	}{
		{map[int]int{}, []int{}},
		{map[int]int{1: 1}, []int{1}},
		{map[int]int{3: 3, 2: 2, 1: 1}, []int{3, 2, 1}},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var out = MapKeys(tt.m)
			assert.True(t, slices.Equal(tt.want, out))
		})
	}
}

func TestMapValues(t *testing.T) {
	tests := []struct {
		m    map[int]string
		want []string
	}{
		{map[int]string{}, []string{}},
		{map[int]string{1: "a"}, []string{"a"}},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []string{"a", "b", "c"}},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var out = MapValues(tt.m)
			assert.True(t, slices.Equal(tt.want, out))
		})
	}
}

func TestMapKeyValues(t *testing.T) {
	tests := []struct {
		m     map[int]string
		want1 []int
		want2 []string
	}{
		{map[int]string{}, nil, nil},
		{map[int]string{1: "a"}, []int{1}, []string{"a"}},
		{map[int]string{1: "a", 2: "b", 3: "c"}, []int{1, 2, 3}, []string{"a", "b", "c"}},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			out1, out2 := MapKeyValues(tt.m)
			assert.True(t, slices.Equal(tt.want1, out1))
			assert.True(t, slices.Equal(tt.want2, out2))
		})
	}
}

func TestMapOrderedKeys(t *testing.T) {
	tests := []struct {
		m    map[int]string
		want []int
	}{
		{map[int]string{}, nil},
		{map[int]string{1: "a"}, []int{1}},
		{map[int]string{3: "a", 2: "b", 1: "c"}, []int{1, 2, 3}},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			out := MapOrderedKeys(tt.m)
			assert.True(t, slices.Equal(tt.want, out))
		})
	}
}

func TestMapOrderedKeyValues(t *testing.T) {
	tests := []struct {
		m     map[int]string
		want1 []int
		want2 []string
	}{
		{map[int]string{}, nil, nil},
		{map[int]string{1: "a"}, []int{1}, []string{"a"}},
		{map[int]string{3: "a", 2: "b", 1: "c"}, []int{1, 2, 3}, []string{"c", "b", "a"}},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			out1, out2 := MapOrderedKeyValues(tt.m)
			assert.True(t, slices.Equal(tt.want1, out1))
			assert.True(t, slices.Equal(tt.want2, out2))
		})
	}
}

func TestMapOrderedValues(t *testing.T) {
	tests := []struct {
		m    map[int]string
		want []string
	}{
		{map[int]string{}, nil},
		{map[int]string{1: "a"}, []string{"a"}},
		{map[int]string{3: "a", 2: "b", 1: "c"}, []string{"c", "b", "a"}},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			out := MapOrderedValues(tt.m)
			assert.True(t, slices.Equal(tt.want, out))
		})
	}
}

func TestMapUnion(t *testing.T) {
	tests := []struct {
		m1   map[int]int
		m2   map[int]int
		want map[int]int
	}{
		{nil, nil, nil},
		{map[int]int{1: 1}, nil, map[int]int{1: 1}},
		{map[int]int{1: 1}, map[int]int{2: 2}, map[int]int{1: 1, 2: 2}},
		{map[int]int{1: 10, 2: 20}, map[int]int{3: 30}, map[int]int{1: 10, 2: 20, 3: 30}},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var out = MapUnion(tt.m1, tt.m2)
			assert.True(t, maps.Equal(tt.want, out))
		})
	}
}

func TestMapIntersect(t *testing.T) {
	tests := []struct {
		m1   map[int]string
		m2   map[int]string
		want map[int]string
	}{
		{nil, nil, nil},
		{map[int]string{1: "a"}, nil, nil},
		{map[int]string{1: "a"}, map[int]string{2: "b"}, nil},
		{map[int]string{1: "a", 2: "b", 3: "c"}, map[int]string{1: "1", 2: "b", 3: "3"}, map[int]string{2: "b"}},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var out = MapIntersect(tt.m1, tt.m2)
			assert.True(t, maps.Equal(tt.want, out))
		})
	}
}
