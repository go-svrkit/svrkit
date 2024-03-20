// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package cluster

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/svrkit.v1/zlog"
)

var (
	ErrEmptyLease    = errors.New("empty lease")
	ErrNodeKeyExist  = errors.New("node key exist")
	ErrNoKeyDeleted  = errors.New("no key deleted")
	ErrEmptyEndpoint = errors.New("empty endpoint")
)

const (
	EventChanCapacity = 256
	OpTimeout         = 5

	VerboseLv1 = 1
	VerboseLv2 = 2
)

// EtcdClient 基于etcd的服务发现
type EtcdClient struct {
	closing   atomic.Bool      //
	verbose   int32            //
	endpoints []string         // etcd server address
	namespace string           // name space of key
	username  string           //
	passwd    string           //
	client    *clientv3.Client // etcd client
}

func NewEtcdClient(endpoints, namespace, username, passwd string) *EtcdClient {
	d := &EtcdClient{
		endpoints: strings.Split(endpoints, ","),
		namespace: namespace,
		username:  username,
		passwd:    passwd,
		verbose:   VerboseLv1,
	}
	return d
}

func (c *EtcdClient) Init(parentCtx context.Context) error {
	if len(c.endpoints) == 0 {
		return ErrEmptyEndpoint
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.endpoints,
		DialTimeout: 5 * time.Second,
		Username:    c.username,
		Password:    c.passwd,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(parentCtx, time.Second*3)
	defer cancel()

	if _, err := cli.Status(ctx, c.endpoints[0]); err != nil {
		return fmt.Errorf("get status of etcd %s: %w", c.endpoints[0], err)
	}
	c.client = cli
	return nil
}

func (c *EtcdClient) Close() {
	if !c.closing.CompareAndSwap(false, true) {
		return
	}
	if c.client != nil {
		c.client.Close()
		c.client = nil
	}
}

func (c *EtcdClient) SetVerbose(v int32) {
	c.verbose = v
}

func (c *EtcdClient) FormatKey(name string) string {
	return path.Join(c.namespace, name)
}

// IsNodeExist 节点是否存在
func (c *EtcdClient) IsNodeExist(ctx context.Context, name string) (bool, error) {
	var key = c.FormatKey(name)
	resp, err := c.client.Get(ctx, key, clientv3.WithCountOnly())
	if err != nil {
		return false, err
	}
	return resp.Count > 0, nil
}

// GetKey 获取key的值
func (c *EtcdClient) GetKey(ctx context.Context, name string) ([]byte, error) {
	var key = c.FormatKey(name)
	resp, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, nil
	}
	return resp.Kvs[0].Value, nil
}

// GetNode 获取节点信息
func (c *EtcdClient) GetNode(ctx context.Context, name string) (*Node, error) {
	var key = c.FormatKey(name)
	resp, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, nil
	}
	var node = new(Node)
	if err := unmarshalNode(resp.Kvs[0].Value, node); err != nil {
		return nil, err
	}
	return node, nil
}

// PutNode 设置节点信息
func (c *EtcdClient) PutNode(ctx context.Context, name string, value any, leaseId int64) error {
	var key = c.FormatKey(name)
	data, err := sonic.MarshalString(value)
	if err != nil {
		return err
	}
	var resp *clientv3.PutResponse
	if leaseId <= 0 {
		resp, err = c.client.Put(ctx, key, data)
	} else {
		resp, err = c.client.Put(ctx, key, data, clientv3.WithLease(clientv3.LeaseID(leaseId)))
	}
	if err != nil {
		return err
	}
	if c.verbose >= VerboseLv1 {
		zlog.Infof("put key [%s] at rev %d", key, resp.Header.Revision)
	}
	return nil
}

// DelKey 删除一个key
func (c *EtcdClient) DelKey(ctx context.Context, name string) error {
	var key = c.FormatKey(name)
	resp, err := c.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	if resp.Deleted == 0 {
		return ErrNoKeyDeleted
	}
	return nil
}

// ListNodes 列出目录下的所有节点
func (c *EtcdClient) ListNodes(ctx context.Context, prefix string) ([]Node, error) {
	var key = c.FormatKey(prefix)
	resp, err := c.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, nil
	}
	var nodes = make([]Node, 0, resp.Count)
	for _, kv := range resp.Kvs {
		var node Node
		if err := unmarshalNode(kv.Value, &node); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// GrantLease 申请一个lease
func (c *EtcdClient) GrantLease(ctx context.Context, ttl int32) (int64, error) {
	lease, err := c.client.Grant(ctx, int64(ttl))
	if err != nil {
		return 0, err
	}
	if lease == nil {
		return 0, ErrEmptyLease
	}
	return int64(lease.ID), nil
}

func (c *EtcdClient) GetLeaseTTL(ctx context.Context, leaseId int64) (int, error) {
	resp, err := c.client.TimeToLive(ctx, clientv3.LeaseID(leaseId))
	if err != nil {
		return 0, nil
	}
	return int(resp.TTL), nil
}

// RevokeLease 撤销一个lease
func (c *EtcdClient) RevokeLease(ctx context.Context, leaseId int64) error {
	_, err := c.client.Revoke(ctx, clientv3.LeaseID(leaseId))
	return err
}

// NodeKeepAliveContext 用于注册并保活节点
type NodeKeepAliveContext struct {
	stopChan   chan struct{}
	LeaseId    int64
	LeaseAlive atomic.Bool
	TTL        int32
	Name       string
	Value      any
}

func NewNodeKeepAliveContext(name string, value any, ttl int32) *NodeKeepAliveContext {
	return &NodeKeepAliveContext{
		stopChan: make(chan struct{}, 1),
		Name:     name,
		Value:    value,
		TTL:      ttl,
	}
}

func (c *EtcdClient) RevokeKeepAlive(ctx context.Context, regCtx *NodeKeepAliveContext) error {
	if c.verbose >= VerboseLv1 {
		zlog.Infof("try revoke node %s lease %d", regCtx.Name, regCtx.LeaseId)
	}
	if regCtx.LeaseId == 0 || !regCtx.LeaseAlive.Load() {
		if c.verbose >= VerboseLv1 {
			zlog.Infof("node %s lease %d is not alive", regCtx.Name, regCtx.LeaseId)
		}
		return nil
	}
	if err := c.RevokeLease(ctx, regCtx.LeaseId); err != nil {
		zlog.Warnf("revoke node %s lease %x failed: %v", regCtx.Name, regCtx.LeaseId, err)
		return err
	} else {
		if c.verbose >= VerboseLv1 {
			zlog.Infof("revoke node %s lease %x done", regCtx.Name, regCtx.LeaseId)
		}
	}
	return nil
}

// RegisterNode 注册一个节点信息，并返回一个ttl秒的lease
func (c *EtcdClient) RegisterNode(rootCtx context.Context, name string, value any, ttl int32) (int64, error) {
	ctx, cancel := context.WithTimeout(rootCtx, time.Second*OpTimeout)
	defer cancel()

	exist, err := c.IsNodeExist(ctx, name)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, ErrNodeKeyExist
	}
	var leaseId int64
	if ttl <= 0 {
		ttl = 7
	}
	if leaseId, err = c.GrantLease(ctx, ttl); err != nil {
		return 0, err
	}
	if err = c.PutNode(ctx, name, value, leaseId); err != nil {
		return 0, err
	}
	return leaseId, nil
}

func revokeLeaseWithTimeout(c *EtcdClient, leaseId int64) {
	if c.verbose >= VerboseLv1 {
		zlog.Infof("try revoke lease %d", leaseId)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*OpTimeout)
	defer cancel()
	if err := c.RevokeLease(ctx, leaseId); err != nil {
		zlog.Warnf("revoke lease %x failed: %v", leaseId, err)
	} else {
		zlog.Infof("revoke lease %x done", leaseId)
	}
}

func (c *EtcdClient) aliveKeeper(ctx context.Context, kaChan <-chan *clientv3.LeaseKeepAliveResponse, stopChan chan struct{}, leaseId int64) {
	defer func() {
		select {
		case stopChan <- struct{}{}:
		default:
			break
		}
	}()

	for {
		select {
		case ka, ok := <-kaChan:
			if !ok || ka == nil {
				zlog.Infof("lease %d is not alive", leaseId)
				return
			}
			if c.verbose >= VerboseLv2 {
				zlog.Infof("lease %d respond alive, ttl %d", ka.ID, ka.TTL)
			}

		case <-ctx.Done():
			zlog.Infof("stop keepalive with lease %d", leaseId)
			return
		}
	}
}

// KeepAlive lease保活，当lease撤销时此stopChan被激活
func (c *EtcdClient) KeepAlive(ctx context.Context, stopChan chan struct{}, leaseId int64) error {
	kaChan, err := c.client.KeepAlive(ctx, clientv3.LeaseID(leaseId))
	if err != nil {
		return nil
	}
	go c.aliveKeeper(ctx, kaChan, stopChan, leaseId)
	return nil
}

func (c *EtcdClient) doRegisterNode(ctx context.Context, regCtx *NodeKeepAliveContext) error {
	var err error
	if c.verbose >= VerboseLv1 {
		zlog.Infof("try register key: %s", c.FormatKey(regCtx.Name))
	}
	regCtx.LeaseAlive.Store(false)
	regCtx.LeaseId = 0

	regCtx.LeaseId, err = c.RegisterNode(ctx, regCtx.Name, regCtx.Value, regCtx.TTL)
	if err != nil {
		return err
	}
	if err = c.KeepAlive(ctx, regCtx.stopChan, regCtx.LeaseId); err != nil {
		return err
	}
	regCtx.LeaseAlive.Store(true)
	if c.verbose >= VerboseLv1 {
		zlog.Infof("register key [%s] with lease %x done", c.FormatKey(regCtx.Name), regCtx.LeaseId)
	}
	return nil
}

func (c *EtcdClient) regAliveKeeper(ctx context.Context, regCtx *NodeKeepAliveContext) {
	var ticker = time.NewTicker(time.Second) // 1s
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if !regCtx.LeaseAlive.Load() {
				if err := c.doRegisterNode(ctx, regCtx); err != nil {
					zlog.Infof("register or keepalive %s failed: %v", regCtx.Name, err)
				}
			}

		case <-regCtx.stopChan:
			var leaseId = regCtx.LeaseId
			regCtx.LeaseAlive.Store(false)
			regCtx.LeaseId = 0
			if c.verbose >= VerboseLv1 {
				zlog.Infof("node %s lease(%d) is not alive, try register later", regCtx.Name, leaseId)
			}

		case <-ctx.Done():
			if c.verbose >= VerboseLv1 {
				zlog.Infof("registration alive keeper with key %s stopped", regCtx.Name)
			}
			return
		}
	}
}

// RegisterAndKeepAliveForever 注册一个节点，并永久保活
func (c *EtcdClient) RegisterAndKeepAliveForever(ctx context.Context, name string, value any, ttl int32) (*NodeKeepAliveContext, error) {
	var regCtx = NewNodeKeepAliveContext(name, value, ttl)
	if err := c.doRegisterNode(ctx, regCtx); err != nil {
		return nil, err
	}
	go c.regAliveKeeper(ctx, regCtx)
	return regCtx, nil
}

func propagateWatchEvent(eventChan chan<- *NodeEvent, ev *clientv3.Event) {
	var event = &NodeEvent{
		Type: EventUnknown,
		Key:  string(ev.Kv.Key),
	}
	switch ev.Type {
	case 0: // PUT
		if ev.IsCreate() {
			event.Type = EventCreate
		} else {
			event.Type = EventUpdate
		}
	case 1: // DELETE
		event.Type = EventDelete
	}
	if len(ev.Kv.Value) > 0 {
		if err := unmarshalNode(ev.Kv.Value, &event.Node); err != nil {
			zlog.Errorf("unmarshal node %s: %v", event.Key, err)
			return
		}
	}

	select {
	case eventChan <- event:
	default:
		zlog.Warnf("watch event channel is full, new event lost: %v", event)
	}
}

// WatchDir 订阅目录下的节点变化
func (c *EtcdClient) WatchDir(ctx context.Context, dir string) <-chan *NodeEvent {
	var key = c.FormatKey(dir)
	var watchCh = c.client.Watch(clientv3.WithRequireLeader(ctx), key, clientv3.WithPrefix())
	var eventChan = make(chan *NodeEvent, EventChanCapacity)
	var watcher = func() {
		defer close(eventChan)
		for {
			select {
			case resp, ok := <-watchCh:
				if !ok {
					return
				}
				if resp.Err() != nil {
					zlog.Warnf("watch key %s canceled: %v", key, resp.Err())
					return
				}
				for _, ev := range resp.Events {
					propagateWatchEvent(eventChan, ev)
				}

			case <-ctx.Done():
				if c.client != nil {
					if err := c.client.Watcher.Close(); err != nil {
						zlog.Warnf("close watcher: %v", err)
					}
				}
				return
			}
		}
	}
	go watcher()
	return eventChan
}

// WatchDirTo 订阅目录下的所有节点变化, 并把节点变化更新到nodeMap
func (c *EtcdClient) WatchDirTo(ctx context.Context, dir string, nodeMap *NodeMap) {
	var evChan = c.WatchDir(ctx, dir)
	var prefix = c.FormatKey(dir)
	var watcher = func() {
		for ev := range evChan {
			updateNodeEvent(nodeMap, prefix, ev)
		}
	}
	go watcher()
}

func updateNodeEvent(nodeMap *NodeMap, rootDir string, ev *NodeEvent) {
	switch ev.Type {
	case EventCreate:
		nodeMap.AddNode(ev.Node)
	case EventUpdate:
		nodeMap.AddNode(ev.Node) // 插入前会先检查是否有重复
	case EventDelete:
		nodeType, id := parseNodeTypeAndID(rootDir, ev.Key)
		if nodeType != "" && id > 0 {
			nodeMap.DeleteNode(nodeType, id)
		}
	}
}

func parseNodeTypeAndID(root, key string) (string, uint32) {
	var idx = strings.Index(key, root)
	if idx < 0 {
		return "", 0
	}
	key = key[len(root)+1:] // root + '/' + key
	idx = strings.Index(key, "/")
	if idx <= 0 {
		return "", 0
	}
	var nodeType = key[:idx]
	var strId = key[idx+1:]
	id, _ := strconv.Atoi(strId)
	return nodeType, uint32(id)
}
