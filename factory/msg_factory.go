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

func IsRequestMessage(name string) bool {
	return strings.HasSuffix(name, "Req")
}

func IsAckMessage(name string) bool {
	return strings.HasSuffix(name, "Ack")
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
		if _, found := name2id[name]; !found {
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
	if kind := rType.Kind(); kind != reflect.Ptr {
		return fmt.Errorf("unexpected type %s", rType.String())
	}
	var name = rType.Elem().Name()
	var hash = NameHash(name)
	// 不能重复注册
	if _, ok := name2id[name]; ok {
		return fmt.Errorf("duplicate registration of %s", name)
	}
	// 不同的名字如果生成了相同的hash，需要更改新名字
	if old, found := id2name[hash]; found {
		return fmt.Errorf("duplicate hash of %s -> %s", old, name)
	}
	name2id[name] = hash
	id2type[hash] = rType
	id2name[hash] = name
	type2id[rType] = hash
	return nil
}

// GetMessageName 根据消息ID获取消息名称
func GetMessageName(msgId uint32) string {
	return id2name[msgId]
}

// GetMessageID 根据消息名称获取消息ID
func GetMessageID(name string) uint32 {
	return name2id[name]
}

// CreateMessageByID 根据消息ID创建消息（使用反射）
func CreateMessageByID(msgId uint32) proto.Message {
	if rtype, found := id2type[msgId]; found {
		var val = reflect.New(rtype.Elem()).Interface()
		return val.(proto.Message)
	}
	return nil
}

// CreateMessageByName 根据消息名称创建消息（使用反射）
func CreateMessageByName(name string) proto.Message {
	if hash, found := name2id[name]; found {
		return CreateMessageByID(hash)
	}
	return nil
}

// GetMessageIDOf 获取proto消息的ID
func GetMessageIDOf(msg proto.Message) uint32 {
	var rtype = reflect.TypeOf(msg)
	return type2id[rtype]
}

// GetPairingAckID 根据Req消息的ID，返回其对应的Ack消息ID
func GetPairingAckID(reqName string) uint32 {
	if !strings.HasSuffix(reqName, "Req") {
		return 0
	}
	var ackName = reqName[:len(reqName)-3] + "Ack"
	return name2id[ackName]
}
