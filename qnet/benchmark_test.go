// Copyright © 2020 ichenq@gmail.com All rights reserved.
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

	"gopkg.in/svrkit.v1/datetime"
	"gopkg.in/svrkit.v1/slog"
)

func serveListen(ln net.Listener, bus chan net.Conn) {
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := time.Second; tempDelay > max {
					tempDelay = max
				}
				slog.Warnf("Accept error: %v, retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			slog.Errorf("accept error: %v", err)
			return
		}
		bus <- conn
	}
}

func startBenchServer(t *testing.T, ctx context.Context, addr string, ready chan struct{}) {
	var incoming = make(chan *NetMessage, 6000)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Errorf("Listen: %s %v", addr, err)
		return
	}

	var errChan = make(chan *Error, 1024)
	var bus = make(chan net.Conn, 1024)
	go serveListen(ln, bus)

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

	var kBenchConnCount = 10
	var kTotalBenchMsgCount = 1000000
	var kEachConnMsgCount = kTotalBenchMsgCount / kBenchConnCount
	t.Logf("total msg %d, total conn %d, each msg %d", kTotalBenchMsgCount, kBenchConnCount, kEachConnMsgCount)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var ready = make(chan struct{})
	go startBenchServer(t, ctx, address, ready)
	<-ready // server listen ready

	var startTime = time.Now()
	fmt.Printf("start QPS benchmark %s\n", datetime.FormatTime(startTime))

	var wg sync.WaitGroup
	var totalRecvMsgCount int64
	for i := 0; i < kBenchConnCount; i++ {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			t.Fatalf("Dial %s: %v", address, err)
		}
		wg.Add(1)
		go startBenchClient(t, conn, kEachConnMsgCount, &totalRecvMsgCount, &wg)
	}

	wg.Wait()

	var stopAt = time.Now()
	fmt.Printf("QPS benchmark finished %s\n", datetime.FormatTime(stopAt))
	var elapsed = stopAt.Sub(startTime)
	fmt.Printf("Send/recv %d message with %d clients cost %v\n", totalRecvMsgCount, kBenchConnCount, elapsed)
	var qps = float64(totalRecvMsgCount) / (float64(elapsed) / float64(time.Second))
	fmt.Printf("avg QPS: %.2f/s\n", qps)
	fmt.Printf("Benchmark finished\n")

	// Output:
	// 	avg QPS: 64243.97/s
}
