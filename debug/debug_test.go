package debug

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceStack(t *testing.T) {
	var sb strings.Builder
	TraceStack("test traceback", &sb)
	var content = sb.String()
	var lines = trimLines(content)
	assert.Greater(t, len(lines), 4)
	assert.Equal(t, lines[0], "test traceback")
	assert.Contains(t, lines[1], "stack traceback")
	assert.Contains(t, lines[2], "TraceStack()")
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
