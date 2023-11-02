// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package logger

import (
	"testing"
)

func TestLogAPI(t *testing.T) {
	Debugf("test debug log")
	Infof("test info log")
	Errorf("test error log")
}
