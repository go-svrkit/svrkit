// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package qnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"reflect"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// 协议说明
// 1. `Ntf`结尾的消息表示通知，Client <-> Server
// 2. `Req`结尾的消息表示请求，Client -> Server；
// 3. `Ack`结尾的消息表示请求的响应，Server -> Client；
// 4. `Req`和`Ack`名字前缀需要匹配，即`FooReq`和`FooAck`表示一对请求/响应协议；

var (
	name2id = make(map[string]uint32)                   // <消息名称, 消息ID>
	id2name = make(map[uint32]string)                   // <消息ID, 消息名全称>
	id2type = make(map[uint32]protoreflect.MessageType) // <消息名全称, 消息类型>
)

const MessagePackagePrefix = "protos."

func ClearFactory() {
	clear(name2id)
	clear(id2name)
	clear(id2type)
}

// HasValidSuffix 指定的后缀才自动注册
func HasValidSuffix(name string) bool {
	nameSuffix := []string{"Req", "Ack", "Ntf"}
	for _, suffix := range nameSuffix {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

func IsReqMessage(name string) bool {
	return strings.HasSuffix(name, "Req")
}

func IsAckMessage(name string) bool {
	return strings.HasSuffix(name, "Ack")
}

// GetPairingAckName 根据Req消息名称，返回其对应的Ack消息名称
func GetPairingAckName(fullReqName string) string {
	if len(fullReqName) > 3 && fullReqName[len(fullReqName)-3:] == "Req" {
		return fullReqName[:len(fullReqName)-3] + "Ack"
	}
	return ""
}

func GetPairingAckNameOf(msgId uint32) string {
	return GetPairingAckName(GetMessageFullName(msgId))
}

func isWellKnown(name string) bool {
	return strings.HasPrefix(name, "google/") ||
		strings.HasPrefix(name, "github.com/") ||
		strings.HasPrefix(name, "grpc/")
}

// HashMsgName 计算字符串的hash值
func HashMsgName(name string) uint32 {
	var h = crc32.NewIEEE()
	h.Write([]byte(name))
	return h.Sum32()
}

func registerByNameHash(mt protoreflect.MessageType) bool {
	var md = mt.Descriptor()
	var fd = md.ParentFile()
	if fd != nil && isWellKnown(fd.Path()) {
		return true
	}
	var fullname = string(md.FullName())
	if !HasValidSuffix(fullname) {
		return true
	}
	if GetMessageId(fullname) == 0 {
		if err := Register(fullname); err != nil {
			log.Printf("register msg %s: %v\n", fullname, err)
		}
	}
	return true
}

// RegisterAllMessages 自动注册所有protobuf消息
// protobuf使用init()注册(RegisterType)，则此API需要在import后调用
func RegisterAllMessages() {
	protoregistry.GlobalTypes.RangeMessages(registerByNameHash)
	log.Printf("registered %d proto messages\n", len(id2name))
}

func Register(fullname string) error {
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(fullname))
	if err != nil {
		return err
	}
	// 不同的名字如果生成了相同的hash，需要更改新名字
	var hash = HashMsgName(fullname)
	if old, found := id2name[hash]; found {
		return fmt.Errorf("duplicate hash of %s -> %s", old, fullname)
	}
	name2id[fullname] = hash
	id2name[hash] = fullname
	id2type[hash] = mt
	return nil
}

// GetMessageFullName 根据消息ID获取消息名称
func GetMessageFullName(msgId uint32) string {
	return id2name[msgId]
}

func GetMessageShortName(msgId uint32) string {
	var name = id2name[msgId]
	if name == "" {
		return ""
	}
	var idx = strings.LastIndex(name, ".")
	if idx > 0 {
		return name[idx+1:]
	}
	return name
}

// GetMessageId 根据消息名称获取消息ID
func GetMessageId(fullName string) uint32 {
	return name2id[fullName]
}

func GetMessageType(msgId uint32) reflect.Type {
	if mt, found := id2type[msgId]; found && mt != nil {
		return reflect.TypeOf(mt.Zero().Interface()).Elem()
	}
	return nil
}

// GetMessageIdOf 获取proto消息的ID
func GetMessageIdOf(msg proto.Message) uint32 {
	var fullname = reflect.TypeOf(msg).Elem().String()
	return name2id[fullname]
}

// CreateMessageByID 根据消息ID创建消息（使用反射）
func CreateMessageByID(msgId uint32) proto.Message {
	//if mt, found := id2type[msgId]; found && mt != nil {
	//	return mt.New().Interface()
	//}
	return nil
}

// CreateMessageByFullName 根据消息名称创建消息（使用反射）
func CreateMessageByFullName(fullName string) proto.Message {
	if hash, found := name2id[fullName]; found {
		return CreateMessageByID(hash)
	}
	return nil
}

func CreateMessageByShortName(name string) proto.Message {
	var fullName = MessagePackagePrefix + name
	if hash, found := name2id[fullName]; found {
		return CreateMessageByID(hash)
	}
	return nil
}

func CreatePairingAck(fullReqName string) proto.Message {
	var fullAckName = GetPairingAckName(fullReqName)
	if fullAckName != "" {
		return CreateMessageByFullName(fullAckName)
	}
	return nil
}

func CreateMsgFromJSON(text []byte) (proto.Message, error) {
	var wsMsg WsRecvMsg
	var dec = json.NewDecoder(bytes.NewReader(text))
	dec.UseNumber()
	if err := dec.Decode(&wsMsg); err != nil {
		return nil, err
	}
	var pbMsg = CreateMessageByShortName(wsMsg.Cmd)
	if pbMsg == nil {
		return nil, fmt.Errorf("cannot create msg %s", wsMsg.Cmd)
	}
	if len(wsMsg.Body) > 0 {
		if err := jsonpb.Unmarshal(bytes.NewReader(wsMsg.Body), pbMsg); err != nil {
			return nil, err
		}
	}
	return pbMsg, nil
}

func CreateJSONFromProto(msgId uint32, data []byte) (*WsWriteMsg, error) {
	var pbMsg = CreateMessageByID(msgId)
	if pbMsg == nil {
		return nil, fmt.Errorf("cannot create msg %d", msgId)
	}
	if err := proto.Unmarshal(data, pbMsg); err != nil {
		return nil, err
	}
	var wsMsg = &WsWriteMsg{
		Cmd:  GetMessageShortName(msgId),
		Body: pbMsg,
	}
	return wsMsg, nil
}
