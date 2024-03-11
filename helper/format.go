// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package helper

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"gopkg.in/svrkit.v1/strutil"
	"gopkg.in/svrkit.v1/zlog"
)

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

// PrettyBytes 打印容量大小
func PrettyBytes(nbytes int) string {
	var sign = ""
	if nbytes < 0 {
		sign = "-"
		nbytes = -nbytes
	}
	if nbytes < KB {
		return fmt.Sprintf("%s%dB", sign, nbytes)
	} else if nbytes < MB {
		return fmt.Sprintf("%s%.2fKB", sign, float64(nbytes)/KB)
	} else if nbytes < GB {
		return fmt.Sprintf("%s%.2fMB", sign, float64(nbytes)/MB)
	}
	return fmt.Sprintf("%s%.2fGB", sign, float64(nbytes)/GB)
}

// JSONParse 避免大数值被解析为float导致的精度丢失
func JSONParse(data []byte, v any) error {
	var dec = json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}

// JSONStringify 序列化为json字符串
func JSONStringify(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		zlog.Errorf("JSONStringify %T: %v", v, err)
		return ""
	}
	return strutil.BytesAsStr(data)
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
