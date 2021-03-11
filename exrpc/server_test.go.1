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
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	var opts []Option
	//opts = append(opts, WithServerCredentialsFromFile("./cer/server.pem", "./cer/server.key"))
	s, err := NewServer(log, opts...)
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	s.RegisterPreprocess(func(appId, appKey string, ctx context.Context) (context.Context, error) {
		fmt.Println(appId, appKey)
		return ctx, nil
	})
	ss := s.Generate()
	pb.RegisterHelloServer(ss, HelloService)
	ss.Serve(listen)
}

func TestClient_Dial(t *testing.T) {
	log := logger.NewLogger("/dev/stdout", "20060102150405")
	opts := []Option{
		//WithClientCredentialsFromFile("./cer/server.pem", "wallet"),
	}
	c, err := NewClient("127.0.0.1:10101", "12323123123", "afafasdfad", log, opts...)
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

// 定义helloService并实现约定的接口
type helloService struct{}

// HelloService ...
var HelloService = helloService{}

// SayHello 实现Hello服务接口
func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s.\n", in.Name)
	return resp, nil
}
