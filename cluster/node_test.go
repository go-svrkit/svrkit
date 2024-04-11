// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package cluster

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	var node = NewNode("GATE", 1)
	node.SetStr("hello", "world")
	assert.Equal(t, "world", node.GetStr("hello"))

	node.SetInt("foo", 123)
	assert.Equal(t, 123, node.GetInt("foo"))

	node.SetBool("boolkey", true)
	assert.Equal(t, true, node.GetBool("boolkey"))

	node.SetFloat("floatkey", 3.14159)
	assert.Equal(t, 3.14159, node.GetFloat("floatkey"))

	var s = node.String()
	t.Logf("node: %v", s)

	var clone = node.Clone()
	clone.SetStr("hello", "clone")
	assert.Equal(t, "world", node.GetStr("hello"))
}

func createNodesMap(text string) *NodeMap {
	var nm = NewNodeMap()
	var parts = strings.Split(text, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		var pair = strings.Split(part, "=")
		if len(pair) == 2 {
			var strKey = strings.TrimSpace(pair[0])
			var strId = strings.TrimSpace(pair[1])
			n, _ := strconv.Atoi(strId)
			nm.AddNode(NewNode(strKey, uint32(n)))
		}
	}
	return nm
}

func TestNodeMap_Count(t *testing.T) {
	tests := []struct {
		input string
		want1 int
		ntype string
		want2 int
	}{
		{"", 0, "", 0},
		{"GATE=1, GATE=2, GAME=1, GAME=2, GAME=3, LOGIN=1", 6, "GATE", 2},
		{"GATE=1, GATE=2, GAME=1, GAME=2, GAME=3, LOGIN=1", 6, "GAME", 3},
		{"GATE=1, GATE=2, GAME=1, GAME=2, GAME=3, LOGIN=1", 6, "LOGIN", 1},
		{"GATE=1, GATE=2, GAME=1, GAME=2, GAME=3, LOGIN=1", 6, "MAIL", 0},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var nm = createNodesMap(tt.input)
			var got1 = nm.Count()
			var got2 = nm.CountOf(tt.ntype)
			if got1 != tt.want1 {
				t.Fatalf("Count() = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Fatalf("CountOf(%s) = %v, want %v", tt.ntype, got2, tt.want2)
			}
		})
	}
}

func TestNodeMap_GetKeys(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"", []string{}},
		{"GATE=1, GATE=2, GAME=1, GAME=2, GAME=3, LOGIN=1, LOGIN=2", []string{"GATE", "GAME", "LOGIN"}},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var nm = createNodesMap(tt.input)
			var got = nm.GetKeys()
			slices.Sort(got)
			slices.Sort(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeMap_FindNodeOf(t *testing.T) {
	tests := []struct {
		input string
		ntype string
		id    uint32
		want  int
	}{
		{"GATE=1, GATE=2, GAME=1, GAME=2, GAME=3, LOGIN=1, LOGIN=2", "MAIL", 1, -1},
		{"GATE=1, GATE=2, GAME=1, GAME=2, GAME=3, LOGIN=1, LOGIN=2", "GATE", 1, 0},
		{"GATE=1, GATE=2, GAME=1, GAME=2, GAME=3, LOGIN=1, LOGIN=2", "GATE", 2, 1},
		{"GATE=1, GATE=2, GAME=1, GAME=2, GAME=3, LOGIN=1, LOGIN=2", "LOGIN", 2, 1},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var nm = createNodesMap(tt.input)
			var got = nm.FindNodeOf(tt.ntype, tt.id)
			if got != tt.want {
				t.Fatalf("FindNodeOf(%s) = %v, want %v", tt.ntype, got, tt.want)
			}
		})
	}
}

func TestNodeMap_AddNode(t *testing.T) {
	var nm = createNodesMap("")
	nm.AddNode(NewNode("GATE", 1))
	nm.AddNode(NewNode("GATE", 1))
	assert.Equal(t, 1, nm.CountOf("GATE"))
	nm.AddNode(NewNode("GATE", 2))
	assert.Equal(t, 2, nm.CountOf("GATE"))
}

func TestNodeMap_DeleteNodesOf(t *testing.T) {
	var nm = createNodesMap("")
	nm.AddNode(NewNode("GATE", 1))
	nm.AddNode(NewNode("GATE", 2))
	nm.AddNode(NewNode("GAME", 1))
	assert.Equal(t, 2, nm.CountOf("GATE"))
	assert.Equal(t, 1, nm.CountOf("GAME"))
	nm.DeleteNodesOf("GATE")
	assert.Equal(t, 0, nm.CountOf("GATE"))
	assert.Equal(t, 1, nm.CountOf("GAME"))
}

func TestNodeMap_DeleteNode(t *testing.T) {
	var nm = createNodesMap("")
	nm.AddNode(NewNode("GATE", 1))
	nm.AddNode(NewNode("GATE", 2))
	nm.AddNode(NewNode("GAME", 1))
	assert.Equal(t, 2, nm.CountOf("GATE"))
	nm.DeleteNode("GATE", 2)
	assert.Equal(t, 1, nm.CountOf("GATE"))
	assert.Equal(t, 1, nm.CountOf("GAME"))
}
