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
	HeaderLength             = 16                           // 消息头大小
	MaxPacketSize            = 0x00FFFFFF                   // 最大消息大小，~16MB
	MaxPayloadSize           = MaxPacketSize - HeaderLength // 最大消息体大小
	MaxClientUpStreamSize    = 1 << 18                      // 最大client上行消息大小，256KB
	DefaultCompressThreshold = 1 << 12                      // 压缩阈值，4KB
)

// wire protocol format
// field |--<len>--|--<crc>--|-<flag>-|--<seq>--|--<cmd>--|--<data>--|
// bytes |----3----|----4----|----1---|----4----|----4----|---N/A----|

// NetHeader 协议头
type NetHeader []byte

func NewNetHeader() NetHeader {
	return make([]byte, HeaderLength)
}

func (h NetHeader) Len() uint32 {
	return bytesToInt(h[:3])
}

func (h NetHeader) CRC() uint32 {
	return binary.LittleEndian.Uint32(h[3:])
}

func (h NetHeader) SetCRC(v uint32) {
	binary.LittleEndian.PutUint32(h[3:], v)
}

func (h NetHeader) Flag() MsgFlag {
	return MsgFlag(h[7])
}

func (h NetHeader) Seq() uint32 {
	return binary.LittleEndian.Uint32(h[8:])
}

func (h NetHeader) Command() uint32 {
	return binary.LittleEndian.Uint32(h[12:])
}

func (h NetHeader) Pack(size uint32, flag MsgFlag, msgId, errno, seq uint32) {
	var cmd = msgId
	if errno > 0 {
		cmd = errno
	}
	intToBytes(size, h[:3])
	// h[3:7] = checksum // set after
	h[7] = uint8(flag)
	binary.LittleEndian.PutUint32(h[8:], seq)
	binary.LittleEndian.PutUint32(h[12:], cmd)
}

// CalcChecksum checksum = f(header[7:]) and f(body)
func (h NetHeader) CalcChecksum(body []byte) uint32 {
	var crc = crc32.NewIEEE()
	crc.Write(h[7:])
	if len(body) > 0 {
		crc.Write(body)
	}
	return crc.Sum32()
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

// ReadHeadBody `maxSize`应该小于`MaxPacketSize`
func ReadHeadBody(rd io.Reader, head NetHeader, maxSize uint32) ([]byte, error) {
	if _, err := io.ReadFull(rd, head); err != nil {
		return nil, err
	}
	var nLen = head.Len()
	if nLen < HeaderLength || nLen > maxSize {
		logger.Errorf("ReadHeadBody: msg size %d out of range", nLen)
		return nil, ErrPktSizeOutOfRange
	}
	var body []byte
	if nLen > HeaderLength {
		body = make([]byte, nLen-HeaderLength)
		if _, err := io.ReadFull(rd, body); err != nil {
			return nil, err
		}
	}

	var checksum = head.CalcChecksum(body)
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

func DecodeMsgFrom(rd io.Reader, maxSize uint32, decrypt Encryptor, netMsg *NetMessage) error {
	var head = NewNetHeader()
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
	if flags.Has(FlagError) {
		netMsg.Errno = int32(head.Command())
	} else {
		netMsg.MsgID = head.Command()
	}
	netMsg.Data = body
	return nil
}

func EncodeMsgTo(netMsg *NetMessage, encrypt Encryptor, w io.Writer) error {
	if err := netMsg.Encode(); err != nil {
		return err
	}
	var body = netMsg.Data
	var flags MsgFlag
	if len(body) > DefaultCompressThreshold {
		if encoded, er := compress(body); er != nil {
			logger.Errorf("msg %v compress failed: %v", netMsg.MsgID, er)
		} else {
			if len(encoded) < len(body) {
				flags |= FlagCompress
				body = encoded
			}
		}
	}
	if encrypt != nil && len(body) > 0 {
		if encrypted, err := encrypt.Encrypt(body); err != nil {
			return err
		} else {
			body = encrypted
		}
		flags |= FlagEncrypt
	}

	var bodySize = len(body)
	if bodySize > MaxPayloadSize {
		return fmt.Errorf("encoded msg %v size %d/%d overflow", netMsg.MsgID, bodySize, MaxPayloadSize)
	}

	var size = bodySize + HeaderLength
	var head = NewNetHeader()
	head.Pack(uint32(size), flags, netMsg.MsgID, uint32(netMsg.Errno), netMsg.Seq)

	var checksum = head.CalcChecksum(body)
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
