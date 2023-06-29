package xrpc

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/mix-go/xrpc/openmix"
	"google.golang.org/grpc"
	"testing"
)

type service struct {
	pb.UnimplementedOrderServer
}

func TestRPCServer_Serve(t *testing.T) {
	s := &RpcServer{
		GrpcAddr:    "0.0.0.0:50000",
		GatewayAddr: "0.0.0.0:50001",
		Logger:      nil,
		GrpcRegistrar: func(s *grpc.Server) {
			pb.RegisterOrderServer(s, &service{})
		},
		GatewayRegistrar: func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
			pb.RegisterOrderHandler(context.Background(), mux, conn)
		},
	}
	s.Serve()
	// s.Shutdown()
}
