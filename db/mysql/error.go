// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package mysql

import (
	"database/sql/driver"
	"net"
	"syscall"

	mydriver "github.com/go-sql-driver/mysql"
)

func IsCanRetryErr(err error) bool {
	if err == mydriver.ErrInvalidConn {
		return true
	}
	return false
}

// IsBadConn returns whether err is a connection error
func IsBadConn(err error) bool {
	if err == driver.ErrBadConn {
		return true
	}
	if e, ok := err.(*net.OpError); ok {
		if errno, ok := e.Err.(syscall.Errno); ok && errno == syscall.ECONNREFUSED {
			return true
		}
	}
	return false
}

func IsDeadlock(err error) bool {
	if e, ok := err.(*mydriver.MySQLError); ok {
		return e.Number == 1213
	}
	return false
}

func IsLockWaitTimeout(err error) bool {
	if e, ok := err.(*mydriver.MySQLError); ok {
		return e.Number == 1205
	}
	return false
}
