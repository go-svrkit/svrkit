// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package zset

import (
	"cmp"
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/collections/util"
)

func uniqueScoreGen() func() int64 {
	var dict = map[int64]bool{}
	return func() int64 {
		for {
			var score = rand.Int63() % 1000000
			if _, found := dict[score]; !found {
				return score
			}
		}
	}
}

func Test_zslRandLevel(t *testing.T) {
	for i := 0; i < 100; i++ {
		var level = zslRandLevel()
		assert.GreaterOrEqual(t, level, 1)
		assert.LessOrEqual(t, level, ZSKIPLIST_MAXLEVEL)
	}
}

func testSkipList(t *testing.T, sl *ZSkipList[string]) {
	var pairs = make([]util.Pair[int64, string], 0, sl.Len())
	sl.Range(func(score int64, ele string) {
		pairs = append(pairs, util.MakePair(score, ele))
	})
	assert.Equal(t, len(pairs), sl.Len())
	var isSorted = slices.IsSortedFunc(pairs, func(a, b util.Pair[int64, string]) int {
		return cmp.Compare(a.First, b.First)
	})
	assert.True(t, isSorted)

	var head = sl.HeadNode()
	var tail = sl.TailNode()
	assert.Equal(t, pairs[0].Second, head.Ele)
	assert.Equal(t, pairs[len(pairs)-1].Second, tail.Ele)
}

func TestZSkipList_Insert(t *testing.T) {
	// 顺序插入
	var sl = NewZSkipList[string](util.OrderedCmp[string])
	for i := 100; i > 0; i-- {
		var score = 100 + int64(i)
		sl.Insert(score, fmt.Sprintf("test%d", i))
	}
	testSkipList(t, sl)

	//t.Logf("%s", sl.String())
	// 随机插入
	sl = NewZSkipList[string](util.OrderedCmp[string])
	var g = uniqueScoreGen()
	for i := 0; i < 100; i++ {
		var score = g()
		sl.Insert(score, fmt.Sprintf("test%d", i))
	}
	testSkipList(t, sl)
}

func TestZSkipList_Delete(t *testing.T) {
	var sl = NewZSkipList[string](util.OrderedCmp[string])
	sl.Insert(10, "test01")
	sl.Insert(20, "test02")
	sl.Insert(30, "test03")
	assert.Equal(t, 3, sl.Len())
	assert.Equal(t, "test01=10,test02=20,test03=30,", sl.String())

	// 删除1个已有的节点
	var node = sl.Delete(10, "test01")
	assert.NotNil(t, node)
	assert.Equal(t, 2, sl.Len())
	assert.Equal(t, "test01", node.Ele)
	assert.Equal(t, int64(10), node.Score)

	// 删除不存在的节点
	node = sl.Delete(40, "test04")
	assert.Nil(t, node)
	assert.Equal(t, 2, sl.Len())
}

func TestZSkipList_UpdateScore(t *testing.T) {
	var sl = NewZSkipList[string](util.OrderedCmp[string])
	sl.Insert(10, "test01")
	sl.Insert(20, "test02")
	sl.Insert(30, "test03")
	sl.Insert(40, "test04")
	assert.Equal(t, 4, sl.Len())
	assert.Equal(t, "test01=10,test02=20,test03=30,test04=40,", sl.String())

	var node = sl.UpdateScore("test02", 20, 40)
	assert.NotNil(t, node)
	assert.Equal(t, int64(40), node.Score)
	assert.Equal(t, "test01=10,test03=30,test02=40,test04=40,", sl.String())
}

func TestZSkipList_DeleteRangeByRank(t *testing.T) {
	var dict = make(map[string]int64)
	var sl = NewZSkipList[string](util.OrderedCmp[string])
	for i := 1; i <= 10; i++ {
		var score = 10 * int64(i)
		var elem = fmt.Sprintf("test%02d", i)
		sl.Insert(score, elem)
		dict[elem] = score
	}
	assert.Equal(t, 10, sl.Len())

	sl.DeleteRangeByRank(1, 5, dict)
	assert.Equal(t, 5, sl.Len())
	var s = "test06=60,test07=70,test08=80,test09=90,test10=100,"
	assert.Equal(t, s, sl.String())

	sl.DeleteRangeByRank(4, 5, dict)
	assert.Equal(t, 3, sl.Len())
	assert.Equal(t, "test06=60,test07=70,test08=80,", sl.String())

	sl.DeleteRangeByRank(1, 1, dict)
	assert.Equal(t, 2, sl.Len())
	assert.Equal(t, "test07=70,test08=80,", sl.String())

	sl.DeleteRangeByRank(1, 10, dict)
	assert.Equal(t, 0, sl.Len())
}

func TestZSkipList_DeleteRangeByScore(t *testing.T) {
	var dict = make(map[string]int64)
	var sl = NewZSkipList[string](util.OrderedCmp[string])
	for i := 1; i <= 10; i++ {
		var score = 10 * int64(i)
		var elem = fmt.Sprintf("test%02d", i)
		sl.Insert(score, elem)
		dict[elem] = score
	}
	assert.Equal(t, 10, sl.Len())

	sl.DeleteRangeByScore(20, 90, dict)
	assert.Equal(t, 2, sl.Len())
	assert.Equal(t, "test01=10,test10=100,", sl.String())
}

func TestZSkipList_GetRank(t *testing.T) {
	var sl = NewZSkipList[string](util.OrderedCmp[string])
	for i := 1; i <= 100; i++ {
		var score = 10 * int64(i)
		var elem = fmt.Sprintf("test%02d", i)
		sl.Insert(score, elem)
	}
	assert.Equal(t, 100, sl.Len())

	for i := 1; i <= 100; i++ {
		var score = 10 * int64(i)
		var elem = fmt.Sprintf("test%02d", i)
		var rank = sl.GetRank(score, elem)
		assert.Equal(t, i, rank)
	}
}

func TestZSkipList_GetElementByRank(t *testing.T) {
	var sl = NewZSkipList[string](util.OrderedCmp[string])
	for i := 1; i <= 100; i++ {
		var score = 10 * int64(i)
		var elem = fmt.Sprintf("test%02d", i)
		sl.Insert(score, elem)
	}
	assert.Equal(t, 100, sl.Len())

	for i := 1; i <= 100; i++ {
		var score = 10 * int64(i)
		var elem = fmt.Sprintf("test%02d", i)
		var node = sl.GetElementByRank(i)
		assert.NotNil(t, node)
		assert.Equal(t, node.Score, score)
		assert.Equal(t, node.Ele, elem)
	}
}
