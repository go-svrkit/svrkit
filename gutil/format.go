// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"gopkg.in/svrkit.v1/zlog"
)

// JSONParse 避免大数值被解析为float导致的精度丢失
func JSONParse(s string, v any) error {
	var dec = decoder.NewDecoder(s)
	dec.UseInt64()
	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}

// JSONStringify 序列化为json字符串
func JSONStringify(v any) string {
	data, err := sonic.MarshalString(v)
	if err != nil {
		zlog.Errorf("JSONStringify %T: %v", v, err)
		return ""
	}
	return data
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
