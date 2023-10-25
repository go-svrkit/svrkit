// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"gopkg.in/svrkit.v1/logger"
)

const (
	UrlFormKey = "pb-data"
)

const PrefixLength = 4 // sizeof uint32

// ReadLenPrefixData 读取长度[4字节]开头的数据
func ReadLenPrefixData(r io.Reader, maxSize uint) ([]byte, error) {
	var tmp [PrefixLength]byte
	if _, err := io.ReadFull(r, tmp[:]); err != nil {
		return nil, fmt.Errorf("read len: %v", err)
	}
	var nLen = int(binary.BigEndian.Uint32(tmp[:]))
	if nLen < PrefixLength || nLen > int(maxSize) {
		return nil, fmt.Errorf("ReadLenPrefixData: msg size %d out of range", nLen)
	}
	var data []byte
	if nLen > PrefixLength {
		data = make([]byte, nLen-4)
		if _, err := io.ReadFull(r, data); err != nil {
			return nil, fmt.Errorf("ReadLenPrefixData: read body of len %d: %v", nLen, err)
		}
	}
	return data, nil
}

// WriteLenPrefixData 写入长度[4字节]开头的数据
func WriteLenPrefixData(w io.Writer, body []byte) error {
	if len(body) == 0 {
		return nil
	}
	var nLen = len(body) + PrefixLength
	if nLen > MaxPacketSize {
		return fmt.Errorf("WriteLenPrefixData: msg size %d out of range", nLen)
	}
	var buf [PrefixLength]byte
	binary.BigEndian.PutUint32(buf[:PrefixLength], uint32(nLen))

	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	if _, err := w.Write(body); err != nil {
		return err
	}
	return nil
}

func ReadProtoMessage(conn net.Conn, msg proto.Message) error {
	conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	body, err := ReadLenPrefixData(conn, MaxClientUpStreamSize)
	if err != nil {
		return err
	}
	if err = proto.Unmarshal(body, msg); err != nil {
		return err
	}
	return nil
}

func WriteProtoMessage(w io.Writer, msg proto.Message) error {
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return WriteLenPrefixData(w, body)
}

// RequestProtoMessage send req and wait for ack
func RequestProtoMessage(conn net.Conn, req, ack proto.Message) error {
	if err := WriteProtoMessage(conn, req); err != nil {
		return err
	}
	return ReadProtoMessage(conn, ack)
}

// GetHTTPRequestIP 获取http请求的来源IP
func GetHTTPRequestIP(req *http.Request) string {
	ip := req.Header.Get("X-Forwarded-For")
	ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	if ip == "" {
		realIP := req.Header.Get("X-Real-IP")
		ip = strings.TrimSpace(realIP)
	}
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// DecodeHTTPRequestBody 解析http请求的body为json
func DecodeHTTPRequestBody(req *http.Request, ptr interface{}) error {
	d := json.NewDecoder(req.Body)
	d.UseNumber()
	defer req.Body.Close()
	if err := d.Decode(ptr); err != nil {
		return err
	}
	return nil
}

// ReadProtoFromHTTPRequest 从http请求中读取proto消息
func ReadProtoFromHTTPRequest(req *http.Request, msg proto.Message) error {
	var contentType = req.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}

	switch contentType {
	case "application/json":
		defer req.Body.Close()
		return jsonpb.Unmarshal(req.Body, msg)

	case "application/octet-stream", "application/binary":
		rawbytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		defer req.Body.Close()
		if err = proto.Unmarshal(rawbytes, msg); err != nil {
			return err
		}
		return nil

	case "application/x-www-form-urlencoded":
		if err := req.ParseForm(); err != nil {
			return err
		}
		var data = req.Form.Get(UrlFormKey)
		if len(data) > 0 {
			var buf = bytes.NewBufferString(data)
			return jsonpb.Unmarshal(buf, msg)
		}
		return nil

	default:
		return fmt.Errorf("unsupported content type %s", contentType)
	}
}

// WriteProtoHTTPResponse 写入proto消息到http响应
func WriteProtoHTTPResponse(w http.ResponseWriter, msg proto.Message, contentType string) error {
	var data []byte
	switch contentType {
	case "json":
		var m jsonpb.Marshaler
		var buf bytes.Buffer
		if err := m.Marshal(&buf, msg); err != nil {
			logger.Errorf("WriteProtoResponse: jsonpb.Marshal: %v", err)
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		data = buf.Bytes()

	default:
		rawbytes, err := proto.Marshal(msg)
		if err != nil {
			logger.Errorf("WriteProtoResponse: pb.Marshal: %v", err)
			return err
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		data = rawbytes
	}
	_, err := w.Write(data)
	return err
}

// GetLocalIPList 获取本地IP列表
func GetLocalIPList() []net.IP {
	inetfaces, err := net.Interfaces()
	if err != nil {
		logger.Errorf("cannot fetch net interfaces: %v", err)
		return nil
	}
	var ipList []net.IP
	for _, inetface := range inetfaces {
		addrs, err := inetface.Addrs()
		if err != nil {
			logger.Errorf("cannot fetch net address: %v", err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if len(ip) > 0 && !ip.IsLoopback() {
				ipList = append(ipList, ip)
			}
		}
	}
	return ipList
}
