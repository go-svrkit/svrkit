// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package fat

import (
	"runtime/debug"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func trimLines(content string) []string {
	var list = strings.Split(content, "\n")
	for i := 0; i < len(list); i++ {
		list[i] = strings.TrimSpace(list[i])
	}
	return list
}

func TestTraceStack(t *testing.T) {
	var sb strings.Builder
	TraceStack(1, "test traceback", "err msg", &sb)
	var content = sb.String()
	var lines = trimLines(content)
	assert.Greater(t, len(lines), 4)
	assert.Equal(t, lines[0], "test traceback")
	assert.Contains(t, lines[1], "err msg")
	assert.Contains(t, lines[2], "stack traceback")
	assert.Contains(t, lines[3], "TestTraceStack()")
}

func didPanic() {
	panic("test panic")
}

func TestCatchPanic(t *testing.T) {
	defer CatchPanic("test catch panic")
	didPanic()
}

func TestStartProfiler(t *testing.T) {
	StartProfiler("localhost:16060")
}

func TestReadGCPercent(t *testing.T) {
	var percent = ReadGCPercent()
	assert.Equal(t, uint64(100), percent)
	percent = 200
	debug.SetGCPercent(int(percent))
	assert.Equal(t, percent, ReadGCPercent())
}

func TestReadMemoryLimit(t *testing.T) {
	var limit = ReadMemoryLimit()
	assert.Equal(t, "MaxInt64", limit)
	debug.SetMemoryLimit(64 << 13)
	assert.Equal(t, "64KiB", ReadMemoryLimit())
}

func TestReadMetrics(t *testing.T) {
	var r = ReadMetrics("gc")
	assert.True(t, len(r) > 0)
}
