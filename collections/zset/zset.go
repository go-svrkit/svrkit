// ZSETs are ordered sets using two data structures to hold the same elements
// in order to get O(log(N)) INSERT and REMOVE operations into a sorted
// data structure.
//
// https://github.com/redis/redis/blob/6.2.14/src/t_zset.c

package zset

import (
	"cmp"
)

// SortedSet 跳表实现的有序字典
type SortedSet[T comparable] struct {
	dict map[T]int64   // Value and score
	zsl  *ZSkipList[T] // indexed linked list
}

func New[T comparable](comparator Comparator[T]) *SortedSet[T] {
	return &SortedSet[T]{
		dict: make(map[T]int64),
		zsl:  NewZSkipList[T](comparator),
	}
}

func NewCmp[T cmp.Ordered]() *SortedSet[T] {
	return New[T](cmp.Compare[T])
}

func (s *SortedSet[T]) Len() int {
	return s.zsl.Len()
}

// Add 添加或者更新一个元素的score
func (s *SortedSet[T]) Add(ele T, score int64) bool {
	curscore, found := s.dict[ele]
	if found {
		// Remove and re-insert when score changes.
		if curscore != score {
			var znode = s.zsl.Delete(curscore, ele)
			s.zsl.Insert(score, znode.Ele)
			s.dict[ele] = score
		}
	} else {
		var znode = s.zsl.Insert(score, ele)
		s.dict[ele] = znode.Score
	}
	return true
}

// Remove 删除一个元素
func (s *SortedSet[T]) Remove(ele T) bool {
	score, found := s.dict[ele]
	if found {
		delete(s.dict, ele)
		s.zsl.Delete(score, ele) // Delete from skiplist
		return true
	}
	return false // No such element found
}

// RemoveRangeByScore 删除score区间[min, max]的元素
func (s *SortedSet[T]) RemoveRangeByScore(min, max int64) int {
	if min > max {
		return 0
	}
	return s.zsl.DeleteRangeByScore(min, max, s.dict)
}

// RemoveRangeByRank 删除排名在[start, end]之间的元素，排名从1开始
func (s *SortedSet[T]) RemoveRangeByRank(start, end int) int {
	var llen = s.zsl.length
	if start < 0 {
		start = llen + start
	}
	if end < 0 {
		end = llen + end
	}
	if start < 0 {
		start = 0
	}
	if start > end || start >= llen {
		return 0
	}
	if end >= llen {
		end = llen - 1
	}
	return s.zsl.DeleteRangeByRank(start+1, end+1, s.dict)
}

// Count score在[min, max]之间的元素数量
func (s *SortedSet[T]) Count(min, max int64) int {
	if min > max {
		return 0
	}
	// Find first element in range
	zn := s.zsl.FirstInRange(min, max)

	// Use rank of first element, if any, to determine preliminary count
	if zn != nil {
		var rank = s.zsl.GetRank(zn.Score, zn.Ele)
		var count = s.zsl.length - (rank - 1)

		// Find last element in range
		zn = s.zsl.LastInRange(min, max)

		// Use rank of last element, if any, to determine the actual count
		if zn != nil {
			rank = s.zsl.GetRank(zn.Score, zn.Ele)
			count -= s.zsl.length - rank
		}
		return count
	}
	return 0
}

// GetRank 返回元素的排名，排名从0开始，如果元素不在zset里，返回-1
func (s *SortedSet[T]) GetRank(ele T, reverse bool) int {
	score, found := s.dict[ele]
	if found {
		var llen = s.zsl.Len()
		var rank = s.zsl.GetRank(score, ele)
		// assert rank != 0
		if reverse {
			return llen - rank
		}
		return rank - 1
	}
	return -1
}

// GetScore 获取元素的score
func (s *SortedSet[T]) GetScore(ele T) int64 {
	if score, found := s.dict[ele]; found {
		return score
	}
	return 0
}

// GetRange 返回排名在[start, end]之间的所有元素
func (s *SortedSet[T]) GetRange(start, end int, reverse bool) []T {
	var llen = s.zsl.length
	if start < 0 {
		start = llen + start
	}
	if end < 0 {
		end = llen + end
	}
	if start < 0 {
		start = 0
	}
	if start > end || start >= llen {
		return nil
	}
	if end >= llen {
		end = llen - 1
	}
	var rangeLen = end - start + 1
	var node *ZSkipListNode[T]
	// Check if starting point is trivial, before doing log(N) lookup.
	if reverse {
		node = s.zsl.tail
		if start > 0 {
			node = s.zsl.GetElementByRank(llen - start)
		}
	} else {
		node = s.zsl.head.level[0].forward
		if start > 0 {
			node = s.zsl.GetElementByRank(start + 1)
		}
	}
	var result = make([]T, 0, rangeLen)
	for rangeLen > 0 {
		result = append(result, node.Ele)
		if reverse {
			node = node.backward
		} else {
			node = node.level[0].forward
		}
		rangeLen--
	}
	return result
}

// GetRangeByScore 获取score在[min, max]之间的所有元素
func (s *SortedSet[T]) GetRangeByScore(min, max int64, reverse bool) []T {
	if min > max {
		return nil
	}
	var node *ZSkipListNode[T]
	// If reversed, get the last node in range as starting point
	if reverse {
		node = s.zsl.LastInRange(min, max)
	} else {
		node = s.zsl.FirstInRange(min, max)
	}
	if node == nil {
		return nil
	}
	var result []T
	for node != nil {
		// Abort when the node is no longer in range
		if reverse {
			if node.Score < min {
				break
			}
		} else {
			if node.Score > max {
				break
			}
		}

		result = append(result, node.Ele)

		// Move to next node
		if reverse {
			node = node.backward
		} else {
			node = node.level[0].forward
		}
	}
	return result
}
