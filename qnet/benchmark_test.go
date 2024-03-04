// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package qnet

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/datetime"
	"gopkg.in/svrkit.v1/slog"
)

const (
	kBenchConnNum       = 10
	kTotalBenchMsgCount = 1000000
)

func listenTestServer(t *testing.T, addr string, n int, bus chan net.Conn) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Errorf("Listen: %s %v", addr, err)
		return
	}
	for n > 0 {
		conn, err := ln.Accept()
		if err != nil {
			slog.Errorf("accept error: %v", err)
			return
		}
		bus <- conn
		n--
	}
}

func startBenchServer(t *testing.T, ctx context.Context, addr string) {
	var incoming = make(chan *NetMessage, 8192)
	var errChan = make(chan *Error, kBenchConnNum)
	var bus = make(chan net.Conn, kBenchConnNum)
	go listenTestServer(t, addr, kBenchConnNum, bus)

	var autoId int32 = 1

	for {
		select {
		case conn := <-bus:
			var session = NewTcpSession(conn)
			session.Node = NodeID(autoId)
			session.RecvQueue = incoming
			session.ErrChan = errChan
			session.SetSendQueue(make(chan *NetMessage, 256))
			session.Go(context.Background(), true, true)
			autoId++

		case netErr := <-errChan:
			// handle connection error
			var endpoint = netErr.Endpoint
			endpoint.Close()
			//t.Logf("session %v close %v", endpoint.Node(), ne)

		case msg := <-incoming:
			var resp = msg.Clone()
			msg.Reply(resp.Command, resp.Data)

		case <-ctx.Done():
			// handle shutdown
			return
		}
	}
}

func startBenchClient(t *testing.T, conn net.Conn, msgCount int, totalRecvMsgCount *int64, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= msgCount; i++ {
		var buf bytes.Buffer
		var msg = new(NetMessage)
		msg.Command = uint32(i)
		msg.Seq = uint32(i)
		if err := EncodeMsgTo(msg, nil, &buf); err != nil {
			t.Errorf("Encode: %v", err)
			break
		}

		if _, err := conn.Write(buf.Bytes()); err != nil {
			t.Errorf("Write: %v", err)
			break
		}
	}
	var rd = bufio.NewReader(conn)
	for i := 1; i <= msgCount; i++ {
		var resp = AllocNetMessage()
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err := DecodeMsgFrom(rd, 8196, nil, resp); err != nil {
			t.Errorf("decode message: %v", err)
			break
		}
		atomic.AddInt64(totalRecvMsgCount, 1)
	}
}

func TestServerClientQPS(t *testing.T) {
	var address = "localhost:15334"

	var kEachConnMsgCount = kTotalBenchMsgCount / kBenchConnNum
	t.Logf("total msg %d, total conn %d, each msg %d", kTotalBenchMsgCount, kBenchConnNum, kEachConnMsgCount)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	go startBenchServer(t, ctx, address)

	var startTime = time.Now()
	fmt.Printf("start QPS benchmark %s\n", datetime.FormatTime(startTime))

	var wg sync.WaitGroup
	var totalRecvMsgCount int64
	for i := 0; i < kBenchConnNum; i++ {
		conn, err := net.DialTimeout("tcp", address, time.Millisecond*500)
		assert.Nil(t, err)
		if conn != nil {
			wg.Add(1)
			go startBenchClient(t, conn, kEachConnMsgCount, &totalRecvMsgCount, &wg)
		}
	}

	wg.Wait()

	var stopAt = time.Now()
	fmt.Printf("QPS benchmark finished %s\n", datetime.FormatTime(stopAt))
	var elapsed = stopAt.Sub(startTime)
	fmt.Printf("Send/recv %d message with %d clients cost %v\n", totalRecvMsgCount, kBenchConnNum, elapsed)
	var qps = float64(totalRecvMsgCount) / (float64(elapsed) / float64(time.Second))
	fmt.Printf("avg QPS: %.2f/s\n", qps)
	fmt.Printf("Benchmark finished\n")

	// Output:
	// 	avg QPS: 64243.97/s
}
