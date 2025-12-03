## Mix XUtil

A set of tools that keep Golang sweet.

## Installation

```
go get github.com/mix-go/xutil
```

## Functions

| Function                                                                                                                                 | Description                                                              |  
|------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|
| xutil.StartAlerter(ntype NotifierType, credential, titlePrefix string, rateLimitInterval time.Duration, mqSize int, logger Logger) error | Launch an Alerter.                                                       |
| xutil.PushAlert(ctx context.Context, title, content, uuid string, mentionAll bool) error                                                 | Use Alerter to push alerts via a queue, and limit frequency.             |
| xutil.SendAlert(ctx context.Context, title, content string, mentionAll bool) error                                                       | Use Alerter to send alerts directly.                                     |
| xutil.ErrorAlert(err error, needStack bool)                                                                                              | Use Alerter to push error alerts via a queue, and limit frequency.       |
| xutil.StartGrpcMonitoring(interval time.Duration, logger Logger)                                                                         | Launch a gRPC monitor.                                                   |
| xutil.StatsServerOptions(logger Logger) []grpc.ServerOption                                                                              | Track the daily request count for the FullMethod of the gRPC server.     |
| xutil.StartPerformanceMonitoring(interval time.Duration, logger Logger, handler func(*PerformanceStats))                                 | Launch a Go performance monitor.                                         |
| xutil.SubString(s string, start int, length int) string                                                                                  | Return part of a string.                                                 |
| xutil.Capitalize(s string) string                                                                                                        | The function converts the first letter of the input string to uppercase. |
| xutil.IsNumeric(s string) bool                                                                                                           | Used to check if the variable is a number or a numeric string.           |
| xutil.StructToMap(i interface{}) map[string]interface{}                                                                                  | Convert struct to map.                                                   |
| xutil.StringToBytes(s string) []byte                                                                                                     | Convert string to bytes (0 copy).                                        |
| xutil.BytesToString(b []byte) string                                                                                                     | Convert bytes to bytes (0 copy).                                         |
| xutil.AESEncrypt(plainText, mode, key, iv string) (string, error)                                                                        | AES encryption.                                                          |
| xutil.AESDecrypt(cipherText, mode, key, iv string) (string, error)                                                                       | AES decryption.                                                          |
| xutil.MD5Hash(b []byte) string                                                                                                           | MD5 encryption.                                                          |

## License

Apache License Version 2.0, http://www.apache.org/licenses/
