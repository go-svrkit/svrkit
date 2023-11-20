// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package debug

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	timestampLayout = "2006-01-02T15:04:05-0700" // IOS8601
	mainPkgName     = "main.main"
)

// code taken from https://github.com/pkg/error with modification

// Frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type Frame uintptr

// PC returns the program counter for this frame;
// multiple frames may have the same PC value.
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
		names = append(names, fnName+"()")
		if len(names) >= limit || fnName == "main.main" {
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
		if fnName == mainPkgName {
			break
		}
	}
	return sb.String()
}

// GetCallerStack 获取当前调用堆栈
func GetCallerStack(stack *Stack, skip int) {
	if stack.pcs == nil {
		stack.pcs = make([]uintptr, 32) // 32 depth is enough
	}
	n := runtime.Callers(skip+1, stack.pcs[:])
	stack.pcs = stack.pcs[0:n]
}

func GetCurrentCallStack(skip int) Stack {
	var stack Stack
	GetCallerStack(&stack, skip+1)
	return stack
}

func Backtrace(title string, val interface{}, w io.Writer) {
	var stack Stack
	GetCallerStack(&stack, 1)
	var now = time.Now()
	fmt.Fprintf(w, "%s\nTraceback[%s] (most recent call last):\n", title, now.Format(timestampLayout))
	fmt.Fprintf(w, "%v %v\n", stack, val)
}

func CatchPanic(title string) {
	if v := recover(); v != nil {
		Backtrace(title, v, os.Stderr)
	}
}
