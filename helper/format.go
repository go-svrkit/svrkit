package helper

import (
	"bytes"
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

func FormatProtoMsg(msg proto.Message) string {
	var jm = jsonpb.Marshaler{EnumsAsInts: true}
	var sb strings.Builder
	if err := jm.Marshal(&sb, msg); err != nil {
		logger.Errorf("marshal %T: %v", msg, err)
	} else {
		return sb.String()
	}
	return msg.String()
}
