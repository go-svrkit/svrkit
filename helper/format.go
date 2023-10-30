// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package helper

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"gopkg.in/svrkit.v1/logger"
)

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

// PrettyBytes 打印容量大小
func PrettyBytes(nbytes int) string {
	if nbytes < KB {
		return fmt.Sprintf("%dB", nbytes)
	} else if nbytes < MB {
		return fmt.Sprintf("%.1fKB", float64(nbytes)/KB)
	} else if nbytes < GB {
		return fmt.Sprintf("%.2fMB", float64(nbytes)/MB)
	}
	return fmt.Sprintf("%.2fGB", float64(nbytes)/GB)
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

func JSONStringify(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		logger.Errorf("JSONStringify %T: %v", v, err)
		return ""
	}
	return BytesAsStr(data)
}

// Proto2JSON 序列化proto消息为json格式
func Proto2JSON(msg proto.Message) string {
	var jm = jsonpb.Marshaler{EnumsAsInts: true}
	var sb strings.Builder
	if err := jm.Marshal(&sb, msg); err != nil {
		logger.Errorf("marshal %T: %v", msg, err)
	} else {
		return sb.String()
	}
	return msg.String()
}

// JSON2Proto 反序列化json字符串为proto消息
func JSON2Proto(body string, dst proto.Message) error {
	var rd = strings.NewReader(body)
	return jsonpb.Unmarshal(rd, dst)
}

func MD5Sum(data []byte) string {
	var hash = md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

func SHA1Sum(data []byte) string {
	var hash = sha1.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
