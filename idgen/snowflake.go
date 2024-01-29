// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package idgen

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

// 一个64位UUID由以下部分组成
//
//	1位符号位
//	1位时钟回拨标记
//	37位时间戳（厘秒），最大可以表示到2065-07-20
//	13位服务器ID，最大服务器ID=8191
//	12位序列号，单个时间单位的最大分配数量（409/毫秒）
const (
	SequenceBits       = 12
	WorkerIDBits       = 13
	TimeUnitBits       = 37
	WorkIDMask         = 1<<WorkerIDBits - 1
	MaxSeqID           = (1 << SequenceBits) - 1
	TimestampShift     = WorkerIDBits + SequenceBits
	MaxTimeUnits       = (1 << TimeUnitBits) - 1
	BackwardsMaskShift = TimeUnitBits + WorkerIDBits + SequenceBits

	TimeUnit    = int64(time.Millisecond * 10)    // 厘秒（10毫秒）
	CustomEpoch = int64(1640966400 * time.Second) // 起始纪元 2022-01-01 00:00:00 UTC
)

var (
	ErrClockGoneBackwards = errors.New("clock gone backwards")
	ErrTimeUnitOverflow   = errors.New("clock time unit overflow")
	ErrUUIDIntOverflow    = errors.New("uuid integer overflow")
)

// 雪花算法的UUID的生成依赖系统时钟，如果系统时钟被回拨，会有潜在的生成重复ID的情况
// 	1，系统时钟被人为回拨，需要业务层提供逻辑时钟机制
// 	2，NTP同步和UTC闰秒(https://en.wikipedia.org/wiki/Leap_second)
//
// 设计中增加了时钟回拨标记位，可以让系统(重启前)在时钟被回拨时仍正确工作
// 但时钟回拨标记位有限并且未存档，时钟回拨后如果在系统重启前没有恢复，仍然会有ID重复的可能
//
// 分布式部署下雪花ID只保证了唯一性，无法保证顺序性（递增）
// 因为ID的分段包含了服务ID，在同个时间单元内(10ms)服务A按时钟后生成的ID会小于服务B按时钟先生成的ID；

type Snowflake struct {
	guard         sync.Mutex // 线程安全
	workerId      int64      // 服务器ID
	seq           int64      // 当前序列号
	lastID        int64      // 最近生成的ID
	lastTimeUnit  int64      // 最近的时间单元
	backwardsMask int64      // 允许时钟被回拨1次
}

func NewSnowflake(workerId uint16) *Snowflake {
	if workerId == 0 {
		workerId = privateIP4()
		log.Printf("snowflake auto set worker id to %d", workerId)
	}
	return &Snowflake{
		workerId:     int64(workerId) & WorkIDMask,
		lastTimeUnit: currentTimeUnit(),
	}
}

func (sf *Snowflake) clockBackwards() error {
	log.Printf("Snowflake: time has gone backwards")
	if sf.backwardsMask > 0 {
		return ErrClockGoneBackwards
	}
	sf.backwardsMask = 1 << BackwardsMaskShift
	return nil
}

// sequence expired, tick to next time unit
func (sf *Snowflake) waitTilNext(lastTs int64) (int64, error) {
	var prevTs int64
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Millisecond)
		var now = currentTimeUnit()
		if now > lastTs {
			return now, nil
		}
		// 时钟是否被回拨
		if now < lastTs || now < prevTs {
			if err := sf.clockBackwards(); err != nil {
				return 0, err
			} else {
				return now, nil // 已经设置了回拨标记
			}
		}
		prevTs = now
	}
	return 0, ErrClockGoneBackwards
}

func (sf *Snowflake) Next() (int64, error) {
	sf.guard.Lock()
	defer sf.guard.Unlock()

	var curTimeUnits = currentTimeUnit()
	if curTimeUnits > MaxTimeUnits {
		return 0, ErrTimeUnitOverflow // 已经是2065年
	}

	// 判断时钟是否被回拨了
	// 注意：这里如果时钟回拨间隔在一个时间单元内（10ms）不会被发觉
	if curTimeUnits < sf.lastTimeUnit {
		if err := sf.clockBackwards(); err != nil {
			return 0, err
		}
	}
	// 当前仍在同一个时间单元内，只增加序列号
	if curTimeUnits == sf.lastTimeUnit {
		sf.seq++
		// 如果序列号满了，需要等待时钟走到下一个时间单元
		if sf.seq > MaxSeqID {
			if ts, err := sf.waitTilNext(curTimeUnits); err != nil {
				return 0, err
			} else {
				curTimeUnits = ts
				sf.seq = 0
			}
		}
	} else {
		// 已经是新的时间单元
		sf.seq = 0
	}

	sf.lastTimeUnit = curTimeUnits
	var id = sf.backwardsMask | (curTimeUnits << TimestampShift) | (sf.workerId << SequenceBits) | sf.seq
	if id <= sf.lastID {
		return 0, ErrUUIDIntOverflow
	}
	sf.lastID = id
	return id, nil
}

func (sf *Snowflake) MustNext() int64 {
	if n, err := sf.Next(); err != nil {
		panic(err)
	} else {
		return n
	}
}

// make it easier to mock time in unit tests
var currentTimeUnit = func() int64 {
	return (time.Now().UTC().UnixNano() - CustomEpoch) / TimeUnit
}

// lower 16-bits of IPv4
func privateIP4() uint16 {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return 0
	}
	var isPrivateIPv4 = func(ip net.IP) bool {
		return ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 &&
			ip[1] < 32) || ip[0] == 192 && ip[1] == 168
	}
	for _, a := range addr {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		var ip = ipnet.IP.To4()
		if ip != nil && isPrivateIPv4(ip) {
			return uint16(ip[2])<<8 + uint16(ip[3])
		}
	}
	return 0
}
