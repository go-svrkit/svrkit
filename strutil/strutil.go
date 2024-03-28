// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package strutil

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-=~!@#$%^&*()_+[]\\;',./{}|:<>?"

const (
	KiB = 1 << 10
	MiB = 1 << 20
	GiB = 1 << 30
	TiB = 1 << 40
)

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

func abs64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// PrettyBytes 打印容量大小
func PrettyBytes(nbytes int64) string {
	if nbytes == 0 {
		return "0B"
	}
	var unit = 1
	var prec = 1
	var suffix string
	var absBytes = abs64(nbytes)
	switch {
	case absBytes < KiB:
		suffix = "B"
	case absBytes < MiB:
		unit = KiB
		suffix = "KiB"
	case absBytes < GiB:
		prec = 2
		unit = MiB
		suffix = "MiB"
	case absBytes < TiB:
		prec = 3
		unit = GiB
		suffix = "GiB"
	default:
		prec = 4
		unit = TiB
		suffix = "TiB"
	}
	var s = strconv.FormatFloat(float64(nbytes)/float64(unit), 'f', prec, 64)
	s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	return s + suffix
}

// ParseByteCount parses a string that represents a count of bytes.
// suffixes include B, KiB, MiB, GiB, and TiB represent quantities of bytes as defined by the IEC 80000-13 standard
func ParseByteCount(s string) (int64, bool) {
	// The empty string is not valid.
	if s == "" {
		return 0, false
	}
	// Handle the easy non-suffix case.
	last := s[len(s)-1]
	if last >= '0' && last <= '9' {
		n, er := strconv.ParseInt(s, 10, 64)
		if er != nil || n < 0 {
			return 0, false
		}
		return n, true
	}
	// Failing a trailing digit, this must always end in 'B'.
	// Also at this point there must be at least one digit before
	// that B.
	if last != 'B' || len(s) < 2 {
		return 0, false
	}
	// The one before that must always be a digit or 'i'.
	if c := s[len(s)-2]; c >= '0' && c <= '9' {
		// Trivial 'B' suffix.
		n, er := strconv.ParseInt(s[:len(s)-1], 10, 64)
		if er != nil || n < 0 {
			return 0, false
		}
		return n, true
	} else if c != 'i' {
		return 0, false
	}
	// Finally, we need at least 4 characters now, for the unit
	// prefix and at least one digit.
	if len(s) < 4 {
		return 0, false
	}
	power := 0
	switch s[len(s)-3] {
	case 'K':
		power = 1
	case 'M':
		power = 2
	case 'G':
		power = 3
	case 'T':
		power = 4
	default:
		// Invalid suffix.
		return 0, false
	}
	m := uint64(1)
	for i := 0; i < power; i++ {
		m *= 1024
	}
	n, er := strconv.ParseInt(s[:len(s)-3], 10, 64)
	if er != nil || n < 0 {
		return 0, false
	}
	un := uint64(n)
	if un > math.MaxUint64/m {
		// Overflow.
		return 0, false
	}
	un *= m
	if un > uint64(math.MaxUint64) {
		// Overflow.
		return 0, false
	}
	return int64(un), true
}
