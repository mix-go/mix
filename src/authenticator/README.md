> Developed by OpenMix: [https://openmix.org](https://openmix.org/mix-go)

## Mix Authenticator

Install

```
go get github.com/mix-go/authenticator@latest
```

Generate Secret

```go
secret := authenticator.GenerateSecret()
```

Generate Code

```go
code := authenticator.GenerateToken(secret)
```

Verify Code

```go
ok := authenticator.VerifyToken(secret, code)
// or
ok := authenticator.VerifyTokenCustom(secret, code, 60)
```

Generate Url

```go
uri := authenticator.GenerateTotpUri("Foo", "bar", secret)
// or
url := authenticator.GenerateQRCodeGoogleUrl("Foo", "bar", secret)
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
