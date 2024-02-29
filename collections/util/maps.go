// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package util

import (
	"maps"
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
