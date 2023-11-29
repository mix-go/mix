> Produced by OpenMix: [https://openmix.org](https://openmix.org/mix-go)

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

Install google proto

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

Add Import Paths: `$GOPATH`/src

## Best Practices

- .proto [style](https://protobuf.dev/programming-guides/style/#message-field-names)
  - service name, rpc name, message name: `AppMessages` PascalCase
  - message field name: `string parse_mode = 1;` snake_case
- urls:
  - website url: `/send-message` kebab-case
  - grpc gateway url: inner api: `/internal/send_message` snake_case
  - grpc gateway url: open api: `/v1/send_message` snake_case
- aws secrets manager:
  - name: `Service-Test` Pascal-Case
  - key: `googleapis_credentials` snake_case
- .yaml
  - file name: `config_test.yaml` snake_case
  - field name: `clientId` camelCase
- mysql:
  - table name: `app_messages` snake_case
  - field name: `client_id` snake_case
- mongodb:
  - table name: `app_messages` snake_case
  - field name: Unrestricted, as it depends on the 3rd party, storing raw data
- docker:
  - container name:  `express-gateway` kebab-case

```protobuf
service AppMessages {
  rpc Send(SendRequest) returns (SendResponse) {
    option (google.api.http) = {
      post: "/internal/send_message"
      body: "*"
    };
  }
}

message SendRequest {
  string text = 1;
  string parse_mode = 2;
}

message SendResponse {
  int64 message_id = 1;
}
```

## Generate code

```
generate-pb.sh
```

## RPC Server

The necessary functions are encapsulated internally for unified management

- Start

```go
s := &xrpc.RpcServer{
    GrpcServer: &xrpc.GrpcServer{
        Addr: "0.0.0.0:50000",
        Registrar: func(s *grpc.Server) {
            pb.RegisterOrderServer(s, &service{})
        },
    },
    GatewayServer: &xrpc.GatewayServer{ // Optional
        Addr: "0.0.0.0:50001",
        Registrar: func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
            pb.RegisterOrderHandler(context.Background(), mux, conn)
        },
    },
    Logger: &RpcLogger{SugaredLogger: zapLogger},
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
tlsConf, err := xrpc.LoadServerTLSConfig("/certificates/ca.pem", "/certificates/server.pem", "/certificates/server.key")
if err != nil {
    log.Fatal(err)
}
tlsClientConf, err := xrpc.LoadClientTLSConfig("/certificates/ca.pem", "/certificates/client.pem", "/certificates/client.key")
if err != nil {
    log.Fatal(err)
}
s := &xrpc.RpcServer{
    GrpcServer: &xrpc.GrpcServer{
        Addr: "0.0.0.0:50000",
        Registrar: func(s *grpc.Server) {
            pb.RegisterOrderServer(s, &service{})
        },
    },
    GatewayServer: &xrpc.GatewayServer{ // Optional
        Addr: "0.0.0.0:50001",
        Registrar: func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
            pb.RegisterOrderHandler(context.Background(), mux, conn)
        },
    },
    Logger: &RpcLogger{SugaredLogger: zapLogger},
    TLSConfig: tlsConf,
    TLSClientConfig: tlsClientConf,
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
    GrpcServer: &xrpc.GrpcServer{
        LoggableEvents: []logging.LoggableEvent{logging.StartCall, logging.FinishCall},
    }
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
    OrderNumber: "123456789",
})
```

## TLS RPC Client

We need to use it when we write core financial services

```go
tlsConf, err := xrpc.LoadClientTLSConfig("/certificates/ca.pem", "/certificates/client.pem", "/certificates/client.key")
if err != nil {
    log.Fatal(err)
}
conn, err := xrpc.NewGrpcClient("127.0.0.1:50000", grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)))
```

Examples of other languages

```PHP
<?php
require __DIR__ . '/vendor/autoload.php';

$opts = [
    'credentials' => Grpc\ChannelCredentials::createSsl(file_get_contents('/certificates/server.pem'), file_get_contents('/certificates/client.key'), file_get_contents('/certificates/client.pem')),
    // 'grpc.ssl_target_name_override' => '127.0.0.1:50000',
    // 'grpc.default_authority' => '127.0.0.1:50000'
];
$client = new \Bitfloww\RexClient('127.0.0.1:50000', $opts);
$request = new \Bitfloww\ReleaseRequest();
$request->setOrderNumber('123456789');
list($reply, $status) = $client->RequestForRelease($request)->wait();
var_dump($reply, $status);
```

## TLS Gateway Client

```go
tlsConf, err := xrpc.LoadClientTLSConfig("/certificates/ca.pem", "/certificates/client.pem", "/certificates/client.key")
if err != nil {
    log.Fatal(err)
}
client := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: tlsConf,
    },
}
defer client.CloseIdleConnections()
resp, err := client.Post("https://127.0.0.1:50001/v1/request_for_release", "application/json", strings.NewReader(`{"order_number":"123456789"}`))
fmt.Println(resp.Body)
```

Examples of other languages

```PHP
<?php
require __DIR__ . '/vendor/autoload.php';

use GuzzleHttp\Client;

$client = new Client([
    'cert' => '/certificates/client.pem',
    'ssl_key' => '/certificates/client.key',
    'verify' => '/certificates/ca.pem'
]);
$response = $client->request('POST', 'https://127.0.0.1:50001/v1/request_for_release', ['body' => '{"order_number":"123456789"}']);
var_dump($response->getBody()->getContents());
```

## Two-way TLS with Self-Signed Certificates

Modify subjectAltName in `generate-rsa.cnf`, the generated files are in the `certificates` directory.

```
generate-rsa.sh
```

Code new TLS Config

- Server

Load from file

```go
tlsConf, err := xrpc.LoadServerTLSConfig("/certificates/ca.pem", "/certificates/server.pem", "/certificates/server.key")
```

New by bytes

```go
tlsConf, err := xrpc.NewServerTLSConfig([]byte{}, []byte{}, []byte{})
```

- Client

Load from file

```go
tlsConf, err := xrpc.LoadClientTLSConfig("/certificates/ca.pem", "/certificates/client.pem", "/certificates/client.key")
```

New by bytes

```go
tlsConf, err := xrpc.NewClientTLSConfig([]byte{}, []byte{}, []byte{})
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
