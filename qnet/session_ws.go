// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"net"

	"github.com/gorilla/websocket"
)

type WebsocketSession struct {
	StreamConnBase
	conn *websocket.Conn
}

func (t *WebsocketSession) UnderlyingConn() net.Conn {
	if t.conn != nil {
		return t.conn.UnderlyingConn()
	}
	return nil
}

func (t *WebsocketSession) Go(reader, writer bool) {
	// TODO:
}

func (t *WebsocketSession) Close() error {
	// TODO:
	return nil
}

func (t *WebsocketSession) ForceClose(err error) {
	// TODO:
}
