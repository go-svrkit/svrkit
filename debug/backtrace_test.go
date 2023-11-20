// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package debug

import (
	"bytes"
	"slices"
	"strings"
	"testing"
)

var gStack Stack

func testCall4() {
	gStack = GetCurrentCallStack(1)
}

func testCall3() {
	testCall4()
}

func testCall2() {
	testCall3()
}

func testCall1() {
	testCall2()
}

func TestGetCurrentCallStack(t *testing.T) {
	testCall1()
	var content = gStack.String()
	//t.Logf("content: %s", content)
	var idx1 = strings.Index(content, "testCall4")
	if idx1 <= 0 {
		t.Fatalf("not found testCall4")
	}
	var idx2 = strings.Index(content, "testCall3")
	if idx2 <= 0 {
		t.Fatalf("not found testCall3")
	}
	var idx3 = strings.Index(content, "testCall2")
	if idx3 <= 0 {
		t.Fatalf("not found testCall2")
	}
	var idx4 = strings.Index(content, "testCall1")
	if idx4 <= 0 {
		t.Fatalf("not found testCall1")
	}
}

func TestGetStackCallerNames(t *testing.T) {
	testCall1()
	var names = gStack.CallerNames(0)
	if len(names) <= 4 {
		t.Fatalf("names count error")
	}

	if slices.Index(names, "server/base/misc.testCall4()") < 0 {
		t.Fatalf("not found testCall4")
	}
	if slices.Index(names, "server/base/misc.testCall3()") < 0 {
		t.Fatalf("not found testCall3")
	}
	if slices.Index(names, "server/base/misc.testCall2()") < 0 {
		t.Fatalf("not found testCall2")
	}
	if slices.Index(names, "server/base/misc.testCall1()") < 0 {
		t.Fatalf("not found testCall1")
	}

	names = gStack.CallerNames(2)
	if len(names) > 2 {
		t.Fatalf("names count error")
	}
}

func TestBacktrace(t *testing.T) {
	defer CatchPanic("test")
	var buf bytes.Buffer
	Backtrace("test", nil, &buf)
	t.Logf("%s", buf.String())
}
