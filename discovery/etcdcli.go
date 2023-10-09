// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	clientv3 "go.etcd.io/etcd/client/v3"

	"gopkg.in/svrkit.v1/logger"
)

var (
	ErrEmptyLease    = errors.New("empty lease")
	ErrNodeKeyExist  = errors.New("node key exist")
	ErrNoKeyDeleted  = errors.New("no key deleted")
	ErrEmptyEndpoint = errors.New("empty endpoint")
)

const (
	EventChanCapacity = 1024
	OpTimeout         = 5

	VerboseLv1 = 1
	VerboseLv2 = 2
)

// Client 基于etcd的服务发现
type Client struct {
	closing   atomic.Int32     //
	verbose   int32            //
	endpoints []string         // etcd server address
	namespace string           // name space of key
	username  string           //
	passwd    string           //
	client    *clientv3.Client // etcd client
}

func NewClient(endpoints, namespace, username, passwd string) *Client {
	d := &Client{
		endpoints: strings.Split(endpoints, ","),
		namespace: namespace,
		username:  username,
		passwd:    passwd,
		verbose:   VerboseLv1,
	}
	return d
}

func (c *Client) Init(parentCtx context.Context) error {
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

func (c *Client) Close() {
	if !c.closing.CompareAndSwap(0, 1) {
		return
	}
	if c.client != nil {
		c.client.Close()
		c.client = nil
	}
}

func (c *Client) SetVerbose(v int32) {
	c.verbose = v
}

func (c *Client) FormatKey(name string) string {
	if name[0] == '/' {
		return c.namespace + name
	}
	return c.namespace + "/" + name
}

// IsNodeExist 节点是否存在
func (c *Client) IsNodeExist(ctx context.Context, name string) (bool, error) {
	var key = c.FormatKey(name)
	resp, err := c.client.Get(ctx, key, clientv3.WithCountOnly())
	if err != nil {
		return false, err
	}
	return resp.Count > 0, nil
}

// GetKeyValue 获取key的值
func (c *Client) GetKeyValue(ctx context.Context, name string) ([]byte, error) {
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
func (c *Client) GetNode(ctx context.Context, name string) (*Node, error) {
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
func (c *Client) PutNode(ctx context.Context, name string, value any, leaseId int64) error {
	var key = c.FormatKey(name)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var resp *clientv3.PutResponse
	if leaseId <= 0 {
		resp, err = c.client.Put(ctx, key, bytesAsString(data))
	} else {
		resp, err = c.client.Put(ctx, key, bytesAsString(data), clientv3.WithLease(clientv3.LeaseID(leaseId)))
	}
	if err != nil {
		return err
	}
	if c.verbose >= VerboseLv1 {
		logger.Infof("put key [%s] at rev %d", key, resp.Header.Revision)
	}
	return nil
}

// DelKey 删除一个key
func (c *Client) DelKey(ctx context.Context, name string) error {
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
func (c *Client) ListNodes(ctx context.Context, prefix string) ([]Node, error) {
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
func (c *Client) GrantLease(ctx context.Context, ttl int) (int64, error) {
	lease, err := c.client.Grant(ctx, int64(ttl))
	if err != nil {
		return 0, err
	}
	if lease == nil {
		return 0, ErrEmptyLease
	}
	return int64(lease.ID), nil
}

func (c *Client) GetLeaseTTL(ctx context.Context, leaseId int64) (int, error) {
	resp, err := c.client.TimeToLive(ctx, clientv3.LeaseID(leaseId))
	if err != nil {
		return 0, nil
	}
	return int(resp.TTL), nil
}

// RevokeLease 撤销一个lease
func (c *Client) RevokeLease(ctx context.Context, leaseId int64) error {
	_, err := c.client.Revoke(ctx, clientv3.LeaseID(leaseId))
	return err
}

// 用于注册并保活节点
type nodeKeepAliveContext struct {
	stopChan   chan struct{}
	LeaseId    int64
	LeaseAlive bool
	Name       string
	Value      any
	TTL        int
}

func NewNodeKeepAliveContext(name string, value any, ttl int) *nodeKeepAliveContext {
	return &nodeKeepAliveContext{
		stopChan: make(chan struct{}, 1),
		Name:     name,
		Value:    value,
		TTL:      ttl,
	}
}

func (c *Client) RevokeKeepAlive(ctx context.Context, regCtx *nodeKeepAliveContext) error {
	if c.verbose >= VerboseLv1 {
		logger.Infof("try revoke node %s lease %d", regCtx.Name, regCtx.LeaseId)
	}
	if regCtx.LeaseId == 0 || !regCtx.LeaseAlive {
		if c.verbose >= VerboseLv1 {
			logger.Infof("node %s lease %d is not alive", regCtx.Name, regCtx.LeaseId)
		}
		return nil
	}
	if err := c.RevokeLease(ctx, regCtx.LeaseId); err != nil {
		logger.Warnf("revoke node %s lease %x failed: %v", regCtx.Name, regCtx.LeaseId, err)
		return err
	} else {
		if c.verbose >= VerboseLv1 {
			logger.Infof("revoke node %s lease %x done", regCtx.Name, regCtx.LeaseId)
		}
	}
	return nil
}

// RegisterNode 注册一个节点信息，并返回一个ttl秒的lease
func (c *Client) RegisterNode(rootCtx context.Context, name string, value any, ttl int) (int64, error) {
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

func revokeLeaseWithTimeout(c *Client, leaseId int64) {
	if c.verbose >= VerboseLv1 {
		logger.Infof("try revoke lease %d", leaseId)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*OpTimeout)
	defer cancel()
	if err := c.RevokeLease(ctx, leaseId); err != nil {
		logger.Warnf("revoke lease %x failed: %v", leaseId, err)
	} else {
		logger.Infof("revoke lease %x done", leaseId)
	}
}

func (c *Client) aliveKeeper(ctx context.Context, kaChan <-chan *clientv3.LeaseKeepAliveResponse, stopChan chan struct{}, leaseId int64) {
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
				logger.Infof("lease %d is not alive", leaseId)
				return
			}
			if c.verbose >= VerboseLv2 {
				logger.Infof("lease %d respond alive, ttl %d", ka.ID, ka.TTL)
			}

		case <-ctx.Done():
			logger.Infof("stop keepalive with lease %d", leaseId)
			return
		}
	}
}

// KeepAlive lease保活，当lease撤销时此stopChan被激活
func (c *Client) KeepAlive(ctx context.Context, stopChan chan struct{}, leaseId int64) error {
	kaChan, err := c.client.KeepAlive(ctx, clientv3.LeaseID(leaseId))
	if err != nil {
		return nil
	}
	go c.aliveKeeper(ctx, kaChan, stopChan, leaseId)
	return nil
}

func (c *Client) doRegisterNode(ctx context.Context, regCtx *nodeKeepAliveContext) error {
	var err error
	if c.verbose >= VerboseLv1 {
		logger.Infof("try register key: %s", c.FormatKey(regCtx.Name))
	}
	regCtx.LeaseAlive = false
	regCtx.LeaseId = 0

	regCtx.LeaseId, err = c.RegisterNode(ctx, regCtx.Name, regCtx.Value, regCtx.TTL)
	if err != nil {
		return err
	}
	if err = c.KeepAlive(ctx, regCtx.stopChan, regCtx.LeaseId); err != nil {
		return err
	}
	regCtx.LeaseAlive = true
	if c.verbose >= VerboseLv1 {
		logger.Infof("register key [%s] with lease %x done", c.FormatKey(regCtx.Name), regCtx.LeaseId)
	}
	return nil
}

func (c *Client) regAliveKeeper(ctx context.Context, regCtx *nodeKeepAliveContext) {
	var ticker = time.NewTicker(time.Second) // 1s
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if !regCtx.LeaseAlive {
				if err := c.doRegisterNode(ctx, regCtx); err != nil {
					logger.Infof("register or keepalive %s failed: %v", regCtx.Name, err)
				}
			}

		case <-regCtx.stopChan:
			var leaseId = regCtx.LeaseId
			regCtx.LeaseAlive = false
			regCtx.LeaseId = 0
			if c.verbose >= VerboseLv1 {
				logger.Infof("node %s lease(%d) is not alive, try register later", regCtx.Name, leaseId)
			}

		case <-ctx.Done():
			if c.verbose >= VerboseLv1 {
				logger.Infof("register alive keeper with key %s stopped", regCtx.Name)
			}
			return
		}
	}
}

// RegisterAndKeepAliveForever 注册一个节点，并永久保活
func (c *Client) RegisterAndKeepAliveForever(ctx context.Context, name string, value any, ttl int) (*nodeKeepAliveContext, error) {
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
			logger.Errorf("unmarshal node %s: %v", event.Key, err)
			return
		}
	}

	select {
	case eventChan <- event:
	default:
		logger.Warnf("watch event channel is full, new event lost: %v", event)
	}
}

// WatchDir 订阅目录下的节点变化
func (c *Client) WatchDir(ctx context.Context, dir string) <-chan *NodeEvent {
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
					logger.Warnf("watch key %s canceled: %v", key, resp.Err())
					return
				}
				for _, ev := range resp.Events {
					propagateWatchEvent(eventChan, ev)
				}

			case <-ctx.Done():
				if c.client != nil {
					if err := c.client.Watcher.Close(); err != nil {
						logger.Warnf("close watcher: %v", err)
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
func (c *Client) WatchDirTo(ctx context.Context, dir string, nodeMap *NodeMap) {
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
		nodeMap.InsertNode(ev.Node)
	case EventUpdate:
		nodeMap.InsertNode(ev.Node) // 插入前会先检查是否有重复
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

func bytesAsString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
