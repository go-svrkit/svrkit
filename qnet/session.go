// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"sync"
	"time"

	"gopkg.in/svrkit.v1/slog"
)

var (
	TCPReadTimeout = 300 * time.Second // 默认读超时
)

type TcpSession struct {
	StreamConnBase
	conn     net.Conn       // tcp Conn
	done     chan struct{}  // sync write pump
	wg       sync.WaitGroup // sync flush write
	intranet bool           // 内联网
}

func NewTcpSession(conn net.Conn) *TcpSession {
	var session = &TcpSession{
		conn: conn,
		done: make(chan struct{}),
	}
	session.RemoteAddr = conn.RemoteAddr().String()
	return session
}

func (t *TcpSession) UnderlyingConn() net.Conn {
	return t.conn
}

func (t *TcpSession) SetIntranet(v bool) {
	t.intranet = v
}

func (t *TcpSession) Go(reader, writer bool) {
	if reader || writer {
		t.Running.Store(true)
	}
	if writer {
		t.wg.Add(1)
		go t.writePump()
	}
	if reader {
		t.wg.Add(1)
		go t.readPump()
	}
}

func (t *TcpSession) Close() error {
	if !t.Running.CompareAndSwap(true, false) {
		return nil // close in progress
	}
	if conn, ok := t.conn.(*net.TCPConn); ok {
		if err := conn.CloseRead(); err != nil {
			slog.Infof("%v close read: %v", t.Node, err)
		}
	}
	t.finally(ErrConnForceClose) // 阻塞等待投递剩余的消息
	return nil
}

func (t *TcpSession) ForceClose(reason error) {
	if !t.Running.CompareAndSwap(true, false) {
		return // close in progress
	}

	if conn, ok := t.conn.(*net.TCPConn); ok {
		if err := conn.CloseRead(); err != nil {
			slog.Infof("%v close read: %v", t.Node, err)
		}
	}
	go t.finally(reason) // 不阻塞等待
}

func (t *TcpSession) notifyErr(reason error) {
	if t.ErrChan != nil {
		var err = NewError(reason, t)
		select {
		case t.ErrChan <- err:
		default:
			return
		}
	}
}

func (t *TcpSession) finally(reason error) {
	close(t.done) // 通知发送线程flush并退出
	t.notifyErr(reason)
	t.wg.Wait()

	t.conn = nil
	t.RecvQueue = nil
	t.SendQueue = nil
	t.ErrChan = nil
	t.Encrypt = nil
	t.Decrypt = nil
	t.Userdata = nil
}

func (t *TcpSession) flush() {
	defer t.conn.Close() // close after flush, and this stops reader pump

	var buf bytes.Buffer
	t.conn.SetWriteDeadline(time.Now().Add(time.Minute))
	for {
		select {
		case netMsg, ok := <-t.SendQueue:
			if !ok {
				return
			}
			buf.Reset()
			if err := t.write(netMsg, &buf); err != nil {
				slog.Errorf("%v flush message %v: %v", t.Node, netMsg.Command, err)
			}
		default:
			return
		}
	}
}

func (t *TcpSession) write(netMsg *NetMessage, buf *bytes.Buffer) error {
	if err := EncodeMsgTo(netMsg, t.Decrypt, buf); err != nil {
		return err
	}
	if _, err := t.conn.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (t *TcpSession) writePump() {
	defer func() {
		t.flush()
		t.wg.Done()
		slog.Debugf("TcpSession: node %v writer stopped", t.Node)
	}()

	//slog.Debugf("TcpSession: node %v(%v) writer started", t.node, t.addr)
	var buf = new(bytes.Buffer)
	for {
		select {
		case netMsg, ok := <-t.SendQueue:
			if !ok {
				return
			}
			// reuse small size buffer
			if buf.Cap() < 16<<10 {
				buf.Reset()
			} else {
				buf = new(bytes.Buffer)
			}
			if err := t.write(netMsg, buf); err != nil {
				slog.Errorf("%v write message %v: %v", t.Node, netMsg.Command, err)
				continue
			}

		case <-t.done:
			return
		}
	}
}

func (t *TcpSession) readMessage(rd io.Reader, netMsg *NetMessage) error {
	var maxBytes uint32 = MaxPayloadSize
	if !t.intranet {
		maxBytes = uint32(MaxClientUpStreamSize)
	}
	var deadline = time.Now().Add(TCPReadTimeout)
	if err := t.conn.SetReadDeadline(deadline); err != nil {
		slog.Errorf("session %v set read deadline: %v", t.Node, err)
	}
	if err := DecodeMsgFrom(rd, maxBytes, t.Decrypt, netMsg); err != nil {
		return err
	}
	netMsg.Session = t
	return nil
}

func (t *TcpSession) readPump() {
	defer func() {
		t.wg.Done()
		slog.Debugf("TcpSession: node %v reader stopped", t.Node)
	}()

	//slog.Debugf("TcpSession: node %v(%v) reader started", t.node, t.addr)
	var rd = bufio.NewReader(t.conn)
	for t.IsRunning() {
		var netMsg = AllocNetMessage()
		if err := t.readMessage(rd, netMsg); err != nil {
			if err != io.EOF {
				slog.Errorf("session %v read packet %v", t.Node, err)
			}
			t.ForceClose(err) // I/O超时或者发生错误，强制关闭连接
			return
		}
		netMsg.Session = t
		netMsg.CreatedAt = time.Now().UnixMicro()
		// 如果channel满了，不能丢弃，需要阻塞等待
		t.RecvQueue <- netMsg
	}
}
