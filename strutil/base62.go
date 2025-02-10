// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package strutil

import (
	"slices"
	"unsafe"
)

var (
	b62Alphabet   = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b62IndexTable = buildIndexTable(b62Alphabet)
)

// build alphabet index
func buildIndexTable(s []byte) map[byte]int64 {
	var table = make(map[byte]int64, len(s))
	for i := 0; i < len(s); i++ {
		table[s[i]] = int64(i)
	}
	return table
}

// EncodeBase62String 编码Base62
func EncodeBase62String(id int64) string {
	if id == 0 {
		return string(b62Alphabet[:1])
	}
	var buf = make([]byte, 0, 12)
	for id > 0 {
		var rem = id % 62
		id /= 62
		buf = append(buf, b62Alphabet[rem])
	}
	slices.Reverse(buf)
	return unsafe.String(unsafe.SliceData(buf), len(buf))
}

// DecodeBase62String 解码Base62
func DecodeBase62String(s string) int64 {
	var n int64
	for i := 0; i < len(s); i++ {
		n = (n * 62) + b62IndexTable[s[i]]
	}
	return n
}
