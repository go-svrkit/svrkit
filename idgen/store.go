// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package idgen

import (
	"context"
	"errors"
)

type StorageType int8

const (
	StorageMongoDB StorageType = 1
	StorageRedis   StorageType = 2
)

var ErrCounterOutOfRange = errors.New("counter out of range")

// CounterStore 表示一个存储组件，维持一个持续递增（不一定连续）的counter
type CounterStore interface {
	Init(context.Context) error
	Close() error
	Incr(context.Context) (int64, error)
}

// IDGenerator ID生成器
type IDGenerator interface {
	Next(context.Context) (int64, error)
}
