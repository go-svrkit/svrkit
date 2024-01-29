// This skiplist implementation is almost a translation of the original
// algorithm described by William Pugh in "Skip Lists: A Probabilistic
// Alternative to Balanced Trees", modified in three ways:
// a) this implementation allows for repeated scores.
// b) the comparison is not just by key (our 'score') but by satellite data.
// c) there is a back pointer, so it's a doubly linked list with the back
// pointers being only at "level 1". This allows to traverse the list
// from tail to head.
//
// https://github.com/redis/redis/blob/6.2/src/t_zset.c
// https://en.wikipedia.org/wiki/Skip_list

package zset

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"

	"gopkg.in/svrkit.v1/collections/util"
)

const (
	ZSKIPLIST_MAXLEVEL = 32   // Should be enough
	ZSKIPLIST_P        = 0.25 // Skiplist P = 1/4
)

// each level of list node
type zskipListLevel[K comparable] struct {
	forward *ZSkipListNode[K] // link to next node
	span    int               // node # between this and forward link
}

// ZSkipListNode
type ZSkipListNode[K comparable] struct {
	Ele      K
	Score    int64
	backward *ZSkipListNode[K]
	level    []zskipListLevel[K]
}

func newZSkipListNode[K comparable](level int, score int64, element K) *ZSkipListNode[K] {
	return &ZSkipListNode[K]{
		Ele:   element,
		Score: score,
		level: make([]zskipListLevel[K], level),
	}
}

func (n *ZSkipListNode[K]) Before() *ZSkipListNode[K] {
	return n.backward
}

// Next return next forward pointer
func (n *ZSkipListNode[K]) Next() *ZSkipListNode[K] {
	return n.level[0].forward
}

// ZSkipList 带索引的排序链表
type ZSkipList[K comparable] struct {
	head       *ZSkipListNode[K]  // 头结点
	tail       *ZSkipListNode[K]  // 尾节点（最大值节点）
	comparator util.Comparator[K] //
	length     int                // 节点数
	level      int                // 层级
}

func NewZSkipList[K comparable](comparator util.Comparator[K]) *ZSkipList[K] {
	return &ZSkipList[K]{
		level:      1,
		comparator: comparator,
		head:       newZSkipListNode[K](ZSKIPLIST_MAXLEVEL, 0, nil),
	}
}

// 返回新节点的随机层级[1-ZSKIPLIST_MAXLEVEL]
func zslRandLevel() int {
	var level = 1
	for {
		var seed = rand.Uint32() & 0xFFFF
		if float32(seed) < ZSKIPLIST_P*0xFFFF {
			level++
		} else {
			break
		}
	}
	if level > ZSKIPLIST_MAXLEVEL {
		level = ZSKIPLIST_MAXLEVEL
	}
	return level
}

// Len 链表的节点数量
func (zsl *ZSkipList[K]) Len() int {
	return zsl.length
}

// Height 链表的层级
func (zsl *ZSkipList[K]) Height() int {
	return zsl.level
}

// HeadNode 头结点
func (zsl *ZSkipList[K]) HeadNode() *ZSkipListNode[K] {
	return zsl.head.level[0].forward
}

// TailNode 尾节点
func (zsl *ZSkipList[K]) TailNode() *ZSkipListNode[K] {
	return zsl.tail
}

// Insert 插入一个不存在的节点
func (zsl *ZSkipList[K]) Insert(score int64, ele K) *ZSkipListNode[K] {
	var update [ZSKIPLIST_MAXLEVEL]*ZSkipListNode[K]
	var rank [ZSKIPLIST_MAXLEVEL]int

	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		// store rank that is crossed to reach the insert position
		if i != zsl.level-1 {
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
	if update[0] != zsl.head {
		x.backward = update[0]
	} else {
		x.backward = nil
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
func (zsl *ZSkipList[K]) deleteNode(x *ZSkipListNode[K], update []*ZSkipListNode[K]) {
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
func (zsl *ZSkipList[K]) Delete(score int64, ele K) *ZSkipListNode[K] {
	var update [ZSKIPLIST_MAXLEVEL]*ZSkipListNode[K]
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
		// log.Printf("zskiplist need delete %v, but found %v\n", ele, x.Ele)
	}
	return nil // not found
}

// UpdateScore 更新分数
func (zsl *ZSkipList[K]) UpdateScore(ele K, curScore, newScore int64) *ZSkipListNode[K] {
	var update [ZSKIPLIST_MAXLEVEL]*ZSkipListNode[K]
	var x = zsl.head
	// We need to seek to element to update to start: this is useful anyway,
	// we'll have to update or remove it.
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
	x.Ele = util.ZeroOf[K]() // free the node now since zsl.Insert created a new one.
	return newNode
}

// DeleteRangeByRank 删除排名在[start-end]之间的节点，排名从1开始
func (zsl *ZSkipList[K]) DeleteRangeByRank(start, end int, dict map[K]int64) int {
	var update [ZSKIPLIST_MAXLEVEL]*ZSkipListNode[K]
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

// DeleteRangeByScore 删除score在[min-max]之间的节点
func (zsl *ZSkipList[K]) DeleteRangeByScore(min, max int64, dict map[K]int64) int {
	var update [ZSKIPLIST_MAXLEVEL]*ZSkipListNode[K]
	var removed int
	var x = zsl.head
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.Score <= min {
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
func (zsl *ZSkipList[K]) GetRank(score int64, ele K) int {
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
		if x.Ele != nil && zsl.comparator(x.Ele, ele) == 0 {
			return rank
		}
	}
	return 0
}

// GetElementByRank 根据排名获得节点，排名从1开始
func (zsl *ZSkipList[K]) GetElementByRank(rank int) *ZSkipListNode[K] {
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
func (zsl *ZSkipList[K]) IsInRange(min, max int64) bool {
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

// FirstInRange
// Find the first node that is contained in the specified range.
// Returns NULL when no element is contained in the range.
func (zsl *ZSkipList[K]) FirstInRange(min, max int64) *ZSkipListNode[K] {
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

// LastInRange
// Find the last node that is contained in the specified range.
// Returns NULL when no element is contained in the range.
func (zsl *ZSkipList[K]) LastInRange(min, max int64) *ZSkipListNode[K] {
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

func (zsl *ZSkipList[K]) String() string {
	var buf bytes.Buffer
	zsl.Dump(&buf)
	return buf.String()
}

// Dump whole list to w, mostly for debugging
func (zsl *ZSkipList[K]) Dump(w io.Writer) {
	var x = zsl.head
	// dump header
	var line bytes.Buffer
	n, _ := fmt.Fprintf(w, "<             head> ")
	prePadding(&line, n)
	for i := 0; i < zsl.level; i++ {
		if i < len(x.level) {
			if x.level[i].forward != nil {
				fmt.Fprintf(w, "[%2d] ", x.level[i].span)
				line.WriteString("  |  ")
			}
		}
	}
	fmt.Fprint(w, "\n")
	line.WriteByte('\n')
	line.WriteTo(w)

	// dump list
	var count = 0
	x = x.level[0].forward
	for x != nil {
		count++
		zsl.dumpNode(w, x, count)
		if len(x.level) > 0 {
			x = x.level[0].forward
		}
	}

	// dump tail end
	fmt.Fprintf(w, "<             end> ")
	for i := 0; i < zsl.level; i++ {
		fmt.Fprintf(w, "  _  ")
	}
	fmt.Fprintf(w, "\n")
}

func (zsl *ZSkipList[K]) dumpNode(w io.Writer, node *ZSkipListNode[K], count int) {
	var line bytes.Buffer
	var ss = fmt.Sprintf("%v", node.Ele)
	n, _ := fmt.Fprintf(w, "<%6d %4d, %s> ", node.Score, count, ss)
	prePadding(&line, n)
	for i := 0; i < zsl.level; i++ {
		if i < len(node.level) {
			fmt.Fprintf(w, "[%2d] ", node.level[i].span)
			line.WriteString("  |  ")
		} else {
			if shouldLinkVertical(zsl.head, node, i) {
				fmt.Fprintf(w, "  |  ")
				line.WriteString("  |  ")
			}
		}
	}
	fmt.Fprint(w, "\n")
	line.WriteByte('\n')
	line.WriteTo(w)
}

func shouldLinkVertical[K comparable](head, node *ZSkipListNode[K], level int) bool {
	if node.backward == nil { // first element
		return head.level[level].span >= 1
	}
	var tranversed = 0
	var prev *ZSkipListNode[K]
	var x = node.backward
	for x != nil {
		if level >= len(x.level) {
			return true
		}
		if x.level[level].span > tranversed {
			return true
		}
		tranversed++
		prev = x
		x = x.backward
	}
	if prev != nil && level < len(prev.level) {
		return prev.level[level].span >= tranversed
	}
	return false
}

func prePadding(line *bytes.Buffer, n int) {
	for i := 0; i < n; i++ {
		line.WriteByte(' ')
	}
}
