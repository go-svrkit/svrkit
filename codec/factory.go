// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package codec

import (
	"fmt"
	"hash/crc32"
	"log"
	"reflect"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// protobuf协议说明
// 1. `Ntf`结尾的消息表示通知，Client <-> Server
// 2. `Req`结尾的消息表示请求，Client -> Server；
// 3. `Ack`结尾的消息表示请求的响应，Server -> Client；
// 4. `Req`和`Ack`名字前缀需要匹配，即`FooReq`和`FooAck`表示一对请求/响应协议；

var (
	name2id = make(map[string]uint32)                   // 消息名称 <-> 消息ID
	id2name = make(map[uint32]string)                   // 消息ID <-> 消息名称
	id2type = make(map[uint32]protoreflect.MessageType) // 消息名称 <-> 消息类型
)

func Clear() {
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
func GetPairingAckName(reqName string) string {
	if len(reqName) > 3 && reqName[len(reqName)-3:] == "Req" {
		return reqName[:len(reqName)-3] + "Ack"
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

// NameHash 计算字符串的hash值
func NameHash(name string) uint32 {
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
	log.Printf("%d proto message registered", len(id2name))
}

func Register(fullname string) error {
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(fullname))
	if err != nil {
		return err
	}
	// 不能重复注册
	if _, ok := name2id[fullname]; ok {
		return fmt.Errorf("duplicate registration of %s", fullname)
	}
	// 不同的名字如果生成了相同的hash，需要更改新名字
	var hash = NameHash(fullname)
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

// GetMessageId 根据消息名称获取消息ID
func GetMessageId(fullName string) uint32 {
	return name2id[fullName]
}

func GetMessageType(msgId uint32) reflect.Type {
	if mt, found := id2type[msgId]; found && mt != nil {
		return reflect.TypeOf(mt.Zero().Interface())
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
	if mt, found := id2type[msgId]; found && mt != nil {
		return mt.New().Interface()
	}
	return nil
}

// CreateMessageByName 根据消息名称创建消息（使用反射）
func CreateMessageByName(fullName string) proto.Message {
	if hash, found := name2id[fullName]; found {
		return CreateMessageByID(hash)
	}
	return nil
}

func CreatePairingAck(reqName string) proto.Message {
	var ackName = GetPairingAckName(reqName)
	if ackName != "" {
		return CreateMessageByName(ackName)
	}
	return nil
}
