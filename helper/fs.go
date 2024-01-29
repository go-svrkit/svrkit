// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package helper

import (
	"bufio"
	"io"
	"os"
)

// IsFileExist test if file exist
func IsFileExist(filename string) bool {
	_, err := os.Lstat(filename)
	return !os.IsNotExist(err)
}

// ReadFileToLines 把文件内容按一行一行读取
func ReadFileToLines(filename string) ([]string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return ReadToLines(fd)
}

// ReadToLines 把文件内容按行读取
func ReadToLines(rd io.Reader) ([]string, error) {
	var lines []string
	var scanner = bufio.NewScanner(rd)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
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
	defer f.Close()
	_, err = io.Copy(df, f)
	return err
}
