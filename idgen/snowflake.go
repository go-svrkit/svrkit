// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

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
//	2位时钟回拨标记，时钟最多被回拨3次
//	36位时间戳（厘秒），~=21y288d, 最大可以表示到2045-10-10
//	13位服务器ID，最大服务器ID=8191
//	12位序列号，单个时间单位的最大分配数量（4096/毫秒）
const (
	SequenceBits       = 12
	WorkerIDBits       = 13
	TimeUnitBits       = 36
	ClockBackwardsBits = 2
	WorkIDMask         = 1<<WorkerIDBits - 1
	MaxSeqID           = (1 << SequenceBits) - 1
	TimestampShift     = WorkerIDBits + SequenceBits
	MaxTimeUnits       = (1 << TimeUnitBits) - 1
	BackwardsMaskShift = TimeUnitBits + WorkerIDBits + SequenceBits

	TimeUnit    = int64(time.Millisecond * 10)    // 厘秒（10毫秒）
	CustomEpoch = int64(1704038400 * time.Second) // 起始纪元 2024-01-01 00:00:00 UTC
)

var (
	ErrClockGoneBackwards = errors.New("clock gone backwards")
	ErrTimeUnitOverflow   = errors.New("clock time unit overflow")
	ErrUUIDIntOverflow    = errors.New("uuid integer overflow")
)

// 1，唯一和递增
// 分布式环境下部署的雪花ID只保证了唯一性，多个节点生成的ID范围无法保证顺序性（递增），因为ID的分段包含了服务ID，在同个时间
// 单元内(10ms)服务A按时钟后生成的ID会小于服务B按时钟先生成的ID。
//
// 2，时钟回拨问题
// 雪花算法的UUID的生成依赖系统时钟，如果系统时钟被回拨，会有潜在的生成重复ID的情况。
//   a, 系统时钟被人为回拨，需要业务层提供逻辑时钟机制
//   b, NTP同步和UTC闰秒(https://en.wikipedia.org/wiki/Leap_second)
//
// 设计中增加了时钟回拨标记位，可以让系统(重启前)在时钟被回拨时仍正确工作，但时钟回拨标记位有限并且未存档，时钟回拨位用完后
// 如果系统有重启，仍然会有ID重复的可能。
//
// 3. Worker ID的分配
// WorkerID跟原版实现一样人工分配，靠配置参数保证唯一性

type Snowflake struct {
	guard        sync.Mutex // 线程安全
	workerId     int64      // 服务ID
	seq          int64      // 当前序列号
	lastID       int64      // 最近生成的ID
	lastTimeUnit int64      // 最近的时间单元
	backwards    int64      // 时钟回拨标记
}

func NewSnowflake(workerId uint16) *Snowflake {
	if workerId == 0 {
		workerId = privateIP4()
		log.Printf("snowflake auto set worker id to %d\n", workerId)
	}
	return &Snowflake{
		workerId:     int64(workerId) & WorkIDMask,
		lastTimeUnit: currentTimeUnit(),
	}
}

func (sf *Snowflake) maskClockBackwards() error {
	log.Printf("Snowflake: time has gone backwards\n")
	if sf.backwards >= (1<<ClockBackwardsBits)-1 {
		return ErrClockGoneBackwards
	}
	sf.backwards++
	return nil
}

// sequence expired, tick to next time unit
func (sf *Snowflake) waitTilNext(lastTs int64) (int64, error) {
	var prevTs int64
	for i := 0; i < 100; i++ {
		time.Sleep(5 * time.Millisecond)
		var now = currentTimeUnit()
		if now > lastTs {
			return now, nil
		}
		// sleep期间时钟是否又被回拨了
		if now < lastTs || now < prevTs {
			if err := sf.maskClockBackwards(); err != nil {
				return 0, err
			} else {
				return now, nil // 标记了回拨后直接返回
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

	// 时钟回拨判断
	if curTimeUnits < sf.lastTimeUnit {
		// 如果只是回拨了1个时间单元，就等待一下时钟
		if curTimeUnits+1 == sf.lastTimeUnit {
			if ts, err := sf.waitTilNext(curTimeUnits); err != nil {
				return 0, err
			} else {
				curTimeUnits = ts
			}
		} else {
			if err := sf.maskClockBackwards(); err != nil {
				return 0, err
			}
		}
	}
	// 当前仍在同一个时间单元内，只增加序列号
	// 注意：时钟也可能发生了一个时间单元内(10ms)的回拨
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
	var id = (sf.backwards << BackwardsMaskShift) | (curTimeUnits << TimestampShift) | (sf.workerId << SequenceBits) | sf.seq
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

func currentTimeUnit() int64 {
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
