// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package idgen

import (
	"context"
	"fmt"
	"math"
	"sync"
)

const (
	DefaultSeqIDStep = 2000 // 默认步长
)

// 分段发号器算法的ID生成器
//
// 1. 算法把一个64位的整数按step范围划分为N个号段；
// 2. 存储组件维护一个持续递增的计数器，表示当前未分配的号段；
// 3. service从存储组件拿到号段（计数器）后才可分配此号段内的ID;
// 4. 在存储组件设置key后勿删除！
//
// 号段ID只保证了唯一性，无法保证顺序性（递增）
// 因为多个服务同时生成， 如果服务1的生成速度如果比服务2快，服务1的ID号段会先用完，
// 那么服务1上按时钟先分配的ID会大于服务2上按时钟后分配的ID；

type SegmentIDGen struct {
	guard   sync.Mutex   // 线程安全
	store   CounterStore // 保存号段
	step    int32        // 当前区间
	counter int64        // 当前号段
	maxID   int64        // 当前分段的最大值
	lastID  int64        // 上次生成的ID
}

func NewSegmentIDGen(store CounterStore, step int32) *SegmentIDGen {
	if step <= 0 {
		step = DefaultSeqIDStep
	}
	return &SegmentIDGen{
		store: store,
		step:  step,
	}
}

func (g *SegmentIDGen) Init(ctx context.Context) error {
	return g.reload(ctx)
}

// Next 分配下一个ID
func (g *SegmentIDGen) Next(ctx context.Context) (int64, error) {
	g.guard.Lock()
	defer g.guard.Unlock()

	var nextId = g.lastID + 1
	// 在当前号段内，直接分配
	if nextId <= g.maxID {
		g.lastID = nextId
		return nextId, nil
	}
	// 需要重新申请号段
	if err := g.reload(ctx); err != nil {
		return 0, err
	}
	nextId = g.lastID + 1
	g.lastID = nextId
	return nextId, nil
}

func (g *SegmentIDGen) MustNext(ctx context.Context) int64 {
	id, err := g.Next(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

func (g *SegmentIDGen) reload(ctx context.Context) error {
	counter, err := g.store.Incr(ctx)
	if err != nil {
		return err
	}
	// 检查是否超过int64数值范围
	var counterMax = math.MaxInt64 / int64(g.step)
	if counter >= counterMax {
		return fmt.Errorf("counter overflow %d", counter)
	}
	g.counter = counter
	g.lastID = counter * int64(g.step)
	g.maxID = (g.counter + 1) * int64(g.step)
	return nil
}
