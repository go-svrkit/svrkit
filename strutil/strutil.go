// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package strutil

import (
	"math/rand"
	"strings"
	"unicode"
	"unsafe"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-=~!@#$%^&*()_+[]\\;',./{}|:<>?"

// StrAsBytes returns the bytes backing a string, it is the caller's responsibility not to mutate the bytes returned.
// see https://pkg.go.dev/unsafe#Pointer rule(6)
func StrAsBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesAsStr returns the string view of byte slice
func BytesAsStr(b []byte) string {
	//return *(*string)(unsafe.Pointer(&b))
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// RandString 随机长度的字符串
func RandString(length int) string {
	if length <= 0 {
		return ""
	}
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		idx := rand.Int() % len(alphabet)
		sb.WriteByte(alphabet[idx])
	}
	return sb.String()
}

// RandBytes 随机长度的字节数组
func RandBytes(length int) []byte {
	if length <= 0 {
		return nil
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		ch := uint8(rand.Int31() % 0xFF)
		result[i] = ch
	}
	return result
}

// FindFirstDigit 查找第一个数字的位置
func FindFirstDigit(s string) int {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return i
		}
	}
	return -1
}

// FindFirstNonDigit 查找第一个非数字的位置
func FindFirstNonDigit(s string) int {
	for i, r := range s {
		if !unicode.IsDigit(r) {
			return i
		}
	}
	return -1
}

// Reverse 反转字符串
func Reverse(str string) string {
	runes := []rune(str)
	l := len(runes)
	for i := 0; i < l/2; i++ {
		runes[i], runes[l-i-1] = runes[l-i-1], runes[i]
	}
	return string(runes)
}

// LongestCommonPrefix 字符串最长共同前缀
func LongestCommonPrefix(s1, s2 string) string {
	if s1 == "" || s2 == "" {
		return ""
	}
	i := 0
	for i < len(s1) && i < len(s2) {
		if s1[i] == s2[i] {
			i++
			continue
		}
		break
	}
	return s1[:i]
}
