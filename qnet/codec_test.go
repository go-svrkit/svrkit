// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.

package qnet

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"gopkg.in/svrkit.v1/helper"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 随机长度的字符串
func randString(length int) string {
	if length <= 0 {
		return ""
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		idx := rand.Int() % len(alphabet)
		result[i] = alphabet[idx]
	}
	return helper.BytesAsStr(result)
}

// 校验账号信息
type testProtoMsgReq struct {
	Account     string `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	ServerId    int32  `protobuf:"varint,2,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty"`
	Session     uint32 `protobuf:"varint,3,opt,name=session,proto3" json:"session,omitempty"`
	AccessToken string `protobuf:"bytes,4,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	Lang        string `protobuf:"bytes,5,opt,name=lang,proto3" json:"lang,omitempty"`
	AppChannel  string `protobuf:"bytes,6,opt,name=app_channel,json=appChannel,proto3" json:"app_channel,omitempty"`
	AppDevice   string `protobuf:"bytes,7,opt,name=app_device,json=appDevice,proto3" json:"app_device,omitempty"`
	DeviceId    string `protobuf:"bytes,8,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
	AppOs       string `protobuf:"bytes,9,opt,name=app_os,json=appOs,proto3" json:"app_os,omitempty"`
	AppVersion  string `protobuf:"bytes,10,opt,name=app_version,json=appVersion,proto3" json:"app_version,omitempty"`
	ResVersion  string `protobuf:"bytes,11,opt,name=res_version,json=resVersion,proto3" json:"res_version,omitempty"`
	SdkVersion  string `protobuf:"bytes,12,opt,name=sdk_version,json=sdkVersion,proto3" json:"sdk_version,omitempty"`
	Timestamp   int64  `protobuf:"varint,13,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	IsGuest     bool   `protobuf:"varint,14,opt,name=is_guest,json=isGuest,proto3" json:"is_guest,omitempty"`
	ClientIp    string `protobuf:"bytes,15,opt,name=client_ip,json=clientIp,proto3" json:"client_ip,omitempty"`
}

func (m *testProtoMsgReq) Reset()         { *m = testProtoMsgReq{} }
func (m *testProtoMsgReq) String() string { return proto.CompactTextString(m) }
func (*testProtoMsgReq) ProtoMessage()    {}

func isMsgEqual(t *testing.T, a, b *NetMessage) bool {
	if !(a.MsgID == b.MsgID && a.Seq == b.Seq) {
		return false
	}
	if err := a.Encode(); err != nil {
		t.Fatalf("encode %v", err)
	}
	if err := b.Encode(); err != nil {
		t.Fatalf("encode %v", err)
	}
	data1, data2 := a.Data, b.Data
	if len(data1) > 0 && len(data2) > 0 {
		return bytes.Equal(data1, data2)
	}
	return a.Body == nil && b.Body == nil
}

func testEncode(t *testing.T, size int) {
	var netMsg = AllocNetMessage()
	netMsg.MsgID = 1001
	netMsg.Seq = 1

	if size > 0 {
		var body = &testProtoMsgReq{
			Account:    "test001",
			ServerId:   2,
			Session:    123456789,
			Lang:       "en/Us",
			Timestamp:  time.Now().Unix(),
			AppChannel: "wx-ios",
			AppDevice:  "huawei mate40",
			AppOs:      "Android 10",
			AppVersion: "1.0.1",
		}
		body.AccessToken = randString(size)
		netMsg.Body = body
	}

	var buf bytes.Buffer
	if err := EncodeMsgTo(netMsg, nil, &buf); err != nil {
		t.Fatalf("%v", err)
	}
	var msg2 = AllocNetMessage()
	if err := DecodeMsgFrom(&buf, MaxPayloadSize, nil, msg2); err != nil {
		t.Fatalf("%v", err)
	}
	if !isMsgEqual(t, netMsg, msg2) {
		t.Fatalf("size %d not equal", size)
	}
	if len(msg2.Data) > 0 {
		var req testProtoMsgReq
		if err := msg2.DecodeTo(&req); err != nil {
			t.Fatalf("decode %v", err)
		}
		//t.Logf("%s", req.String())
	}
}

func TestCodecEncode(t *testing.T) {
	testEncode(t, 0)
	testEncode(t, 64)
	testEncode(t, DefaultCompressThreshold)
	testEncode(t, MaxPayloadSize-HeaderLength-850)
	//testEncode(t, MaxPayloadSize)
}
