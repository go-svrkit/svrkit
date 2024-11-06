package qnet

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeBackendNode(t *testing.T) {
	tests := []struct {
		service  uint16
		instance uint32
	}{
		{1, 1},
		{math.MaxUint8, math.MaxUint16},
		{math.MaxUint16, math.MaxUint32},
	}
	for _, tt := range tests {
		var node = MakeBackendNode(tt.service, tt.instance)
		assert.True(t, node.IsBackend())
		assert.False(t, node.IsSession())
		assert.Equal(t, tt.service, node.ServiceType())
		assert.Equal(t, tt.instance, node.InstanceID())
	}
}

func TestMakeGateSession(t *testing.T) {
	tests := []struct {
		instance uint16
		session  uint32
	}{
		{1, 1},
		{math.MaxUint8, math.MaxUint32},
	}
	for _, tt := range tests {
		var node = MakeGateSession(tt.instance, tt.session)
		assert.False(t, node.IsBackend())
		assert.True(t, node.IsSession())
		assert.Equal(t, tt.instance, node.GateID())
		assert.Equal(t, tt.session, node.Session())
	}
}

func TestNodeID_String(t *testing.T) {
	tests := []struct {
		service  uint16
		instance uint32
		expected string
	}{
		{1, 1, "01#1"},
		{math.MaxUint8, math.MaxUint32, "FF#4294967295"},
	}
	for _, tt := range tests {
		var node = MakeBackendNode(tt.service, tt.instance)
		assert.Equal(t, node.String(), tt.expected)
	}
}
