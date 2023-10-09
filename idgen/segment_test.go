// Copyright © 2020-present ichenq@gmail.com All rights reserved.
// Distributed under the terms and conditions of the BSD License.
// See accompanying files LICENSE.

package idgen

import (
	"context"
	"testing"
	"time"
)

func TestSeqIDRedisSimple(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	runSeqIDTestSimple(t, ctx, "redis", "/uuid/counter1")
	// Output:
	//  QPS 4792038.70/s
}

func TestSeqIDRedisConcurrent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	runSeqIDTestConcurrent(t, ctx, "redis", "/uuid/counter2")
	// Output:
	//  QPS 2222480.03/s
}

func TestSeqIDRedisDistributed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	runSeqIDTestDistributed(t, ctx, "redis", "/uuid/counter3")
	// Output:
	//  QPS 2462537.90/s
}

func TestSeqIDMongoSimple(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	runSeqIDTestSimple(t, ctx, "mongo", "uuid_counter1")
	// Output:
	//  QPS 3325821.41/s
}

func TestSeqIDMongoConcurrent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	runSeqIDTestConcurrent(t, ctx, "mongo", "uuid_counter2")
	// Output:
	//  QPS 1880948.45/s
}

func TestSeqIDMongoDistributed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	runSeqIDTestDistributed(t, ctx, "mongo", "uuid_counter3")
	// Output:
	//  QPS 2380477.59/s
}
