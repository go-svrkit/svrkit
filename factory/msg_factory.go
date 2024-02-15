// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package factory

import (
	"fmt"
	"hash/crc32"
	"log"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// protobuf协议说明
// 1. `Ntf`结尾的消息表示通知，Client <-> Server
// 2. `Req`结尾的消息表示请求，Client -> Server；
// 3. `Ack`结尾的消息表示请求的响应，Server -> Client；
// 4. `Req`和`Ack`名字前缀需要匹配，即`FooReq`和`FooAck`表示一对请求/响应协议；

var (
	name2id = make(map[string]uint32)       // 消息名称 <-> 消息ID
	id2name = make(map[uint32]string)       // 消息ID <-> 消息名称
	id2type = make(map[uint32]reflect.Type) // 消息ID <-> reflect
	type2id = make(map[reflect.Type]uint32) // reflect <-> 消息ID
)

func Clear() {
	name2id = make(map[string]uint32)
	id2name = make(map[uint32]string)
	id2type = make(map[uint32]reflect.Type)
	type2id = make(map[reflect.Type]uint32)
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

func registerByNameHash(fd protoreflect.FileDescriptor) bool {
	if isWellKnown(fd.Path()) {
		return true
	}
	// log.Printf("register %s", fd.Path())
	var descriptors = fd.Messages()
	for i := 0; i < descriptors.Len(); i++ {
		var descriptor = descriptors.Get(i)
		var fullname = string(descriptor.FullName())
		if !HasValidSuffix(fullname) {
			continue
		}
		var name = string(descriptor.Name())
		var rtype = proto.MessageType(fullname)
		if rtype == nil {
			log.Printf("message %s cannot be reflected\n", fullname)
			continue
		}
		if GetMessageId(name) == 0 {
			if err := Register(rtype); err != nil {
				log.Printf("register msg %s: %v\n", name, err)
			}
		}
	}
	return true
}

// RegisterAllMessages 自动注册所有protobuf消息
// protobuf使用init()注册(RegisterType)，则此API需要在import后调用
func RegisterAllMessages() {
	protoregistry.GlobalFiles.RangeFiles(registerByNameHash)
	log.Printf("%d proto message registered", len(id2type))
}

func Register(rType reflect.Type) error {
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}
	// 不能重复注册
	var name = rType.String()
	if _, ok := name2id[name]; ok {
		return fmt.Errorf("duplicate registration of %s", name)
	}
	// 不同的名字如果生成了相同的hash，需要更改新名字
	var hash = NameHash(name)
	if old, found := id2name[hash]; found {
		return fmt.Errorf("duplicate hash of %s -> %s", old, name)
	}
	name2id[name] = hash
	id2type[hash] = rType
	id2name[hash] = name
	type2id[rType] = hash
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
	return id2type[msgId]
}

// GetMessageIdOf 获取proto消息的ID
func GetMessageIdOf(msg proto.Message) uint32 {
	var rtype = reflect.TypeOf(msg)
	return type2id[rtype.Elem()]
}

// CreateMessageByID 根据消息ID创建消息（使用反射）
func CreateMessageByID(msgId uint32) proto.Message {
	if rtype, found := id2type[msgId]; found {
		var val = reflect.New(rtype).Interface()
		return val.(proto.Message)
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
