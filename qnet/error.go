// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"errors"
	"fmt"
)

var (
	ErrPktSizeOutOfRange    = errors.New("packet size out of range")
	ErrPktChecksumMismatch  = errors.New("packet checksum mismatch")
	ErrCannotDecryptPkt     = errors.New("cannot decrypt packet")
	ErrConnNotRunning       = errors.New("connection not running")
	ErrConnOutboundOverflow = errors.New("connection outbound queue overflow")
	ErrConnForceClose       = errors.New("connection forced to close")
	ErrBufferOutOfRange     = errors.New("buffer out of range")
)

type Error struct {
	Err      error
	Endpoint Endpoint
}

func NewError(err error, endpoint Endpoint) *Error {
	return &Error{
		Err:      err,
		Endpoint: endpoint,
	}
}

func (e Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("node %v(%s) EOF", e.Endpoint.GetNode(), e.Endpoint.GetRemoteAddr())
	}
	return fmt.Sprintf("node %v(%s) %s", e.Endpoint.GetNode(), e.Endpoint.GetRemoteAddr(), e.Err.Error())
}
