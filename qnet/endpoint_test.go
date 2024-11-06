package qnet

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
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
	var m = NewEndpointMap()
	assert.Equal(t, m.Size(), 0)
	m.Put(1, &FakeSession{})
	assert.Equal(t, m.Size(), 1)
}

func TestEndpointMap_IsEmpty(t *testing.T) {
	var m = NewEndpointMap()
	assert.True(t, m.IsEmpty())
	m.Put(1, &FakeSession{})
	assert.False(t, m.IsEmpty())
}

func TestEndpointMap_Has(t *testing.T) {
	var m = NewEndpointMap()
	assert.False(t, m.Has(1))
	m.Put(1, &FakeSession{})
	assert.True(t, m.Has(1))
}
