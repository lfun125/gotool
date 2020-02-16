package exrpc

import (
	"context"
	"fmt"

	pb "github.com/jergoo/go-grpc-example/proto/hello"
)

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
