// This skiplist implementation is almost a translation of the original
// algorithm described by William Pugh in "Skip Lists: A Probabilistic
// Alternative to Balanced Trees", modified in three ways:
// a) this implementation allows for repeated scores.
// b) the comparison is not just by Key (our 'score') but by satellite data.
// c) there is a back pointer, so it's a doubly linked list with the back
// pointers being only at "level 1". This allows to traverse the list
// from tail to head.
//
// https://github.com/redis/redis/blob/6.2.14/src/t_zset.c
// https://en.wikipedia.org/wiki/Skip_list

package zset

import (
	"cmp"
	"fmt"
	"io"
	"math"
	"math/rand"
	"strings"
)

const (
	ZSKIPLIST_MAXLEVEL = 32   // Should be enough for 2^64 elements
	ZSKIPLIST_P        = 0.25 // Skiplist P = 1/4
)

// each level of list node
type zskipListLevel[T comparable] struct {
	forward *ZSkipListNode[T] // link to next node
	span    int               // node # between this and forward link
}

// ZSkipListNode
type ZSkipListNode[T comparable] struct {
	Ele      T
	Score    int64
	backward *ZSkipListNode[T]
	level    []zskipListLevel[T]
}

func newZSkipListNode[T comparable](level int, score int64, elem T) *ZSkipListNode[T] {
	return &ZSkipListNode[T]{
		Ele:   elem,
		Score: score,
		level: make([]zskipListLevel[T], level),
	}
}

func (n *ZSkipListNode[T]) Before() *ZSkipListNode[T] {
	return n.backward
}

// Next return next forward pointer
func (n *ZSkipListNode[T]) Next() *ZSkipListNode[T] {
	return n.level[0].forward
}

type Comparator[T any] func(a, b T) int

// ZSkipList 带索引的排序链表
type ZSkipList[T comparable] struct {
	head       *ZSkipListNode[T] // 头结点
	tail       *ZSkipListNode[T] // 尾节点（最大值节点）
	comparator Comparator[T]     //
	length     int               // 节点数
	level      int               // 层级
}

func NewZSkipList[T comparable](comparator Comparator[T]) *ZSkipList[T] {
	var zero T
	return &ZSkipList[T]{
		level:      1,
		comparator: comparator,
		head:       newZSkipListNode[T](ZSKIPLIST_MAXLEVEL, 0, zero),
	}
}

func NewZSkipListCmp[T cmp.Ordered]() *ZSkipList[T] {
	return NewZSkipList[T](cmp.Compare[T])
}

// 返回新节点的随机层级[1-ZSKIPLIST_MAXLEVEL]
func zslRandLevel() int {
	var level = 1
	for (rand.Int31() & 0xFFFF) < int32(math.Floor(ZSKIPLIST_P*float64(0xFFFF))) {
		level++
	}
	return min(level, ZSKIPLIST_MAXLEVEL)
}

func (zsl *ZSkipList[T]) Len() int {
	return zsl.length
}

func (zsl *ZSkipList[T]) Height() int {
	return zsl.level
}

func (zsl *ZSkipList[T]) HeadNode() *ZSkipListNode[T] {
	return zsl.head.level[0].forward
}

func (zsl *ZSkipList[T]) TailNode() *ZSkipListNode[T] {
	return zsl.tail
}

func (zsl *ZSkipList[T]) Range(action func(score int64, elem T)) {
	var node = zsl.head.level[0].forward
	for node != nil {
		action(node.Score, node.Ele)
		if len(node.level) > 0 {
			node = node.level[0].forward
		} else {
			node = nil
		}
	}
}

// Insert 插入一个不存在的节点
func (zsl *ZSkipList[T]) Insert(score int64, ele T) *ZSkipListNode[T] {
	var update [ZSKIPLIST_MAXLEVEL]*ZSkipListNode[T]
	var rank [ZSKIPLIST_MAXLEVEL]int

	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		// store rank that is crossed to reach the insert position
		if i == zsl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for x.level[i].forward != nil &&
			(x.level[i].forward.Score < score ||
				(x.level[i].forward.Score == score &&
					zsl.comparator(x.level[i].forward.Ele, ele) < 0)) {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}
	// we assume the element is not already inside, since we allow duplicated
	// scores, reinserting the same element should never happen since the
	// caller of zslInsert() should test in the hash table if the element is
	// already inside or not.
	var level = zslRandLevel()
	if level > zsl.level {
		for i := zsl.level; i < level; i++ {
			rank[i] = 0
			update[i] = zsl.head
			update[i].level[i].span = zsl.length
		}
		zsl.level = level
	}
	x = newZSkipListNode(level, score, ele)
	for i := 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		// update span covered by update[i] as x is inserted here
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}
	// increment span for untouched levels
	for i := level; i < zsl.level; i++ {
		update[i].level[i].span++
	}
	if update[0] == zsl.head {
		x.backward = nil
	} else {
		x.backward = update[0]
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		zsl.tail = x
	}
	zsl.length++
	return x
}

// 删除一个节点
func (zsl *ZSkipList[T]) deleteNode(x *ZSkipListNode[T], update []*ZSkipListNode[T]) {
	for i := 0; i < zsl.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span -= 1
		}
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		zsl.tail = x.backward
	}
	for zsl.level > 1 && zsl.head.level[zsl.level-1].forward == nil {
		zsl.level--
	}
	zsl.length--
}

// Delete 删除对应score的节点
func (zsl *ZSkipList[T]) Delete(score int64, ele T) *ZSkipListNode[T] {
	var update [ZSKIPLIST_MAXLEVEL]*ZSkipListNode[T]
	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.Score < score ||
				(x.level[i].forward.Score == score &&
					zsl.comparator(x.level[i].forward.Ele, ele) < 0)) {
			x = x.level[i].forward
		}
		update[i] = x
	}

	// We may have multiple elements with the same score, what we need
	// is to find the element with both the right score and object.
	x = x.level[0].forward
	if x != nil {
		if score == x.Score && zsl.comparator(x.Ele, ele) == 0 {
			zsl.deleteNode(x, update[:])
			return x
		}
	}
	return nil // not found
}

// UpdateScore 更新分数
func (zsl *ZSkipList[T]) UpdateScore(ele T, curScore, newScore int64) *ZSkipListNode[T] {
	var update [ZSKIPLIST_MAXLEVEL]*ZSkipListNode[T]

	// We need to seek to element to update to start: this is useful anyway,
	// we'll have to update or remove it.
	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.Score < curScore ||
				(x.level[i].forward.Score == curScore &&
					zsl.comparator(x.level[i].forward.Ele, ele) < 0)) {
			x = x.level[i].forward
		}
		update[i] = x
	}
	// Jump to our element: note that this function assumes that the
	// element with the matching score exists.
	x = x.level[0].forward

	// If the node, after the score update, would be still exactly
	// at the same position, we can just update the score without
	// actually removing and re-inserting the element in the skiplist.
	if (x.backward == nil || x.backward.Score < newScore) &&
		(x.level[0].forward == nil || x.level[0].forward.Score > newScore) {
		x.Score = newScore
		return x
	}
	// No way to reuse the old node: we need to remove and insert a new
	// one at a different place.
	zsl.deleteNode(x, update[:])
	var newNode = zsl.Insert(newScore, x.Ele)
	var zero T
	x.Ele = zero
	return newNode
}

// DeleteRangeByRank delete nodes with rank [rank >= start && rank <= end]
func (zsl *ZSkipList[T]) DeleteRangeByRank(start, end int, dict map[T]int64) int {
	var update [ZSKIPLIST_MAXLEVEL]*ZSkipListNode[T]
	var traversed, removed int
	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (traversed+x.level[i].span < start) {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}
	traversed++
	x = x.level[0].forward
	for x != nil && traversed <= end {
		var next = x.level[0].forward
		zsl.deleteNode(x, update[:])
		delete(dict, x.Ele)
		removed++
		traversed++
		x = next
	}
	return removed
}

// DeleteRangeByScore delete nodes with [score >= min && score <= max]
func (zsl *ZSkipList[T]) DeleteRangeByScore(min, max int64, dict map[T]int64) int {
	var update [ZSKIPLIST_MAXLEVEL]*ZSkipListNode[T]
	var removed int
	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.Score < min {
			x = x.level[i].forward
		}
		update[i] = x
	}

	// Current node is the last with score < or <= min
	x = x.level[0].forward

	// Delete nodes while in range
	for x != nil && x.Score <= max {
		var next = x.level[0].forward
		zsl.deleteNode(x, update[0:])
		delete(dict, x.Ele)
		removed++
		x = next
	}
	return removed
}

// GetRank 获取score所在的排名，排名从1开始
func (zsl *ZSkipList[T]) GetRank(score int64, ele T) int {
	var rank = 0
	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.Score < score ||
				(x.level[i].forward.Score == score &&
					zsl.comparator(x.level[i].forward.Ele, ele) <= 0)) {
			rank += x.level[i].span
			x = x.level[i].forward
		}

		// x might be equal to zsl->header, so test if obj is non-nil
		if zsl.comparator(x.Ele, ele) == 0 {
			return rank
		}
	}
	return 0
}

// GetElementByRank 根据排名获得节点，排名从1开始
func (zsl *ZSkipList[T]) GetElementByRank(rank int) *ZSkipListNode[T] {
	var tranversed = 0
	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (tranversed+x.level[i].span <= rank) {
			tranversed += x.level[i].span
			x = x.level[i].forward
		}
		if tranversed == rank {
			return x
		}
	}
	return nil
}

// IsInRange Returns if there is a part of the zset is in range.
func (zsl *ZSkipList[T]) IsInRange(min, max int64) bool {
	if min > max {
		return false
	}
	var x = zsl.tail
	if x == nil || x.Score < min {
		return false
	}
	x = zsl.head.level[0].forward
	if x == nil || x.Score > max {
		return false
	}
	return true
}

// FirstInRange find the first node that is contained in the specified range.
// Returns NULL when no element is contained in the range.
func (zsl *ZSkipList[T]) FirstInRange(min, max int64) *ZSkipListNode[T] {
	if !zsl.IsInRange(min, max) {
		return nil
	}
	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		// Go forward while *OUT* of range.
		for x.level[i].forward != nil && x.level[i].forward.Score < min {
			x = x.level[i].forward
		}
	}
	// This is an inner range, so the next node cannot be NULL.
	x = x.level[0].forward
	if x != nil && x.Score > max {
		return nil
	}
	return x
}

// LastInRange find the last node that is contained in the specified range.
// Returns NULL when no element is contained in the range.
func (zsl *ZSkipList[T]) LastInRange(min, max int64) *ZSkipListNode[T] {
	if !zsl.IsInRange(min, max) {
		return nil
	}
	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		// Go forward while *OUT* of range.
		for x.level[i].forward != nil && x.level[i].forward.Score <= max {
			x = x.level[i].forward
		}
	}
	// Check if score >= min.
	if x.Score < min {
		return nil
	}
	return x
}

func (zsl *ZSkipList[T]) Clear() {
	var zero T
	zsl.level = 1
	zsl.head = newZSkipListNode[T](ZSKIPLIST_MAXLEVEL, 0, zero)
	zsl.tail = nil
}

func (zsl *ZSkipList[T]) ToMap() map[T]int64 {
	var dict = make(map[T]int64, zsl.Len())
	zsl.Range(func(score int64, elem T) {
		dict[elem] = score
	})
	return dict
}

func (zsl *ZSkipList[T]) Dump(w io.Writer) {
	zsl.Range(func(score int64, elem T) {
		fmt.Fprintf(w, "%v=%d,", elem, score)
	})
}

func (zsl *ZSkipList[T]) String() string {
	var sb strings.Builder
	zsl.Dump(&sb)
	return sb.String()
}
