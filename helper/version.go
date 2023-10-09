// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package helper

import (
	"fmt"
	"runtime"
	"strings"
)

// 一个版本号（如1.0.1）由`major.minor.patch`三部分组成
//
// `major`: 主版本号
// `minor`: 次版本号
// `patch`: 修订号
//
// 第一个初始开发版本使用`0.1.0`
// 第一个可以对外发布的版本使用`1.0.0`
//

var (
	_Version   = "?.?.?"
	_GitBranch = "?"
	_CommitRev = "???"
	_BuildTime = "???"
)

// Version 版本号
func Version() string {
	return _Version
}

// GitBranch 代码分支
func GitBranch() string {
	return _GitBranch
}

// CommitRev 提交版本
func CommitRev() string {
	return _CommitRev
}

// BuildTime 编译时间
func BuildTime() string {
	return _BuildTime
}

func VersionString() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Version: %s\n", _Version)
	fmt.Fprintf(&sb, "Revision: %s-%s\n", _GitBranch, _CommitRev)
	fmt.Fprintf(&sb, "Built at: %s\n", _BuildTime)
	fmt.Fprintf(&sb, "Powered by: %s", runtime.Version())
	return sb.String()
}

func PrintVersion() {
	println(VersionString())
}
