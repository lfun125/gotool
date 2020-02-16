package exrpc_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	pb "github.com/jergoo/go-grpc-example/proto/hello"
	. "github.com/lfun125/gotool/exrpc"
	"github.com/lfun125/gotool/logger"
)

func TestServer(t *testing.T) {
	log := logger.NewLogger("/dev/stdout", "20060102150405")
	listen, err := net.Listen("tcp", ":10101")
	fmt.Println(123)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	s, err := NewServer(log)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	s.SetAuth(func(appId, appKey string) error {
		fmt.Println(appId, appKey)
		return nil
	})
	ss := s.Generate()
	pb.RegisterHelloServer(ss, HelloService)
	ss.Serve(listen)
}

func TestClient_Dial(t *testing.T) {
	log := logger.NewLogger("/dev/stdout", "20060102150405")
	c, err := NewClient("127.0.0.1:10101", "12323123123", "afafasdfad", log)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := c.Dial()
	if err != nil {
		t.Fatal(err)
	}
	hc := pb.NewHelloClient(conn)
	in := &pb.HelloRequest{
		Name: "abcdefg",
	}
	resp, err := hc.SayHello(context.Background(), in)
	fmt.Println(err)
	fmt.Println(resp)
}
