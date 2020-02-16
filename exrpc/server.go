package exrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/lfun125/gotool/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type AuthFunc func(appId, appKey string) error

type Server struct {
	tls         bool
	log         *logger.Logger
	credentials credentials.TransportCredentials
	auth        AuthFunc
}

func NewServer(log *logger.Logger, options ...option) (*Server, error) {
	s := &Server{}
	s.log = log
	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Server) setCredentials(credentials credentials.TransportCredentials) {
	s.credentials = credentials
}

func (s *Server) setTls(tls bool) {
	s.tls = tls
}

func (s *Server) SetAuth(auth AuthFunc) {
	s.auth = auth
}

func (s Server) Generate() *grpc.Server {
	var opts []grpc.ServerOption
	if s.tls && s.credentials != nil {
		opts = append(opts, grpc.Creds(s.credentials))
	}
	opts = append(opts, grpc.UnaryInterceptor(s.serverInterceptor))
	return grpc.NewServer(opts...)
}

func (s Server) serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (data interface{}, err error) {
	requestTime := time.Now()
	defer func() {
		if e := recover(); e != nil {
			tracks := Tracks()
			s.log.With("track_list", tracks).Error(e)
			err = fmt.Errorf("panic: %v", e)
		}
		st := "SUCCESS"
		var e string
		if err != nil {
			st = "FAILED"
			e = fmt.Sprintf("error: %s", err.Error())
		}
		s.log.Info(fmt.Sprintf("[%s] [%s] [%s] %s", st, time.Now().Sub(requestTime), info.FullMethod, e))
	}()
	appID, appKey := s.getAuthInfo(ctx)
	if s.auth != nil {
		if err = s.auth(appID, appKey); err != nil {
			return
		}
	}
	data, err = handler(ctx, req)
	return
}

func (s Server) getAuthInfo(ctx context.Context) (appId, appKey string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}
	if val, ok := md["app_id"]; ok {
		appId = val[0]
	}
	if val, ok := md["app_key"]; ok {
		appKey = val[0]
	}
	return
}
