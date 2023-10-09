// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package debug

import (
	"bytes"
	"testing"
)

func TestBacktrace(t *testing.T) {
	defer CatchPanic("test")
	var buf bytes.Buffer
	Backtrace("test", nil, &buf)
	t.Logf("%s", buf.String())
}
