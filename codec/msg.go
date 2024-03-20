// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package codec

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Message = proto.Message

var (
	Marshal   = proto.Marshal
	Unmarshal = proto.Unmarshal
)

func UnmarshalProtoJSON(b []byte, m proto.Message) error {
	var opt = protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
	return opt.Unmarshal(b, m)
}

// MarshalProtoJSON 序列化proto消息为json格式
func MarshalProtoJSON(msg proto.Message) ([]byte, error) {
	var opt = protojson.MarshalOptions{
		AllowPartial:   true,
		UseEnumNumbers: true,
	}
	return opt.Marshal(msg)
}
