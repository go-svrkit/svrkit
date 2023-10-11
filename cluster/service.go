package cluster

import (
	"context"
	"reflect"

	"gopkg.in/svrkit.v1/qnet"
)

// IService 服务接口
type IService interface {
	ServiceType() uint16
	ServiceID() uint32
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
	if service == nil {
		panic("invalid service")
	}
	var serviceType = service.ServiceType()
	if dup := serviceRegistry[serviceType]; dup != nil {
		panic("duplicate registration")
	}
	serviceRegistry[serviceType] = reflect.TypeOf(service).Elem()
}

// CreateService 创建服务
func CreateService(serviceType uint16) IService {
	var rtype = serviceRegistry[serviceType]
	var rval = reflect.New(rtype)
	return rval.Interface().(IService)
}
