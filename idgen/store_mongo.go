// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package idgen

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionName = "uuid"
)

// Counter 计数器
type Counter struct {
	Label string `bson:"label"`   // 识别符
	Value int64  `bson:"counter"` // 计数器
}

// MongoDBCounter 使用MongoDB的自增计数器
type MongoDBCounter struct {
	guard       sync.Mutex    //
	uri         string        // 连接uri
	db          string        // DB名称
	label       string        // 识别符
	lastCounter int64         // 保存最近一次生成的counter
	cli         *mongo.Client //
}

func NewMongoDBCounter(uri, db, label string) CounterStore {
	return &MongoDBCounter{
		uri:         uri,
		db:          db,
		label:       label,
		lastCounter: -1,
	}
}

func (s *MongoDBCounter) Init(ctx context.Context) error {
	if err := s.createConn(ctx); err != nil {
		return err
	}
	return nil
}

func (s *MongoDBCounter) createConn(ctx context.Context) error {
	clientOpts := options.Client().ApplyURI(s.uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return err
	}
	s.cli = client
	return nil
}

func (s *MongoDBCounter) Close() error {
	s.guard.Lock()
	s.cli.Disconnect(context.Background())
	s.cli = nil
	s.guard.Unlock()
	return nil
}

func (s *MongoDBCounter) Incr(ctx context.Context) (int64, error) {
	s.guard.Lock()
	defer s.guard.Unlock()

	var ctr = &Counter{
		Label: s.label,
	}
	// 把counter自增再读取最新的counter
	if err := s.IncrementAndLoad(ctx, 1, ctr); err != nil {
		return 0, err
	}
	if s.lastCounter >= ctr.Value {
		return 0, ErrCounterOutOfRange
	}
	s.lastCounter = ctr.Value
	return ctr.Value, nil
}

// IncrementAndLoad 把counter自增再读取最新的counter
// https://docs.mongodb.com/manual/core/write-operations-atomicity/
func (s *MongoDBCounter) IncrementAndLoad(ctx context.Context, n int, ctr *Counter) error {
	var filter = bson.M{"label": s.label}
	var update = bson.M{
		"$setOnInsert": bson.M{
			"label": ctr.Label,
		},
		"$inc": bson.M{"counter": n},
	}
	var opt = options.FindOneAndUpdate()
	opt.SetUpsert(true).SetReturnDocument(options.After)

	// https://www.mongodb.com/docs/manual/reference/method/db.collection.findOneAndUpdate/
	result := s.cli.Database(s.db).Collection(CollectionName).FindOneAndUpdate(ctx, filter, update, opt)
	if err := result.Err(); err != nil {
		return err
	}
	if err := result.Decode(ctr); err != nil {
		return err
	}
	return nil
}
