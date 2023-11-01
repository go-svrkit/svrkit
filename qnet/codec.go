// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.

package qnet

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"

	"gopkg.in/svrkit.v1/logger"
)

const (
	V1HeaderLength        = 16                             // 消息头大小
	MaxPacketSize         = 0x00FFFFFF                     // 最大消息大小，~16MB
	MaxPayloadSize        = MaxPacketSize - V1HeaderLength //
	MaxClientUpStreamSize = 1 << 18                        // 最大client上行消息大小，256KB
)

var (
	DefaultCompressThreshold = 1 << 12 // 压缩阈值，4KB
)

type MsgFlag uint8

const (
	FlagCompress MsgFlag = 0x01 // 压缩
	FlagEncrypt  MsgFlag = 0x02 // 加密
	FlagError    MsgFlag = 0x04 // 错误
	FlagRequest  MsgFlag = 0x08 //
	FlagFrame    MsgFlag = 0x20 //
	FlagExtent   MsgFlag = 0x40 //
)

func (g MsgFlag) Has(n MsgFlag) bool {
	return g&n != 0
}

func (g MsgFlag) Clear(n MsgFlag) MsgFlag {
	return g &^ n
}

// wire protocol
// ------|---------|---------|--------|---------|---------|----------|
// field |--<len>--|--<crc>--|-<flag>-|--<seq>--|--<cmd>--|--<data>--|
// bytes |----3----|----4----|----1---|----4----|----4----|---N/A----|
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
	return MsgFlag(h[7])
}

func (h NetV1Header) Seq() uint32 {
	return binary.LittleEndian.Uint32(h[8:])
}

func (h NetV1Header) Command() uint32 {
	return binary.LittleEndian.Uint32(h[12:])
}

func (h NetV1Header) CRC() uint32 {
	return binary.LittleEndian.Uint32(h[3:])
}

func (h NetV1Header) SetCRC(v uint32) {
	binary.LittleEndian.PutUint32(h[3:], v)
}

// CalcCRC checksum = f(header[7:]) and f(body)
func (h NetV1Header) CalcCRC(body []byte) uint32 {
	var crc = crc32.NewIEEE()
	crc.Write(h[7:])
	if len(body) > 0 {
		crc.Write(body)
	}
	return crc.Sum32()
}

func (h NetV1Header) Pack(size uint32, flag MsgFlag, netMsg *NetMessage) {
	intToBytes(size, h[:3])
	// h[3:7] = checksum // set after
	h[7] = uint8(flag)
	binary.LittleEndian.PutUint32(h[8:], netMsg.Seq)
	binary.LittleEndian.PutUint32(h[12:], netMsg.Command)
}

// ReadHeadBody @maxSize should less than MaxPacketSize
func ReadHeadBody(rd io.Reader, head NetV1Header, maxSize uint32) ([]byte, error) {
	if _, err := io.ReadFull(rd, head); err != nil {
		return nil, err
	}
	var nLen = head.Len()
	if nLen < V1HeaderLength || nLen > maxSize {
		logger.Errorf("ReadHeadBody: msg size %d out of range", nLen)
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
		logger.Errorf("ReadHeadBody: msg %v checksum mismatch %x != %x", head.Command(), checksum, crc)
		return nil, ErrPktChecksumMismatch
	}
	return body, nil
}

func ProcessHeaderFlags(flags MsgFlag, body []byte, decrypt Encryptor) ([]byte, error) {
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
		if decoded, err := uncompress(body); err != nil {
			return nil, err
		} else {
			body = decoded
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
	var flags = head.Flag()
	body, err = ProcessHeaderFlags(flags, body, decrypt)
	if err != nil {
		return err
	}
	netMsg.Seq = head.Seq()
	netMsg.Command = head.Command()

	if flags.Has(FlagError) {
		n, i := binary.Uvarint(body)
		if i <= 0 {
			return fmt.Errorf("decode msg %d errno negative %d", netMsg.Command, i)
		}
		netMsg.Errno = uint32(n)
	} else {
		netMsg.Data = body
	}
	return nil
}

// EncodeMsgTo encode message to writer
func EncodeMsgTo(netMsg *NetMessage, encrypt Encryptor, w io.Writer) error {
	var flags MsgFlag
	if netMsg.Errno != 0 {
		flags |= FlagError
	}
	if err := netMsg.Encode(); err != nil {
		return err
	}
	var body = netMsg.Data
	if len(body) > DefaultCompressThreshold {
		if encoded, err := compress(body); err == nil {
			if len(encoded) < len(body) {
				flags |= FlagCompress
				body = encoded
			}
		} else {
			logger.Errorf("msg %d compress failed: %v", netMsg.Command, err)
		}
	}
	if encrypt != nil && len(body) > binary.MaxVarintLen32 {
		if encrypted, err := encrypt.Encrypt(body); err == nil {
			body = encrypted
			flags |= FlagEncrypt
		} else {
			return err
		}
	}

	var bodySize = len(body)
	if bodySize > MaxPayloadSize {
		return fmt.Errorf("encoded msg %d size %d/%d overflow", netMsg.Command, bodySize, MaxPayloadSize)
	}

	var head = NewNetV1Header()
	head.Pack(uint32(bodySize+V1HeaderLength), flags, netMsg)
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

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	var w = zlib.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		if er := w.Close(); er != nil {
			return nil, er
		}
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func uncompress(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, r); err != nil {
		if er := r.Close(); er != nil {
			return nil, er
		}
		return nil, err
	}
	if err = r.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
