// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package idgen

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	redisAddr = os.Getenv("REDIS_ADDR")
	mongoUri  = os.Getenv("MONGO_URI")
)

func init() {
	println("redis addr is", redisAddr)
	println("mongo uri is", mongoUri)
}

func createCounterStorage(storeTye string, label string) CounterStore {
	switch storeTye {
	case "mongo":
		var db = "testdb"
		return NewMongoDBCounter(mongoUri, db, label)
	case "redis":
		return NewRedisCounter(redisAddr, label)
	default:
		panic(fmt.Sprintf("invalid storage type %s", storeTye))
	}
}

func putIfAbsent(guard sync.Locker, uuids map[int64]bool, id int64) bool {
	guard.Lock()
	defer guard.Unlock()

	if _, found := uuids[id]; !found {
		uuids[id] = true
		return true
	} else {
		return false
	}
}

type IDGenWorkerContext struct {
	wg           sync.WaitGroup
	guard        sync.Mutex
	uuids        map[int64]bool
	eachMaxCount int
	genMaker     func() IDGenerator
	startAt      time.Time
	stopAt       time.Time
}

func NewWorkerContext(eachMaxCount int, f func() IDGenerator) *IDGenWorkerContext {
	return &IDGenWorkerContext{
		genMaker:     f,
		eachMaxCount: eachMaxCount,
		uuids:        make(map[int64]bool, 10000),
		startAt:      time.Now(),
	}
}

func (ctx *IDGenWorkerContext) serve(t *testing.T, c context.Context, gid int) {
	defer ctx.wg.Done()
	var idGen = ctx.genMaker()
	for i := 0; i < ctx.eachMaxCount; i++ {
		id, err := idGen.Next(c)
		if err != nil {
			t.Fatalf("worker %d generate error: %v", gid, err)
		}
		// fmt.Printf("worker %d generate id %d\n", worker, id)
		if !putIfAbsent(&ctx.guard, ctx.uuids, id) {
			t.Fatalf("worker %d: tick %d, id %d is already produced by worker", gid, i, id)
		}
	}
}

func (ctx *IDGenWorkerContext) Go(t *testing.T, c context.Context, gid int) {
	ctx.wg.Add(1)
	go ctx.serve(t, c, gid)
}

func (ctx *IDGenWorkerContext) Wait() {
	ctx.wg.Wait()
	ctx.stopAt = time.Now()
}

func (ctx *IDGenWorkerContext) Duration() time.Duration {
	return ctx.stopAt.Sub(ctx.startAt)
}

func runSeqIDTestSimple(t *testing.T, ctx context.Context, storeTye, label string) {
	var store = createCounterStorage(storeTye, label)
	var seq = NewSegmentIDGen(store, DefaultSeqIDStep)
	if err := seq.Init(ctx); err != nil {
		t.Fatalf("Init: %v", err)
	}
	var m = make(map[int64]bool)
	var start = time.Now()
	const tetLoad = 2000000
	for i := 0; i < tetLoad; i++ {
		uid := seq.MustNext(ctx)
		if _, found := m[uid]; found {
			t.Fatalf("duplicate key %d exist", uid)
		}
		m[uid] = true
	}
	var elapsed = time.Since(start).Seconds()
	t.Logf("etcd QPS %.2f/s", float64(tetLoad)/elapsed)
}

// N个并发worker，共享一个生成器, 测试生成id的一致性
func runSeqIDTestConcurrent(t *testing.T, ctx context.Context, storeTye, label string) {
	var gcnt = 20
	var eachMax = 500000
	var store = createCounterStorage(storeTye, label)
	var seq = NewSegmentIDGen(store, DefaultSeqIDStep)
	if err := seq.Init(ctx); err != nil {
		t.Fatalf("Init: %v", err)
	}
	var workerCtx = NewWorkerContext(eachMax, func() IDGenerator { return seq })
	for i := 0; i < gcnt; i++ {
		workerCtx.Go(t, ctx, i)
	}
	workerCtx.Wait()

	var elapsed = workerCtx.Duration().Seconds()
	if !t.Failed() {
		t.Logf("QPS %.2f/s", float64(gcnt*eachMax)/elapsed)
	}
}

// N个并发worker，每个worker单独生成器, 测试生成id的一致性
func runSeqIDTestDistributed(t *testing.T, ctx context.Context, storeTye, label string) {
	var gcnt = 20
	var eachMax = 500000
	var generator = func() IDGenerator {
		var store = createCounterStorage(storeTye, label)
		var seq = NewSegmentIDGen(store, DefaultSeqIDStep)
		if err := seq.Init(ctx); err != nil {
			t.Fatalf("Init: %v", err)
		}
		return seq
	}

	var workerCtx = NewWorkerContext(eachMax, generator)
	for i := 0; i < gcnt; i++ {
		workerCtx.Go(t, ctx, i)
	}
	workerCtx.Wait()

	var elapsed = workerCtx.Duration().Seconds()
	if !t.Failed() {
		t.Logf("QPS %.2f/s", float64(gcnt*eachMax)/elapsed)
	}
}
