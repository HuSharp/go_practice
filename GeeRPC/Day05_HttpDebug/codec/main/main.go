package main

import (
	"context"
	"geerpc"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)


type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}


// 在 startServer 中使用了信道 addr，确保服务端端口监听成功，客户端再发起请求。
func startServer(addr chan string)  {
	var foo Foo
	// pick up a free port
	listen, err := net.Listen("tcp", ":12580")
	if err != nil {
		log.Fatal("[startServer] network err: ", err)
	}

	if err := geerpc.Register(&foo); err != nil {
		log.Fatal("[startServer] register err: ", err)
	}
	
	geerpc.HandleHTTP()
	log.Println("[startServer] start rpc server on", listen.Addr())
	addr <- listen.Addr().String()
	_ = http.Serve(listen, nil)
}

func call(addrCh chan string) {
	client, _ := geerpc.DialHTTP("tcp", <-addrCh)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

func main() {
	// 客户端首先发送 Option 进行协议交换，
	// 接下来发送消息头 h := &codec.Header{}，和消息体 geeRpc req ${h.Seq}。

	// 删除所有标志，包括时间戳
	log.SetFlags(0)
	addr := make(chan string)
	go call(addr)
	startServer(addr)
}