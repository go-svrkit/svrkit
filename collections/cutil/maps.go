// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package cutil

import (
	"cmp"
	"maps"
	"slices"
)

// MapKeys 返回map的key列表
func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	var keys = make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapValues 返回map的value列表
func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	var values = make([]V, 0, len(m))
	for _, val := range m {
		values = append(values, val)
	}
	return values
}

// MapKeyValues 返回map的key和value列表
func MapKeyValues[M ~map[K]V, K comparable, V any](m M) ([]K, []V) {
	var keys = make([]K, 0, len(m))
	var values = make([]V, 0, len(m))
	for k, val := range m {
		keys = append(keys, k)
		values = append(values, val)
	}
	return keys, values
}

// MapOrderedKeys 返回map里已排序的key列表
func MapOrderedKeys[M ~map[K]V, K cmp.Ordered, V any](m M) []K {
	var keys = MapKeys(m)
	slices.Sort(keys)
	return keys
}

// MapOrderedValues 返回map里按key排序的value列表
func MapOrderedValues[M ~map[K]V, K cmp.Ordered, V any](m M) []V {
	var keys = MapKeys(m)
	slices.Sort(keys)
	var values = make([]V, 0, len(m))
	for _, k := range keys {
		values = append(values, m[k])
	}
	return values
}

// MapOrderedKeyValues 返回map里已排序的key和value列表
func MapOrderedKeyValues[M ~map[K]V, K cmp.Ordered, V any](m M) ([]K, []V) {
	var keys = MapKeys(m)
	slices.Sort(keys)
	var values = make([]V, 0, len(m))
	for _, k := range keys {
		values = append(values, m[k])
	}
	return keys, values
}

// MapUnion 返回两个map的并集, copy of a ∪ b
func MapUnion[M ~map[K]V, K comparable, V any](a, b M) M {
	var result = make(M, len(a)+len(b))
	maps.Copy(result, a)
	maps.Copy(result, b)
	return result
}

// MapIntersect 返回两个map的交集, a ∩ b
func MapIntersect[M ~map[K]V, K, V comparable](a, b M) M {
	var result = make(map[K]V)
	for k, v := range a {
		if val, ok := b[k]; ok && val == v {
			result[k] = v
		}
	}
	return result
}
