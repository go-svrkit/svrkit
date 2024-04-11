// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package codec

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Message = proto.Message

type VTProtoMessage interface {
	SizeVT() int
	MarshalVT() ([]byte, error)
	MarshalToVT(data []byte) (int, error)
	MarshalToSizedBufferVT(data []byte) (int, error)
	UnmarshalVT([]byte) error
}

var (
	Merge = proto.Merge
	Clone = proto.Clone
)

func Marshal(m proto.Message) ([]byte, error) {
	if vtM, ok := m.(VTProtoMessage); ok {
		return vtM.MarshalVT()
	}
	return proto.Marshal(m)
}

func Unmarshal(b []byte, m Message) error {
	if vtM, ok := m.(VTProtoMessage); ok {
		return vtM.UnmarshalVT(b)
	}
	return proto.Unmarshal(b, m)
}

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
