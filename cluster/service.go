// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package cluster

import (
	"context"
	"reflect"

	"gopkg.in/svrkit.v1/qnet"
)

// IService 服务接口
type IService interface {
	ServiceType() uint16
	Node() qnet.NodeID
	SetNode(node qnet.NodeID)
	SetAttr(name string, v any)
	Init(ctx context.Context) error
	Startup(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown()
}

var (
	serviceRegistry = make(map[uint16]reflect.Type)
)

// Register 注册服务
func Register(service IService) {
	var serviceType = service.ServiceType()
	if dup := serviceRegistry[serviceType]; dup != nil {
		panic("duplicate service type registration")
	}
	serviceRegistry[serviceType] = reflect.TypeOf(service).Elem()
}

// CreateService 创建服务
func CreateService(serviceType uint16) IService {
	rtype, found := serviceRegistry[serviceType]
	if !found {
		return nil
	}
	var rval = reflect.New(rtype)
	return rval.Interface().(IService)
}
