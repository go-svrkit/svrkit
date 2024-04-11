// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package idgen

import (
	"sync"
	"testing"
	"time"
)

type NoLock struct {
}

func (NoLock) Lock() {
}

func (NoLock) Unlock() {
}

func TestSnowflakeLimit(t *testing.T) {
	var sf = NewSnowflake(0)
	t.Logf("a typical uuid is %d", sf.MustNext())

	var epoch = time.Unix(CustomEpoch/int64(time.Second), 0)
	var endOfWorld = epoch.Add(time.Duration(TimeUnit) * MaxTimeUnits)
	t.Logf("the end time of uuid is %v", endOfWorld.UTC())
}

func TestSnowflake_ClockBackwards(t *testing.T) {
	var count = 0
	var old = currentTimeUnit
	currentTimeUnit = func() int64 {
		count++
		if count%10 == 0 {
			return old() - 50
		}
		return old()
	}
	var dict = make(map[int64]bool)
	var sf = NewSnowflake(1)
	for i := 0; i < 10; i++ {
		if uuid, err := sf.Next(); err != nil {
			t.Fatalf("generate uuid failed: %v", err)
		} else {
			if dict[uuid] {
				t.Fatalf("duplicate uuid %d", uuid)
			}
			dict[uuid] = true
		}
	}
}

func TestSnowflakeNext(t *testing.T) {
	const count = 1000000
	var dict = make(map[int64]bool)
	var sf = NewSnowflake(1234)
	var start = time.Now()
	var l NoLock
	for i := 0; i < count; i++ {
		id := sf.MustNext()
		if !putIfAbsent(&l, dict, id) {
			t.Fatalf("duplicate id %d after %d", id, i)
			return
		}
	}
	if len(dict) != count {
		t.Fatalf("duplicate id found")
	}
	var expired = time.Since(start)
	t.Logf("QPS: %.02f/s", float64(len(dict))/expired.Seconds())
	// Output:
	//   QPS: 288022.46/s
}

var (
	uuidMap   = make(map[int64]bool, 1000000)
	uuidGuard sync.Mutex
)

func newSnowflakeIDWorker(t *testing.T, sf *Snowflake, wg *sync.WaitGroup, gid int, count int) {
	defer wg.Done()
	//t.Logf("snowflake worker %d started", gid)
	for i := 0; i < count; i++ {
		id := sf.MustNext()
		if !putIfAbsent(&uuidGuard, uuidMap, id) {
			t.Errorf("duplicate id %d after %d", id, i)
			return
		}
	}
	//t.Logf("snowflake worker %d done", gid)
}

// 开启N个goroutine，测试UUID的并发性
func TestSnowflakeConcurrent(t *testing.T) {
	var gcount = 20
	var eachCnt = 100000
	var start = time.Now()
	var sf = NewSnowflake(1234)
	var wg sync.WaitGroup
	wg.Add(gcount)
	for i := 0; i < gcount; i++ {
		go newSnowflakeIDWorker(t, sf, &wg, i, eachCnt)
	}
	wg.Wait()
	if len(uuidMap) != gcount*eachCnt {
		t.Fatalf("duplicate id found")
	}
	var expired = time.Since(start)
	t.Logf("QPS: %.02f/s", float64(len(uuidMap))/expired.Seconds())
	// Output:
	//   QPS: 288876.86/s
}

func BenchmarkSnowflakeGen(b *testing.B) {
	var sf = NewSnowflake(1234)
	for i := 0; i < 10000; i++ {
		sf.MustNext()
	}
}
