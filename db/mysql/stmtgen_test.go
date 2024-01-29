// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package mysql

import (
	"testing"
	"time"
)

// 用户账号信息
type UserAccountInfo struct {
	Uid               int64     `db:"name=uid type=bigint index=primary,auto_incr"`
	Account           string    `db:"type=varchar(40) index=unique,order:desc"`
	Salt              []byte    `db:"type=varchar(20) index=name:idx_hash,type:hash"`
	CertificateString string    `db:"type=varchar(40)"`
	AccountType       uint16    `db:"index=order:desc"`
	Status            uint16    ``
	RegChannel        uint32    `db:"index=name:idx_recent,priority:4"`
	Email             string    `db:"type=varchar(255) index=name:idx_recent,priority:3"`
	RegIP             string    `db:"type=varchar(50) index=name:idx_recent,priority:2"`
	RegTime           time.Time `db:"index=name:idx_recent,priority:1"`
	RegDeviceType     string    `db:"type=varchar(50)"`
	RegDeviceOs       string    `db:"type=varchar(100)"`
	RegAppVer         string    `db:"type=varchar(20)"`
	RegDeviceId       string    `db:"type=varchar(100)"`
}

func init() {
	var ptr *UserAccountInfo
	_gen.RegisterStruct(DefaultDBTagKey, ptr, "")
}

func TestGenSelectStmt(t *testing.T) {
	var ptr *UserAccountInfo
	var stmt = _gen.SelectStmtOf(ptr)
	t.Logf(stmt)
}

func TestGenInsertStmt(t *testing.T) {
	var ptr *UserAccountInfo
	var stmt = _gen.InsertStmtOf(ptr)
	t.Logf(stmt)
}

func TestGenUpdateStmt(t *testing.T) {
	var info = &UserAccountInfo{
		Uid:               123456789,
		Account:           "guest001",
		Salt:              []byte("xfjskxz"),
		CertificateString: "",
		RegChannel:        1001,
		Email:             "admin@example.com",
		RegIP:             "192.168.0.123",
		RegTime:           time.Now(),
		RegDeviceType:     "PC",
		RegDeviceOs:       "Windows 7",
		RegAppVer:         "1.2.3",
		RegDeviceId:       "x0x0x0x0x0x0",
	}
	var keys = []string{"Uid", "Account"}
	var sql = _gen.UpdateStmtOf(info)
	t.Logf(sql)

	stmt, args := _gen.UpdateQueryOf(info, keys)
	t.Logf(stmt)
	t.Logf("%#v", args)
}
