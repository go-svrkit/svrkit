// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

// ZeroOf returns the zero value of the type T
func ZeroOf[T any]() T {
	var zero T
	return zero
}

// MD5Sum 计算MD5值
func MD5Sum(data []byte) string {
	var hash = md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA1Sum 计算SHA1值
func SHA1Sum(data []byte) string {
	var hash = sha1.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA256Sum 计算SHA256值
func SHA256Sum(data []byte) string {
	var hash = sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
