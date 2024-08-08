// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gext

import (
	"fmt"
	"runtime"
	"strings"
)

// code taken from https://github.com/pkg/errors with modification

// Frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type Frame uintptr

// PC returns the program counter for this frame; multiple frames may have the same PC value.
func (f Frame) PC() uintptr { return uintptr(f) - 1 }

// Stack represents a stack of program counters.
type Stack struct {
	pcs []uintptr
}

// CallerNames 获取堆栈函数名
func (s Stack) CallerNames(limit int) []string {
	if limit <= 0 || limit > len(s.pcs) {
		limit = len(s.pcs)
	}
	var names = make([]string, 0, limit)
	for _, v := range s.pcs {
		pc := Frame(v).PC()
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}
		fnName := fn.Name()
		if i := strings.LastIndex(fnName, "/"); i > 0 {
			fnName = fnName[i+1:]
		}
		names = append(names, fnName+"()")
		if len(names) >= limit || fnName == "runtime·goexit" {
			break
		}
	}
	return names
}

func (s Stack) String() string {
	var sb strings.Builder
	for i, v := range s.pcs {
		pc := Frame(v).PC()
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}
		file, line := fn.FileLine(pc)
		fnName := fn.Name()
		fmt.Fprintf(&sb, "% 3d. %s() %s:%d\n", i+1, fnName, file, line)
		if fnName == "runtime·goexit" {
			break
		}
	}
	return sb.String()
}

// GetCallerStack 获取当前调用栈
func GetCallerStack(skip int) Stack {
	var stack = Stack{
		pcs: make([]uintptr, 32), // 32 depth is enough
	}
	n := runtime.Callers(skip+1, stack.pcs)
	stack.pcs = stack.pcs[0:n]
	return stack
}

// GetCallStackNames 当前调用栈得名称
func GetCallStackNames(skip, limit int) []string {
	var stack = GetCallerStack(skip + 1)
	return stack.CallerNames(limit)
}
