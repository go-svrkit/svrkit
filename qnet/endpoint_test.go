package qnet

import (
	"context"
	"net"
	"testing"
)

type FakeSession struct {
	StreamConnBase
}

var _ Endpoint = (*FakeSession)(nil)

func NewFakeSession(queue chan *NetMessage) *FakeSession {
	return &FakeSession{
		StreamConnBase: StreamConnBase{
			SendQueue: queue,
		},
	}
}

func (s *FakeSession) GetNode() NodeID {
	return 0
}

func (s *FakeSession) SetNode(NodeID) {

}

func (s *FakeSession) GetRemoteAddr() string {
	return ""
}

func (s *FakeSession) UnderlyingConn() net.Conn {
	return nil
}

func (s *FakeSession) Go(ctx context.Context, reader, writer bool) {

}

func (s *FakeSession) Close() error {
	return nil
}

func (s *FakeSession) ForceClose(error) {

}

func TestEndpointMap_Size(t *testing.T) {
}
