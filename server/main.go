package main

import (
	"context"
	"os"
	"time"

	"github.com/cloudwego/netpoll"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})

	log.SetOutput(os.Stdout)
}

func main() {
	network, address := "tcp", "127.0.0.1:8888"
	log.Info(network, "@", address)

	// 创建 listener
	listener, err := netpoll.CreateListener(network, address)
	if err != nil {
		panic("create netpoll listener fail")
	}

	// handle: 连接读数据和处理逻辑
	var onRequest netpoll.OnRequest = handler
	var onPrepare netpoll.OnPrepare = prepare

	// options: EventLoop 初始化自定义配置项
	var opts = []netpoll.Option{
		netpoll.WithReadTimeout(1 * time.Second),
		netpoll.WithIdleTimeout(10 * time.Minute),
		netpoll.WithOnPrepare(onPrepare),
	}

	// 创建 EventLoop
	eventLoop, err := netpoll.NewEventLoop(onRequest, opts...)
	if err != nil {
		panic("create netpoll event-loop fail")
	}

	// 运行 Server
	err = eventLoop.Serve(listener)
	if err != nil {
		panic("netpoll server exit")
	}
}

// 读事件处理
func handler(ctx context.Context, connection netpoll.Connection) error {
	reader := connection.Reader()
	bts, err := reader.Next(reader.Len())
	if err != nil {
		return err
	}
	log.Infof("Key: %s, data: %s", ctx.Value(ctxKey{}), string(bts))

	connection.Write([]byte("dd"))
	return connection.Writer().Flush()
}

type ctxKey struct{}

var DefCtxKey ctxKey

func prepare(connection netpoll.Connection) context.Context {
	return context.WithValue(context.Background(), DefCtxKey, "context")
}
