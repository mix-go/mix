package xrpc

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/mix-go/xrpc/testdata"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
)

type service struct {
	pb.UnimplementedOrderServer
}

func TestRPCServer_Serve(t *testing.T) {
	s := &RpcServer{
		Grpc: &Grpc{
			Addr: "0.0.0.0:50000",
			Registrar: func(s *grpc.Server) {
				pb.RegisterOrderServer(s, &service{})
			},
		},
		Gateway: &Gateway{ // Optional
			Addr: "0.0.0.0:50001",
			Registrar: func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
				pb.RegisterOrderHandler(context.Background(), mux, conn)
			},
		},
		Logger: nil,
	}
	s.Serve()
	// s.Shutdown()
}

func TestRpcServerTLS_Serve(t *testing.T) {
	dir, _ := os.Getwd()
	tlsConf, err := LoadServerTLSConfig(dir+"/certificates/ca.pem", dir+"/certificates/server.pem", dir+"/certificates/server.key")
	if err != nil {
		log.Fatal(err)
	}
	tlsClientConf, err := LoadClientTLSConfig(dir+"/certificates/ca.pem", dir+"/certificates/client.pem", dir+"/certificates/client.key")
	if err != nil {
		log.Fatal(err)
	}
	s := &RpcServer{
		Grpc: &Grpc{
			Addr: "0.0.0.0:50000",
			Registrar: func(s *grpc.Server) {
				pb.RegisterOrderServer(s, &service{})
			},
		},
		Gateway: &Gateway{ // Optional
			Addr: "0.0.0.0:50001",
			Registrar: func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
				pb.RegisterOrderHandler(context.Background(), mux, conn)
			},
		},
		Logger:          nil,
		TLSConfig:       tlsConf,
		TLSClientConfig: tlsClientConf,
	}
	s.Serve()
}
