// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"gopkg.in/svrkit.v1/qlog"
)

const (
	KiB = 1 << 10
	MiB = 1 << 20
	GiB = 1 << 30
	TiB = 1 << 40
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

// IsFileExist test if file exist
func IsFileExist(filename string) bool {
	_, err := os.Lstat(filename)
	return !os.IsNotExist(err)
}

func IsFileNotExist(filename string) bool {
	_, err := os.Lstat(filename)
	return os.IsNotExist(err)
}

// ReadFileToLines 把文件内容按一行一行读取
func ReadFileToLines(filename string) ([]string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fd.Close()
	}()

	return ReadToLines(fd)
}

// ReadToLines 把内容按行读取
func ReadToLines(rd io.Reader) ([]string, error) {
	var lines []string
	var scanner = bufio.NewScanner(rd)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// CopyFile writes the contents of the given source file to dest.
func CopyFile(dest, source string) error {
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = io.Copy(df, f)
	return err
}

// JSONParse 避免大数值被解析为float导致的精度丢失
func JSONParse(s string, v any) error {
	var dec = json.NewDecoder(strings.NewReader(s))
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
		qlog.Errorf("JSONStringify %T: %v", v, err)
		return ""
	}
	return unsafe.String(unsafe.SliceData(data), len(data))
}

func UnmarshalProtoJSON(b []byte, m proto.Message) error {
	return jsonpb.Unmarshal(bytes.NewReader(b), m)
}

// MarshalProtoJSON 序列化proto消息为json格式
func MarshalProtoJSON(msg proto.Message) ([]byte, error) {
	var buf bytes.Buffer
	var jm = jsonpb.Marshaler{
		EnumsAsInts:  true,
		EmitDefaults: true,
	}
	if err := jm.Marshal(&buf, msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
