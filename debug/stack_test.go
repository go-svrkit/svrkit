// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testCall4() Stack {
	return GetCallerStack(1)
}

func testCall3() Stack {
	return testCall4()
}

func testCall2() Stack {
	return testCall3()
}

func testCall1() Stack {
	return testCall2()
}

func TestGetCurrentCallStack(t *testing.T) {
	var stack = testCall1()
	var lines = trimLines(stack.String())
	//t.Logf("stack trace:\n%s", lines)
	assert.Greater(t, len(lines), 4)
	assert.Contains(t, lines[0], "debug.testCall4")
	assert.Contains(t, lines[1], "debug.testCall3")
	assert.Contains(t, lines[2], "debug.testCall2")
	assert.Contains(t, lines[3], "debug.testCall1")
}

func TestGetStackCallerNames(t *testing.T) {
	var stack = testCall1()
	var names = stack.CallerNames(4)
	assert.LessOrEqual(t, len(names), 4)

	assert.Equal(t, "debug.testCall4()", names[0])
	assert.Equal(t, "debug.testCall3()", names[1])
	assert.Equal(t, "debug.testCall2()", names[2])
	assert.Equal(t, "debug.testCall1()", names[3])
}

func TestGetCallStackNames(t *testing.T) {
	var names = GetCallStackNames(1, 0)
	assert.Greater(t, len(names), 0)
	assert.Equal(t, names[0], "debug.TestGetCallStackNames()")
}
