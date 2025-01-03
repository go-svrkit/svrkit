// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package cluster

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	etcdHostAddr = os.Getenv("ETCD_ADDR")
	etcdKeyspace = "/choyd"
	nodeId       string
)

func init() {
	nodeId = strconv.Itoa(rand.Int() % 100000)
	println("etcd addr is", etcdHostAddr)
}

func connectClient(t *testing.T) *EtcdClient {
	var client = NewEtcdClient(etcdHostAddr, etcdKeyspace, "", "")
	client.verbose = VerboseLv2
	if err := client.Init(context.Background()); err != nil {
		t.Fatalf("connect server: %v", err)
	}
	return client
}

func TestEtcdClient_PutNode(t *testing.T) {
	var client = connectClient(t)
	defer client.Close()

	var node = createNode(nodeId)
	var name = "bingo/" + nodeId
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := client.PutNode(ctx, name, node, 0); err != nil {
		t.Fatalf("set node: %v\n", err)
	}
}

func TestEtcdClient_GetNode(t *testing.T) {
	var client = connectClient(t)
	defer client.Close()

	var name = "bingo/" + nodeId
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	node, err := client.GetNode(ctx, name)
	if err != nil {
		t.Fatalf("get node: %v\n", err)
	}
	t.Logf("node %s: %v\n", name, node)
}

func TestEtcdClient_IsNodeExist(t *testing.T) {
	var client = connectClient(t)
	defer client.Close()

	var name = "bingo/" + nodeId
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	found, err := client.IsNodeExist(ctx, name)
	if err != nil {
		t.Fatalf("is exist: %v\n", err)
	}
	t.Logf("node %s exist: %v\n", name, found)

	if found {
		if err := client.DeleteKey(ctx, name); err != nil {
			t.Fatalf("delete node: %v\n", err)
		}
	}
}

func TestEtcdClient_ListNodes(t *testing.T) {
	var client = connectClient(t)
	defer client.Close()

	var dir = "service"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	for i := 1; i < 5; i++ {
		var key = fmt.Sprintf("%s/node%d", dir, i)
		var node = createNode(strconv.Itoa(i))
		if err := client.PutNode(ctx, key, node, 0); err != nil {
			t.Fatalf("set node: %v\n", err)
		}
	}

	nodes, err := client.ListNodes(ctx, dir)
	if err != nil {
		t.Fatalf("list dir %s: %v\n", dir, err)
	}
	t.Logf("%d nodes in dir %s", len(nodes), dir)
	for _, node := range nodes {
		t.Logf("  node: %v", node)
	}
}

func createNode(id string) Node {
	var node = NewNode("BINGO", 1)
	return node
}

func TestEtcdClient_RegisterNode(t *testing.T) {
	var client = connectClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var node = createNode(nodeId)
	var name = "bingo/" + nodeId
	var regCtx = NewNodeKeepAliveContext(name, node, 5)
	var err error
	var job = func() {
		t.Logf("try to register %s", name)
		regCtx.LeaseId, err = client.RegisterNode(ctx, name, node, regCtx.TTL)
		if err != nil {
			t.Logf("register: %v\n", err)
		} else {
			if err = client.KeepAlive(ctx, regCtx.stopChan, regCtx.LeaseId); err != nil {
				t.Logf("keepalive: %v", err)
			} else {
				regCtx.LeaseAlive.Store(true)
				t.Logf("register %s with lease %d done", name, regCtx.LeaseId)
			}
		}
	}

	job()
	defer client.DeleteKey(context.Background(), name)

	var ticker = time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	var ticks = 0
	for {
		select {
		case <-ticker.C:
			ticks++
			fmt.Printf("ticks %d\n", ticks)
			if !regCtx.LeaseAlive.Load() {
				fmt.Printf("re-register worker at tick %d, in case of etcd server lost\n", ticks)
				client.DeleteKey(context.Background(), name)
				job()
			}

		case <-regCtx.stopChan:
			regCtx.LeaseAlive.Store(false)
			fmt.Printf("lease %d is dead, try re-register later\n", regCtx.LeaseId)

		case <-ctx.Done():
			return
		}
	}
}

func TestEtcdClient_RegisterAndKeepAliveForever(t *testing.T) {
	var client = connectClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var node = createNode(nodeId)
	var name = "bingo/" + nodeId
	t.Logf("register and keepalive forever, only for 30s")
	regCtx, err := client.RegisterAndKeepAliveForever(ctx, name, node, 5)
	if err != nil {
		t.Fatalf("register forever: %v", err)
	}

	defer client.DeleteKey(context.Background(), name)

	// wait until timed-out
	<-ctx.Done()
	client.RevokeKeepAlive(context.Background(), regCtx)
	t.Logf("done")
}

func TestEtcdClient_WatchDir(t *testing.T) {
	var client = connectClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*40)
	defer cancel()

	var dir = "service"
	var eventChan = client.WatchDir(ctx, dir)

	var ticker = time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	var tick = 0

	var modKey = func() {
		var id = rand.Int() % tick
		if id == 0 {
			id += 1
		}
		var key = fmt.Sprintf("%s/node%d", dir, id)
		var node = createNode(strconv.Itoa(id))
		if err := client.PutNode(ctx, key, node, 0); err != nil {
			t.Fatalf("set node: %v\n", err)
		}
		if tick > 0 && tick%5 == 0 {
			if err := client.DeleteKey(ctx, key); err != nil {
				t.Fatalf("del node: %v\n", err)
			}
		}
	}

	for {
		select {
		case <-ticker.C:
			tick++
			modKey()

		case event, ok := <-eventChan:
			if !ok {
				return
			}
			fmt.Printf("event: %v, key: %s, node: %v\n", event.Type, event.Key, event.Node)

		case <-ctx.Done():
			return
		}
	}
}

func TestEtcdClient_WatchDirTo(t *testing.T) {
	var client = connectClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*40)
	defer cancel()

	var nodeMap = NewNodeMap()
	var dir = "service"

	// list all nodes, and insert to map
	nodes, err := client.ListNodes(ctx, dir)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	for _, node := range nodes {
		nodeMap.AddNode(node)
	}

	client.WatchDirTo(ctx, dir, nodeMap)

	var showNodeMap = func() {
		fmt.Printf("now we have %d nodes\n", nodeMap.Count())
		for _, name := range nodeMap.GetKeys() {
			var nn = nodeMap.GetNodes(name)
			for _, node := range nn {
				fmt.Printf("  %v\n", node)
			}
		}
	}

	showNodeMap()

	var ticker = time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	var tick = 0

	var modKey = func() {
		var id = rand.Int() % tick
		if id == 0 {
			id += 1
		}
		var key = fmt.Sprintf("%s/node%d", dir, id)
		var node = createNode(strconv.Itoa(id))
		if err := client.PutNode(ctx, key, node, 0); err != nil {
			t.Fatalf("set node: %v\n", err)
		}
		if tick > 0 && tick%5 == 0 {
			if err := client.DeleteKey(ctx, key); err != nil {
				t.Fatalf("del node: %v\n", err)
			}
		}
	}

	for {
		select {
		case <-ticker.C:
			tick++
			modKey()

		case <-ctx.Done():
			showNodeMap()
			return
		}
	}
}
