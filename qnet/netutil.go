// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"gopkg.in/svrkit.v1/codec"
	"gopkg.in/svrkit.v1/strutil"
	"gopkg.in/svrkit.v1/zlog"
)

const (
	UrlFormKey = "pb-data"
)

const PrefixLength = 2 // sizeof uint16

// ReadLenData 读取长度[2字节]开头的数据
func ReadLenData(r io.Reader, maxSize uint16) ([]byte, error) {
	var tmp [PrefixLength]byte
	if _, err := io.ReadFull(r, tmp[:]); err != nil {
		return nil, fmt.Errorf("ReadLenPrefixData: read len: %v", err)
	}
	var nLen = int(binary.BigEndian.Uint16(tmp[:]))
	if nLen < PrefixLength || nLen > int(maxSize) {
		return nil, fmt.Errorf("ReadLenPrefixData: msg size %d out of range", nLen)
	}
	var data []byte
	if nLen > PrefixLength {
		data = make([]byte, nLen-PrefixLength)
		if _, err := io.ReadFull(r, data); err != nil {
			return nil, fmt.Errorf("ReadLenPrefixData: read body of len %d: %v", nLen, err)
		}
	}
	return data, nil
}

// WriteLenData 写入长度[2字节]开头的数据
func WriteLenData(w io.Writer, body []byte) error {
	if len(body) == 0 {
		return nil
	}
	var nLen = len(body) + PrefixLength
	if nLen > math.MaxUint16 {
		return fmt.Errorf("WriteLenPrefixData: msg size %d out of range", nLen)
	}
	var buf [PrefixLength]byte
	binary.BigEndian.PutUint16(buf[:PrefixLength], uint16(nLen))

	if _, err := w.Write(buf[:]); err != nil {
		return err
	}
	if _, err := w.Write(body); err != nil {
		return err
	}
	return nil
}

func ReadProtoMessage(conn net.Conn, msg codec.Message) error {
	conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	body, err := ReadLenData(conn, math.MaxUint16)
	if err != nil {
		return err
	}
	if err = codec.Unmarshal(body, msg); err != nil {
		return err
	}
	return nil
}

func WriteProtoMessage(w io.Writer, msg codec.Message) error {
	body, err := codec.Marshal(msg)
	if err != nil {
		return err
	}
	return WriteLenData(w, body)
}

// RequestProtoMessage send req and wait for ack
func RequestProtoMessage(conn net.Conn, req, ack codec.Message) error {
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
func ReadProtoFromHTTPRequest(req *http.Request, msg codec.Message) error {
	var contentType = req.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}

	switch contentType {
	case "application/json":
		defer req.Body.Close()
		if data, err := io.ReadAll(req.Body); err != nil {
			return err
		} else {
			return codec.UnmarshalProtoJSON(data, msg)
		}

	case "application/octet-stream", "application/binary":
		rawbytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		defer req.Body.Close()
		if err = codec.Unmarshal(rawbytes, msg); err != nil {
			return err
		}
		return nil

	case "application/x-www-form-urlencoded":
		if err := req.ParseForm(); err != nil {
			return err
		}
		var data = req.Form.Get(UrlFormKey)
		if len(data) > 0 {
			return codec.UnmarshalProtoJSON(strutil.StrAsBytes(data), msg)
		}
		return nil

	default:
		return fmt.Errorf("unsupported content type %s", contentType)
	}
}

// WriteProtoHTTPResponse 写入proto消息到http响应
func WriteProtoHTTPResponse(w http.ResponseWriter, msg codec.Message, contentType string) error {
	var data []byte
	switch contentType {
	case "json":
		var err error
		data, err = codec.MarshalProtoJSON(msg)
		if err != nil {
			zlog.Errorf("WriteProtoResponse: MarshalProtoJSON: %v", err)
			return err
		}
		w.Header().Set("Content-Type", "application/json")

	default:
		rawbytes, err := codec.Marshal(msg)
		if err != nil {
			zlog.Errorf("WriteProtoResponse: pb.Marshal: %v", err)
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
		zlog.Errorf("cannot fetch net interfaces: %v", err)
		return nil
	}
	var ipList []net.IP
	for _, inetface := range inetfaces {
		addrs, err := inetface.Addrs()
		if err != nil {
			zlog.Errorf("cannot fetch net address: %v", err)
			continue
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
