// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package codec

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// TODO: switch to different protobuf implementations
// 	1. google.golang.org/protobuf
// 	2. github.com/gogo/protobuf
// 	3. github.com/planetscale/vtprotobuf

type Message = proto.Message

var (
	Marshal            = proto.Marshal
	Unmarshal          = proto.Unmarshal
	UnmarshalProtoJSON = jsonpb.Unmarshal
)

// MarshalProtoJSON 序列化proto消息为json格式
func MarshalProtoJSON(msg proto.Message) ([]byte, error) {
	var jm = jsonpb.Marshaler{EnumsAsInts: true}
	var sb bytes.Buffer
	if err := jm.Marshal(&sb, msg); err != nil {
		return nil, err
	} else {
		return sb.Bytes(), nil
	}
}
