// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gext

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
	assert.Contains(t, lines[0], "gext.testCall4")
	assert.Contains(t, lines[1], "gext.testCall3")
	assert.Contains(t, lines[2], "gext.testCall2")
	assert.Contains(t, lines[3], "gext.testCall1")
}

func TestGetStackCallerNames(t *testing.T) {
	var stack = testCall1()
	var names = stack.CallerNames(4)
	assert.LessOrEqual(t, len(names), 4)

	assert.Equal(t, "gext.testCall4()", names[0])
	assert.Equal(t, "gext.testCall3()", names[1])
	assert.Equal(t, "gext.testCall2()", names[2])
	assert.Equal(t, "gext.testCall1()", names[3])
}

func TestGetCallStackNames(t *testing.T) {
	var names = GetCallStackNames(1, 0)
	assert.Greater(t, len(names), 0)
	assert.Equal(t, names[0], "gext.TestGetCallStackNames()")
}
