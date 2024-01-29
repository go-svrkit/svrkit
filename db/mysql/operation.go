// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package mysql

import (
	"bytes"
	"fmt"
)

// DBOperation
type DBOperation struct {
	command string
	args    []interface{}
}

func NewDBOperation(command string, args ...interface{}) *DBOperation {
	return &DBOperation{
		command: command,
		args:    args,
	}
}

// Get operation statement
func (o *DBOperation) Command() string {
	return o.command
}

// Get operation arguments
func (o *DBOperation) Args() []interface{} {
	return o.args
}

func (o DBOperation) String() string {
	if len(o.args) == 0 {
		return o.command
	}

	var buf bytes.Buffer
	buf.WriteString(o.command)
	buf.WriteByte(' ')
	fmt.Fprintf(&buf, "%v", o.args)
	return buf.String()
}

// DBOperationSet
type DBOperationSet struct {
	Ops  []*DBOperation
	Done func()
}

func NewDBOperationSet(capacity int) *DBOperationSet {
	if capacity <= 0 {
		capacity = 16
	}
	return &DBOperationSet{
		Ops: make([]*DBOperation, 0, capacity),
	}
}

func (s *DBOperationSet) Len() int {
	return len(s.Ops)
}

func (s *DBOperationSet) List() []*DBOperation {
	return s.Ops
}

func (s *DBOperationSet) Add(op *DBOperation) {
	s.Ops = append(s.Ops, op)
}

func (s *DBOperationSet) Reset() {
	s.Ops = s.Ops[0:]
}

func (s *DBOperationSet) Set(ops []*DBOperation) {
	s.Ops = ops
}
