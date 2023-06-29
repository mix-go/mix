package xrpc

import (
	"context"
	"crypto/tls"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"net"
	"net/http"
	"strings"
	"time"
)

type RpcServer struct {
	GrpcAddr    string
	GatewayAddr string

	// Optional
	Logger logging.Logger
	// No content: logging.StartCall, logging.FinishCall
	// With content: logging.PayloadReceived, logging.PayloadSent
	LoggableEvents []logging.LoggableEvent

	GrpcRegistrar    func(server *grpc.Server)
	GatewayRegistrar func(mux *runtime.ServeMux, conn *grpc.ClientConn)

	GrpcServer    *grpc.Server
	GatewayServer *http.Server

	// Use xrpc.NewTLSConfig or xrpc.LoadTLSConfig to create
	TLSConfig *tls.Config
	// Use xrpc.NewTLSClientConfig or xrpc.LoadTLSClientConfig to create
	// Not empty, gateway require
	TLSClientConfig *tls.Config

	// Additional server config
	ServerOptions []grpc.ServerOption
}

func (t *RpcServer) Serve() error {
	// listen
	listen, err := net.Listen("tcp", t.GrpcAddr)
	if err != nil {
		return err
	}

	// server
	srvOpts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
			PermitWithoutStream: true,            // Allow pings even when there are no active streams
		}),
	}
	if t.Logger != nil {
		logOpts := []logging.Option{
			logging.WithLogOnEvents(t.LoggableEvents...),
		}
		srvOpts = append(srvOpts,
			grpc.ChainUnaryInterceptor(
				logging.UnaryServerInterceptor(t.Logger, logOpts...),
			),
			grpc.ChainStreamInterceptor(
				logging.StreamServerInterceptor(t.Logger, logOpts...),
			))
	}
	if t.TLSConfig != nil {
		srvOpts = append(srvOpts, grpc.Creds(credentials.NewTLS(t.TLSConfig)))
	}
	if len(t.ServerOptions) > 0 {
		srvOpts = append(srvOpts, t.ServerOptions...)
	}
	s := grpc.NewServer(srvOpts...)
	t.GrpcRegistrar(s)
	go func() {
		if err := s.Serve(listen); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			panic(err)
		}
	}()

	// gRPC-Gateway 就是通过它来代理请求（将HTTP请求转为RPC请求）
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: true,
		}),
	}
	if t.TLSClientConfig != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(t.TLSClientConfig)))
	}
	addr := strings.ReplaceAll(t.GrpcAddr, "0.0.0.0", "127.0.0.1")
	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		return err
	}

	mux := runtime.NewServeMux()
	t.GatewayRegistrar(mux, conn)
	gwServer := &http.Server{
		Addr:    t.GatewayAddr,
		Handler: mux,
	}
	if t.TLSConfig != nil {
		gwServer.TLSConfig = t.TLSConfig
	}
	return gwServer.ListenAndServe()
}

func (t *RpcServer) Shutdown() error {
	t.GrpcServer.Stop()
	return t.GatewayServer.Shutdown(context.Background())
}
