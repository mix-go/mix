package xrpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

var (
	DialTimeout = 5 * time.Second

	CallTimeout = 5 * time.Second
)

func NewGrpcClient(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()
	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: true,
		}),
	}
	if len(opts) > 0 {
		dialOpts = append(dialOpts, opts...)
	}
	c, err := grpc.DialContext(ctx, addr, dialOpts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}
