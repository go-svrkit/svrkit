package qnet

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"

	"gopkg.in/svrkit.v1/qlog"
)

const (
	V1HeaderLength = 16                             // 消息头大小
	MaxPacketSize  = 0x00FFFFFF                     // 最大消息大小，~16MB
	MaxPayloadSize = MaxPacketSize - V1HeaderLength //
)

var (
	DefaultCompressThreshold = 1 << 12 // 压缩阈值，4KB
	MaxClientUpStreamSize    = 1 << 18 // 最大client上行消息大小，256KB
)

// a simple wire protocol
// ------|---------|---------|--------|---------|---------|----------|
// field |--<len>--|--<flag>--|-<crc>-|--<seq>--|--<cmd>--|--<data>--|
// bytes |----3----|----1----|----4---|----4----|----4----|---N/A----|
// ------|---------|---------|--------|---------|---------|----------|

// NetV1Header 协议头
type NetV1Header []byte

func NewNetV1Header() NetV1Header {
	return make([]byte, V1HeaderLength)
}

func (h NetV1Header) Len() uint32 {
	return bytesToInt(h[:3])
}

func (h NetV1Header) Flag() MsgFlag {
	return MsgFlag(h[3])
}

func (h NetV1Header) Seq() uint32 {
	return binary.LittleEndian.Uint32(h[8:])
}

func (h NetV1Header) Command() uint32 {
	return binary.LittleEndian.Uint32(h[12:])
}

func (h NetV1Header) CRC() uint32 {
	return binary.LittleEndian.Uint32(h[4:])
}

func (h NetV1Header) SetCRC(v uint32) {
	binary.LittleEndian.PutUint32(h[4:], v)
}

func (h NetV1Header) Clear() {
	for i := 0; i < len(h); i++ {
		h[i] = 0
	}
}

// CalcCRC checksum = f(head) and f(body)
func (h NetV1Header) CalcCRC(body []byte) uint32 {
	var crc = crc32.NewIEEE()
	crc.Write(h[:3])
	crc.Write(h[8:])
	if len(body) > 0 {
		crc.Write(body)
	}
	return crc.Sum32()
}

func (h NetV1Header) Pack(size uint32, flag MsgFlag, seq, cmd uint32) {
	intToBytes(size, h[:3])
	h[3] = uint8(flag)
	// h[4:8] = checksum // set after
	binary.LittleEndian.PutUint32(h[8:], seq)
	binary.LittleEndian.PutUint32(h[12:], cmd)
}

type MsgFlag uint8

const (
	FlagCompress MsgFlag = 0x01
	FlagEncrypt  MsgFlag = 0x02
	FlagExtent   MsgFlag = 0x08
	FlagError    MsgFlag = 0x10
	FlagCache    MsgFlag = 0x20
)

func (g MsgFlag) Has(n MsgFlag) bool {
	return g&n != 0
}

func (g MsgFlag) Clear(n MsgFlag) MsgFlag {
	return g &^ n
}

// ReadHeadBody read header and body less than `maxSize`
func ReadHeadBody(rd io.Reader, head NetV1Header, maxSize uint32) ([]byte, error) {
	if _, err := io.ReadFull(rd, head); err != nil {
		return nil, err
	}
	var nLen = head.Len()
	if nLen < V1HeaderLength || nLen > maxSize {
		qlog.Errorf("ReadHeadBody: msg size %d out of range", nLen)
		return nil, ErrPktSizeOutOfRange
	}
	var body []byte
	if nLen > V1HeaderLength {
		body = make([]byte, nLen-V1HeaderLength)
		if _, err := io.ReadFull(rd, body); err != nil {
			return nil, err
		}
	}
	var checksum = head.CalcCRC(body)
	if crc := head.CRC(); crc != checksum {
		qlog.Errorf("ReadHeadBody: msg %v checksum mismatch %x != %x", head.Command(), checksum, crc)
		return nil, ErrPktChecksumMismatch
	}
	return body, nil
}

func DecodeMsgBody(flags MsgFlag, body []byte, decrypt Encryptor) ([]byte, error) {
	if flags.Has(FlagEncrypt) {
		if decrypt == nil {
			return nil, ErrCannotDecryptPkt
		}
		if decrypted, err := decrypt.Decrypt(body); err != nil {
			return nil, err
		} else {
			body = decrypted
		}
	}
	if flags.Has(FlagCompress) {
		var decoded bytes.Buffer
		if err := uncompress(body, &decoded); err != nil {
			return nil, err
		} else {
			body = decoded.Bytes()
		}
	}
	return body, nil
}

// DecodeMsgFrom decode message from reader
func DecodeMsgFrom(rd io.Reader, maxSize uint32, decrypt Encryptor, netMsg *NetMessage) error {
	var head = NewNetV1Header()
	body, err := ReadHeadBody(rd, head, maxSize)
	if err != nil {
		return err
	}
	return DecodeNetMsg(head, body, decrypt, netMsg)
}

func DecodeNetMsg(head NetV1Header, body []byte, decrypt Encryptor, netMsg *NetMessage) error {
	var err error
	body, err = DecodeMsgBody(head.Flag(), body, decrypt)
	if err != nil {
		return err
	}
	netMsg.Seq = head.Seq()
	netMsg.Command = head.Command()
	netMsg.Data = body
	return nil
}

// EncodeMsgTo encode message to writer
func EncodeMsgTo(netMsg *NetMessage, encrypt Encryptor, w io.Writer) error {
	var flags MsgFlag
	if err := netMsg.Encode(); err != nil {
		return err
	}
	var body = netMsg.Data
	if len(body) > DefaultCompressThreshold {
		var encoded bytes.Buffer
		if err := compress(body, &encoded); err != nil {
			qlog.Errorf("msg %d compress failed: %v", netMsg.Command, err)
			// 不压缩，plain data
		} else {
			if encoded.Len() < len(body) {
				flags |= FlagCompress
				body = encoded.Bytes()
			}
		}
	}
	if encrypt != nil {
		if encrypted, err := encrypt.Encrypt(body); err != nil {
			return err
		} else {
			body = encrypted
			flags |= FlagEncrypt
		}
	}

	var bodySize = len(body)
	if bodySize > MaxPayloadSize {
		return fmt.Errorf("encoded msg %d size %d/%d overflow", netMsg.Command, bodySize, MaxPayloadSize)
	}

	var head = NewNetV1Header()
	head.Pack(uint32(bodySize+V1HeaderLength), flags, netMsg.Seq, netMsg.Command)
	var checksum = head.CalcCRC(body)
	head.SetCRC(checksum)

	if _, err := w.Write(head); err != nil {
		return err
	}
	if bodySize > 0 {
		if _, err := w.Write(body); err != nil {
			return err
		}
	}
	return nil
}

func compress(input []byte, buf *bytes.Buffer) error {
	if len(input) == 0 {
		return nil
	}
	var w = zlib.NewWriter(buf)
	if _, err := w.Write(input); err != nil {
		return err
	}
	return w.Close()
}

func uncompress(input []byte, buf *bytes.Buffer) error {
	if len(input) == 0 {
		return nil
	}
	rd, err := zlib.NewReader(bytes.NewReader(input))
	if err != nil {
		return err
	}
	_, err = io.Copy(buf, rd)
	if err != nil {
		return err
	}
	return rd.Close()
}

// 3-bytes little endian to uint32
func bytesToInt(b []byte) uint32 {
	_ = b[2] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

// uint32 to little endian 3-bytes
func intToBytes(v uint32, b []byte) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
}
