// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package qnet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/websocket"

	"gopkg.in/svrkit.v1/qlog"
)

type WsRecvMsg struct {
	Cmd  string          `json:"cmd"`
	Seq  uint32          `json:"seq,omitempty"`
	Body json.RawMessage `json:"body,omitempty"`
}

type WsWriteMsg struct {
	Cmd  string        `json:"cmd"`
	Seq  uint32        `json:"seq,omitempty"`
	Body proto.Message `json:"body,omitempty"`
}

type WebsockSession struct {
	StreamConnBase

	wg   sync.WaitGroup  // sync flush write
	done chan struct{}   // sync write pump
	conn *websocket.Conn // websocket connection
}

var _ Endpoint = (*WebsockSession)(nil)

func NewWebsockSession(conn *websocket.Conn, sendQSize int) *WebsockSession {
	var session = &WebsockSession{
		conn: conn,
		done: make(chan struct{}),
	}
	if sendQSize > 0 {
		session.SendQueue = make(chan *NetMessage, sendQSize)
	}
	session.RemoteAddr = conn.RemoteAddr().String()
	return session
}

func (t *WebsockSession) UnderlyingConn() net.Conn {
	if t.conn != nil {
		return t.conn.NetConn()
	}
	return nil
}

func (t *WebsockSession) SetIntranet(v bool) {
	// do nothing
}

func (t *WebsockSession) notifyErr(reason error) {
	if t.ErrChan != nil {
		var err = NewError(reason, t)
		select {
		case t.ErrChan <- err:
		default:
			return
		}
	}
}

func (t *WebsockSession) Close() error {
	if !t.Running.CompareAndSwap(true, false) {
		return nil // close in progress
	}
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
	close(t.done)
	t.notifyErr(ErrConnClosed)
	t.wg.Wait()
	return nil
}

func (t *WebsockSession) ForceClose(reason error) {
	if !t.Running.CompareAndSwap(true, false) {
		return // close in progress
	}
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
	close(t.done)
	t.notifyErr(reason)
	t.wg.Wait()
}

func (t *WebsockSession) Go(ctx context.Context, reader, writer bool) {
	if reader || writer {
		t.Running.Store(true)
	}
	if writer {
		t.wg.Add(1)
		go t.writePump(ctx)
	}
	if reader {
		t.wg.Add(1)
		go t.readPump(ctx)
	}
}

func (t *WebsockSession) write(netMsg *NetMessage) error {
	if t.conn == nil {
		return ErrConnClosed
	}
	var wsMsg = &WsWriteMsg{
		Seq:  netMsg.Seq,
		Cmd:  GetMessageShortName(netMsg.Command),
		Body: netMsg.Body,
	}
	data, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}
	if err = t.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}
	return nil
}

func (t *WebsockSession) writePump(ctx context.Context) {
	defer func() {
		t.wg.Done()
		qlog.Debugf("WebsockSession: node %v writer stopped", t.Node)
	}()

	for {
		select {
		case netMsg, ok := <-t.SendQueue:
			if !ok {
				return
			}
			if err := t.write(netMsg); err != nil {
				qlog.Errorf("%v write message %v: %v", t.Node, netMsg.Command, err)
				continue
			}

		case <-t.done:
			return

		case <-ctx.Done():
			return
		}
	}
}

func (t *WebsockSession) ReadMessage(netMsg *NetMessage) error {
	if t.conn == nil {
		return ErrConnClosed
	}
	var deadline = time.Now().Add(TCPReadTimeout)
	if err := t.conn.SetReadDeadline(deadline); err != nil {
		qlog.Errorf("session %v set read deadline: %v", t.Node, err)
	}
	msgType, data, err := t.conn.ReadMessage()
	if err != nil {
		return err
	}
	if msgType != websocket.TextMessage {
		return ErrInvalidWSMsgType
	}

	var wsMsg WsRecvMsg
	var dec = json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err = dec.Decode(&wsMsg); err != nil {
		return err
	}
	var pbMsg = CreateMessageByShortName(wsMsg.Cmd)
	if pbMsg == nil {
		return fmt.Errorf("cannot create msg %s", wsMsg.Cmd)
	}
	if len(wsMsg.Body) > 0 {
		if err = jsonpb.Unmarshal(bytes.NewReader(wsMsg.Body), pbMsg); err != nil {
			return err
		}
	}
	netMsg.Command = GetMessageId(MessagePackagePrefix + wsMsg.Cmd)
	netMsg.Body = pbMsg
	netMsg.Seq = wsMsg.Seq
	netMsg.Session = t
	netMsg.CreatedAt = time.Now().UnixMicro()
	return nil
}

func (t *WebsockSession) readPump(ctx context.Context) {
	defer func() {
		t.wg.Done()
		qlog.Debugf("WebsockSession: node %v reader stopped", t.Node)
	}()

	for t.IsRunning() {
		var netMsg = AllocNetMessage()
		if err := t.ReadMessage(netMsg); err != nil {
			if err != io.EOF {
				qlog.Errorf("session %v read packet %v", t.Node, err)
			}
			t.ForceClose(err) // I/O超时或者发生错误，强制关闭连接
			return
		}
		netMsg.Session = t
		netMsg.CreatedAt = time.Now().UnixMicro()
		// 如果channel满了，不能丢弃，需要阻塞等待
		t.RecvQueue <- netMsg

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
