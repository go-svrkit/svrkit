// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// IsFileExist test if file exist
func IsFileExist(filename string) bool {
	_, err := os.Lstat(filename)
	return !os.IsNotExist(err)
}

func IsFileNotExist(filename string) bool {
	_, err := os.Lstat(filename)
	return os.IsNotExist(err)
}

// ReadFileToLines 把文件内容按一行一行读取
func ReadFileToLines(filename string) ([]string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fd.Close()
	}()

	return ReadToLines(fd)
}

// ReadToLines 把内容按行读取
func ReadToLines(rd io.Reader) ([]string, error) {
	var lines []string
	var scanner = bufio.NewScanner(rd)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// CopyFile writes the contents of the given source file to dest.
func CopyFile(dest, source string) error {
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = io.Copy(df, f)
	return err
}

var EssentialProjDirs = []string{"bin", "config", "data", "logs"}

// IsProjRootDir 如果有以下几个目录，就认为是项目的根路径
func IsProjRootDir(path string) bool {
	for _, dir := range EssentialProjDirs {
		var fullpath = filepath.Join(path, dir)
		if IsFileNotExist(fullpath) {
			return false
		}
	}
	return true
}

const defaultProjRootPath = "./" // 默认当前目录

func GetProjRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return defaultProjRootPath
	}
	path, err := filepath.Abs(cwd)
	if err != nil {
		return defaultProjRootPath
	}
	// 一直往上找直到根目录
	for depth := 16; depth > 0 && !IsProjRootDir(path); depth-- {
		var idx = strings.LastIndexByte(path, filepath.Separator)
		if idx >= 0 {
			path = path[:idx]
		} else {
			return defaultProjRootPath // not found
		}
	}
	return path
}

func GetProjRootDirOf(path string) string {
	var rootDir = GetProjRootDir()
	return filepath.Join(rootDir, path)
}
