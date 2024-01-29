// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package mysql

import (
	"testing"
)

func TestSimpleOperationSet(t *testing.T) {
	var set = NewDBOperationSet(0)
	if set.Len() != 0 {
		t.Fatalf("set should be empty")
	}
	var op1 = NewDBOperation("INSERT INTO(`?`, `?`, `?`)", 1, 2, 3)
	var op2 = NewDBOperation("INSERT INTO(`?`, `?`, `?`)", 4, 5, 6)
	set.Add(op1)
	set.Add(op2)
	if set.Len() != 2 {
		t.Fatalf("set op count failure")
	}
}
