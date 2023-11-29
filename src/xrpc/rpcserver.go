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
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"strings"
	"time"
)

type RpcServer struct {
	// Required
	*Grpc

	// Optional
	*Gateway

	// Optional
	Logger logging.Logger

	// Optional, Use xrpc.NewServerTLSConfig or xrpc.LoadServerTLSConfig to create
	TLSConfig *tls.Config

	// Optional, Use xrpc.NewClientTLSConfig or xrpc.LoadClientTLSConfig to create
	TLSClientConfig *tls.Config
}

type Gateway struct {
	// Required
	Addr string

	// Required
	Registrar func(mux *runtime.ServeMux, conn *grpc.ClientConn)

	Server *http.Server

	// Additional server config
	// Optional
	ServeMuxOptions []runtime.ServeMuxOption
}

type Grpc struct {
	// Required
	Addr string

	// No content: logging.StartCall, logging.FinishCall
	// With content: logging.PayloadReceived, logging.PayloadSent
	LoggableEvents []logging.LoggableEvent

	// Required
	Registrar func(server *grpc.Server)

	Listener net.Listener

	Server *grpc.Server

	// Additional server config
	// Optional
	ServerOptions []grpc.ServerOption
}

func (t *RpcServer) Serve() error {
	// listen
	listen, err := net.Listen("tcp", t.Grpc.Addr)
	if err != nil {
		return err
	}
	t.Grpc.Listener = listen

	// grpc server
	srvOpts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
			PermitWithoutStream: true,            // Allow pings even when there are no active streams
		}),
	}
	if t.Logger != nil {
		logOpts := []logging.Option{
			logging.WithLogOnEvents(t.Grpc.LoggableEvents...),
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
	if len(t.Grpc.ServerOptions) > 0 {
		srvOpts = append(srvOpts, t.Grpc.ServerOptions...)
	}
	s := grpc.NewServer(srvOpts...)
	t.Grpc.Registrar(s)
	serve := func() {
		if err := s.Serve(listen); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			panic(err)
		}
	}
	if t.Gateway == nil {
		serve()
		return nil
	}
	go serve()

	// grpc client
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
	addr := strings.ReplaceAll(t.Grpc.Addr, "0.0.0.0", "127.0.0.1")
	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		return err
	}

	// gateway server
	muxOpts := []runtime.ServeMuxOption{
		// Format for using proto names in json https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/#using-proto-names-in-json
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: false, // Omit 0 values, such as 0, "" or null
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	}
	if t.Logger != nil {
		customHTTPError := func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			defer runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
			if err == nil {
				return
			}
			st := status.Convert(err)
			t.Logger.Log(r.Context(), logging.LevelInfo, "gateway error", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr, "code", st.Code(), "message", st.Message(), "details", st.Details())
		}
		muxOpts = append(muxOpts, runtime.WithErrorHandler(customHTTPError))
	}
	if len(t.Grpc.ServerOptions) > 0 {
		muxOpts = append(muxOpts, t.Gateway.ServeMuxOptions...)
	}
	mux := runtime.NewServeMux(muxOpts...)
	t.Gateway.Registrar(mux, conn)
	requestLogger := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			statusCode := w.(interface{ StatusCode() int }).StatusCode()
			t.Logger.Log(r.Context(), logging.LevelInfo, "gateway request", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr, "status_code", statusCode)
		})
	}
	gateway := &http.Server{
		Addr:    t.Gateway.Addr,
		Handler: requestLogger(mux),
	}
	if t.TLSConfig != nil {
		gateway.TLSConfig = t.TLSConfig
		return gateway.ListenAndServeTLS("", "")
	}
	return gateway.ListenAndServe()
}

func (t *RpcServer) Shutdown() error {
	t.Grpc.Server.Stop()
	_ = t.Grpc.Listener.Close()
	return t.Gateway.Server.Shutdown(context.Background())
}
