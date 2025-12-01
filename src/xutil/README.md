## Mix XUtil

A set of tools that keep Golang sweet.

## Installation

```
go get github.com/mix-go/xutil
```

| Function                                                           | Description                                                              |  
|--------------------------------------------------------------------|--------------------------------------------------------------------------|
| xutil.IsNumeric(s string) bool                                     | Used to check if the variable is a number or a numeric string.           |
| xutil.SubString(s string, start int, length int) string            | Return part of a string                                                  |
| xutil.Capitalize(s string) string                                  | The function converts the first letter of the input string to uppercase. |
| xutil.StructToMap(i interface{}) map[string]interface{}            | Convert struct to map.                                                   |
| xutil.StringToBytes(s string) []byte                               | Convert string to bytes (0 copy).                                        |
| xutil.BytesToString(b []byte) string                               | Convert bytes to bytes (0 copy).                                         |
| xutil.AESEncrypt(plainText, mode, key, iv string) (string, error)  | AES encryption                                                           |
| xutil.AESDecrypt(cipherText, mode, key, iv string) (string, error) | AES Decryption                                                           |

## License

Apache License Version 2.0, http://www.apache.org/licenses/
