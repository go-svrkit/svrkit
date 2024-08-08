// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package strutil

import (
	"math/rand"
	"unicode"
	"unsafe"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-=~!@#$%^&*()_+[]\\;',./{}|:<>?"

// RandString 随机长度的字符串
func RandString(length int) string {
	if length <= 0 {
		return ""
	}
	var buf = make([]byte, length)
	for i := 0; i < length; i++ {
		idx := rand.Int() % len(alphabet)
		buf[i] = alphabet[idx]
	}
	return unsafe.String(unsafe.SliceData(buf), len(buf))
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
