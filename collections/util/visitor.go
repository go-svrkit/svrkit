// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package util

type Visitor[V any] func(value V) bool

type KeyValVisitor[K, V any] func(key K, value V) bool
