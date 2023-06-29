# XRPC

Assistant for gRPC and Gateway

## Install

go mod install

```
go get github.com/mix-go/xrpc@latest
```

Install the proto compiler

Download proto: https://github.com/protocolbuffers/protobuf/releases

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Install the grpc-gateway compilation tool

```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

Install googleapis proto

```
mkdir `go env GOPATH`/src/google
```

```
wget https://github.com/googleapis/googleapis/archive/refs/heads/master.zip -O googleapis.zip
unzip googleapis.zip
cp -R googleapis-master/google/api `go env GOPATH`/src/google
```

```
wget https://github.com/protocolbuffers/protobuf/archive/refs/heads/main.zip -O protobuf.zip
unzip protobuf.zip
cp -R protobuf-main/src/google/protobuf `go env GOPATH`/src/google
```

Goland settings: Settings > Languages & Frameworks > Protocol Buffers

Add Import Paths: `go env GOPATH`/src

## Best Practices

- method name: `RequestForRelease` PascalCase
- field Name: `string order_number = 1;` snake_case
- grpc gateway url (api url): `/v1/request_for_release` snake_case
- website url: `/request-for-release` kebab-case

```
// .proto
service Order {
  rpc RequestForRelease(ReleaseRequest) returns (ReleaseResponse) {
    option (google.api.http) = {
      post: "/v1/request_for_release"
      body: "*"
    };
  }
}

message ReleaseRequest {
  string order_number = 1;
  string requester = 2;
  string approver = 3;
}

message ReleaseResponse {
  int64 code = 1;
  string message = 2;
}
```

## Generate code

- go

```
generate-go.sh
```

## RPC Server

The necessary functions are encapsulated internally for unified management

- Start

```go
s := &xrpc.RpcServer{
    GrpcAddr:    "0.0.0.0:50000",
    GatewayAddr: "0.0.0.0:50001",
    Logger:      &RpcLogger{SugaredLogger: zapLogger},
    GrpcRegistrar: func(s *grpc.Server) {
        pb.RegisterOrderServer(s, &service{})
    },
    GatewayRegistrar: func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
        pb.RegisterOrderHandler(context.Background(), mux, conn)
    },
}
s.Serve()
```

- Shutdown

```go
s.Shutdown()
```

### TLS Server

We need to use it when we write core financial services

```go
tlsConf, err := xrpc.LoadTLSConfig("/certificates/ca.pem", "/certificates/server.pem", "/certificates/server.key")
if err != nil {
    log.Fatal(err)
}
tlsCliConf, err := xrpc.LoadTLSClientConfig("/certificates/ca.pem", "/certificates/client.pem", "/certificates/client.key")
if err != nil {
    log.Fatal(err)
}
s := &xrpc.RpcServer{
    GrpcAddr:    "0.0.0.0:50000",
    GatewayAddr: "0.0.0.0:50001",
    Logger:      &RpcLogger{SugaredLogger: zapLogger},
    GrpcRegistrar: func(s *grpc.Server) {
        pb.RegisterOrderServer(s, &service{})
    },
    GatewayRegistrar: func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
        pb.RegisterOrderHandler(context.Background(), mux, conn)
    },
    TLSConfig: tlsConf,
    TLSClientConfig: tlsCliConf, // not empty, gateway requires
}
s.Serve()
```

### Logger

Implement the following interfaces

```go
type Logger interface {
    Log(ctx context.Context, level Level, msg string, fields ...any)
}
```

Loggable Events

- No content: logging.StartCall, logging.FinishCall
- With content: logging.PayloadReceived, logging.PayloadSent

```go
s := &xrpc.RpcServer{
    LoggableEvents: []logging.LoggableEvent{logging.StartCall, logging.FinishCall},
}
```

Zap Logger

```go
type RpcLogger struct {
    *zap.SugaredLogger
}

func (t *RpcLogger) Log(ctx context.Context, level logging.Level, msg string, fields ...any) {
    f := make([]zap.Field, 0, len(fields)/2)
    for i := 0; i < len(fields); i += 2 {
        key := fields[i]
        value := fields[i+1]
        switch v := value.(type) {
        case string:
            f = append(f, zap.String(key.(string), v))
        case int:
            f = append(f, zap.Int(key.(string), v))
        case bool:
            f = append(f, zap.Bool(key.(string), v))
        default:
            f = append(f, zap.Any(key.(string), v))
        }
    }
    logger := t.Desugar().WithOptions(zap.AddCallerSkip(1)).With(f...)
    switch level {
    case logging.LevelDebug:
        logger.Debug(msg)
    case logging.LevelInfo:
        logger.Info(msg)
    case logging.LevelWarn:
        logger.Warn(msg)
    case logging.LevelError:
        logger.Error(msg)
    default:
        panic(fmt.Sprintf("unknown level %v", level))
    }
}
```

```go
logger := &RpcLogger{SugaredLogger: zapLogger}
```

## RPC Client

The necessary functions are encapsulated internally for unified management

> Please reuse this connection, you don't need to handle disconnect and reconnect, but you do need to handle request retry after errors.

```go
conn, err := xrpc.NewGrpcClient("127.0.0.1:50000")
if err != nil {
    log.Fatal(err)
}
defer conn.Close()
client := pb.NewOrderClient(conn)
ctx, _ := context.WithTimeout(context.Background(), xrpc.CallTimeout)
resp, err := client.RequestForRelease(ctx, &pb.ReleaseRequest{
    OrderNumber: "131243234",
})
```

## TLS RPC Client

We need to use it when we write core financial services

```go
tlsConf, err := xrpc.LoadTLSClientConfig("/certificates/ca.pem", "/certificates/client.pem", "/certificates/client.key")
if err != nil {
    log.Fatal(err)
}
conn, err := xrpc.NewGrpcClient("127.0.0.1:50000", grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)))
```

## TLS Gateway Client

```go
tlsConf, err := xrpc.LoadTLSClientConfig("/certificates/ca.pem", "/certificates/client.pem", "/certificates/client.key")
if err != nil {
    log.Fatal(err)
}
client := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: tlsConf,
    },
}
resp, err := client.Post("http://127.0.0.1:50001/v1/request_for_release", "application/json", strings.NewReader(`{"order_number":"123456789"}`))
fmt.Println(resp.Body)
```

Examples of other languages

```PHP
<?php
require __DIR__ . '/vendor/autoload.php';

use GuzzleHttp\Client;

$client = new Client([
    'cert' => [ '/certificates/client.pem', '']
]);
$response = $client->request('POST', 'http://127.0.0.1:50001/v1/request_for_release', ['body' => '{"order_number":"123456789"}']);
var_dump($response->getBody()->getContents());
```

## Self Signed Certificates

Modify subjectAltName in `generate-rsa.cnf`, the generated files are in the `certificates` directory.

```
generate-rsa.sh
```

Code new TLS Config

- Server

Load from file

```go
tlsConf, err := xrpc.LoadTLSConfig("/certificates/ca.pem", "/certificates/server.pem", "/certificates/server.key")
```

New by bytes

```go
tlsConf, err := xrpc.NewTLSConfig([]byte{}, []byte{}, []byte{})
```

- Client

Load from file

```go
tlsConf, err := xrpc.LoadTLSClientConfig("/certificates/ca.pem", "/certificates/client.pem", "/certificates/client.key")
```

New by bytes

```go
tlsConf, err := xrpc.NewTLSClientConfig([]byte{}, []byte{}, []byte{})
```
