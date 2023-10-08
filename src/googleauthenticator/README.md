> Produced by OpenMix: [https://openmix.org](https://openmix.org/mix-go)

## Mix Google Authenticator

Install

```
go get github.com/mix-go/googleauthenticator@latest
```

Generate Secret

```go
secret := googleauthenticator.GenerateSecret()
```

Generate Code

```go
code := googleauthenticator.GenerateToken(secret)
```

Verify Code

```go
ok := googleauthenticator.VerifyToken(secret, code)
// or
ok := googleauthenticator.VerifyTokenCustom(secret, code, 60)
```

Generate Url

```go
uri := googleauthenticator.GenerateTotpUri("Foo", "bar", secret)
// or
url := googleauthenticator.GenerateQRCodeGoogleUrl("Foo", "bar", secret)
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
