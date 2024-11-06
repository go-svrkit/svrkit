// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package qnet

import (
	"encoding/json"

	"github.com/gogo/protobuf/proto"
)

type WsRecvMsg struct {
	Cmd  string          `json:"cmd"`
	Body json.RawMessage `json:"body"`
}

type WsWriteMsg struct {
	Cmd  string        `json:"cmd"`
	Body proto.Message `json:"body"`
}

type WebsockSession struct {
}

var _ Endpoint = (*TcpSession)(nil)
