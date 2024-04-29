package zset

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/collections/cutil"
)

func TestSortedSet_Add(t *testing.T) {
	var st = NewSortedSet(cutil.OrderedCmp[string])
	assert.Equal(t, 0, st.Len())
	st.Add("key", 10)
	st.Add("key", 20)
	st.Add("key", 30)
	assert.Equal(t, 1, st.Len())
	assert.Equal(t, int64(30), st.GetScore("key"))
}

func TestSortedSet_Remove(t *testing.T) {
	var st = NewSortedSet(cutil.OrderedCmp[string])
	assert.False(t, st.Remove("key"))
	assert.Equal(t, 0, st.Len())
	st.Add("key", 10)
	assert.Equal(t, 1, st.Len())
	assert.True(t, st.Remove("key"))
	assert.Equal(t, 0, st.Len())

	for i := 1; i <= 100; i++ {
		var key = fmt.Sprintf("test%d", i)
		st.Add(key, int64(i))
	}
	assert.Equal(t, 100, st.Len())

	for i := 100; i > 0; i-- {
		var key = fmt.Sprintf("test%d", i)
		assert.True(t, st.Remove(key))
	}
	assert.Equal(t, 0, st.Len())
}

func TestSortedSet_Count(t *testing.T) {
	var st = NewSortedSet(cutil.OrderedCmp[string])
	for i := 1; i <= 100; i++ {
		var key = fmt.Sprintf("test%d", i)
		st.Add(key, int64(i))
	}
	assert.Equal(t, 100, st.Len())

	tests := []struct {
		min      int64
		max      int64
		expected int
	}{
		{0, 1000, 100},
		{1, 100, 100},
		{2, 100, 99},
		{100, 100, 1},
		{101, 100, 0},
		{101, 200, 0},
		{0, 0, 0},
		{0, -1, 0},
		{-100, -1, 0},
		{101, 0, 0},
		{101, -1, 0},
	}
	for _, tc := range tests {
		assert.Equal(t, tc.expected, st.Count(tc.min, tc.max))
	}
}

func TestSortedSet_GetRank(t *testing.T) {
	var st = NewSortedSet(cutil.OrderedCmp[string])
	for i := 100; i > 0; i-- {
		var key = fmt.Sprintf("test%d", i)
		st.Add(key, int64(i))
	}
	assert.Equal(t, 100, st.Len())

	tests := []struct {
		input    string
		reverse  bool
		expected int // 排名从0开始
	}{
		{"test100", false, 99},
		{"test100", true, 0},
		{"test1", false, 0},
		{"test1", true, 99},
		{"test50", false, 49},
		{"test1234", false, -1},
		{"test1234", true, -1},
	}
	for _, tc := range tests {
		assert.Equal(t, tc.expected, st.GetRank(tc.input, tc.reverse))
	}
}

func TestSortedSet_GetScore(t *testing.T) {
	var st = NewSortedSet(cutil.OrderedCmp[string])
	for i := 100; i > 0; i-- {
		var key = fmt.Sprintf("test%d", i)
		st.Add(key, int64(i))
	}
	assert.Equal(t, 100, st.Len())

	for i := 1; i <= 100; i++ {
		var key = fmt.Sprintf("test%d", i)
		assert.Equal(t, int64(i), st.GetScore(key))
	}
}

func TestSortedSet_GetRange(t *testing.T) {
	var st = NewSortedSet(cutil.OrderedCmp[string])
	for i := 1; i <= 20; i++ {
		st.Add(strconv.Itoa(i), int64(i))
	}
	assert.Equal(t, 20, st.Len())

	// 排名从0开始
	tests := []struct {
		start    int
		end      int
		reverse  bool
		expected string
	}{
		{1, 10, false, "2,3,4,5,6,7,8,9,10,11"},
		{1, 10, true, "19,18,17,16,15,14,13,12,11,10"},
		{0, 0, false, "1"},
		{0, 0, true, "20"},
		{0, 1, false, "1,2"},
		{0, 1, true, "20,19"},
		{0, 2, false, "1,2,3"},
		{0, 2, true, "20,19,18"},
		{0, 3, false, "1,2,3,4"},
		{0, 3, true, "20,19,18,17"},
		{0, 4, false, "1,2,3,4,5"},
		{0, 4, true, "20,19,18,17,16"},
	}
	for _, tc := range tests {
		var list = st.GetRange(tc.start, tc.end, tc.reverse)
		var output = strings.Join(list, ",")
		assert.Equal(t, tc.expected, output)
	}
}

func TestSortedSet_GetRangeByScore(t *testing.T) {
	var st = NewSortedSet(cutil.OrderedCmp[string])
	for i := 1; i <= 20; i++ {
		st.Add(strconv.Itoa(i), int64(i))
	}
	assert.Equal(t, 20, st.Len())

	// 排名从0开始
	tests := []struct {
		min      int64
		max      int64
		reverse  bool
		expected string
	}{
		{1, 10, false, "1,2,3,4,5,6,7,8,9,10"},
		{1, 10, true, "10,9,8,7,6,5,4,3,2,1"},
		{0, 0, false, ""},
		{0, 0, true, ""},
		{0, 1, false, "1"},
		{18, 20, false, "18,19,20"},
		{18, 20, true, "20,19,18"},
	}
	for _, tc := range tests {
		var list = st.GetRangeByScore(tc.min, tc.max, tc.reverse)
		var output = strings.Join(list, ",")
		assert.Equal(t, tc.expected, output)
	}
}
