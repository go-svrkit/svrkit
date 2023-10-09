// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteFileLog(filename, format string, a ...interface{}) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, format, a...)
	return err
}

func AppFileErrorLog(format string, a ...interface{}) error {
	appname := filepath.Base(os.Args[0])
	if i := strings.LastIndex(appname, "."); i > 0 {
		appname = appname[:i]
	}
	if i := strings.LastIndex(appname, "_"); i > 0 {
		appname = appname[i+1:]
	}
	var filename = fmt.Sprintf("logs/%s_error.log", appname)
	return WriteFileLog(filename, format, a...)
}
