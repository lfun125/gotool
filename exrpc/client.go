package exrpc

import (
	"context"

	"github.com/lfun125/gotool/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	tls         bool
	log         *logger.Logger
	credentials credentials.TransportCredentials
	address     string
	appID       string
	appKey      string
}

func NewClient(address, appID, appKey string, log *logger.Logger, options ...option) (*Client, error) {
	c := &Client{}
	c.address = address
	c.appID = appID
	c.appKey = appKey
	c.log = log
	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Client) setCredentials(credentials credentials.TransportCredentials) {
	c.credentials = credentials
}

func (c *Client) setTls(tls bool) {
	c.tls = tls
}

func (c *Client) Dial() (conn *grpc.ClientConn, err error) {
	var opts []grpc.DialOption
	if c.tls && c.credentials != nil {
		opts = append(opts, grpc.WithTransportCredentials(c.credentials))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(
		opts,
		grpc.WithPerRPCCredentials(&customCredential{
			AppID:  c.appID,
			AppKey: c.appKey,
		}),
		grpc.WithUnaryInterceptor(ClientInterceptor),
	)
	conn, err = grpc.Dial(c.address, opts...)
	return
}

// customCredential 自定义认证
type customCredential struct {
	AppID  string
	AppKey string
}

// GetRequestMetadata 实现自定义认证接口
func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_id":  c.AppID,
		"app_key": c.AppKey,
	}, nil
}

// RequireTransportSecurity 自定义认证是否开启TLS
func (c customCredential) RequireTransportSecurity() bool {
	return false
}

// interceptor 客户端拦截器
func ClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
