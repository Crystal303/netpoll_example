package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/cloudwego/netpoll"
)

func main() {
	network, address := "tcp", "127.0.0.1:8888"

	// 直接创建连接
	conn, err := netpoll.DialConnection(network, address, 50*time.Millisecond)
	if err != nil {
		panic("dial netpoll connection fail")
	}

	// 通过 dialer 创建连接
	dialer := netpoll.NewDialer()
	conn, err = dialer.DialConnection(network, address, 50*time.Millisecond)
	if err != nil {
		panic("dialer netpoll connection fail")
	}

	// 设置读事件回调
	conn.SetOnRequest(cb())

	// conn write & flush message
	conn.Writer().WriteBinary([]byte("hello world"))
	conn.Writer().Flush()

	time.Sleep(400 * time.Microsecond)
}

func cb() netpoll.OnRequest {
	return func(ctx context.Context, connection netpoll.Connection) error {
		log.Info("cb")
		return nil
	}
}
