// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package fat

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

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
