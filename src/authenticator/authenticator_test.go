package authenticator_test

import (
	"fmt"
	"github.com/mix-go/authenticator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyToken(t *testing.T) {
	a := assert.New(t)

	secret := authenticator.GenerateSecret()
	code := authenticator.GenerateToken(secret)
	ok := authenticator.VerifyToken(secret, code)
	uri := authenticator.GenerateTotpUri("foo", "bar", secret)
	url := authenticator.GenerateQRCodeGoogleUrl("foo", "bar", secret)
	fmt.Printf("%v\n%s\n%s\n", ok, uri, url)

	a.Equal(ok, true)
}

func TestPHPVerifyToken(t *testing.T) {
	a := assert.New(t)

	secret := "OQB6ZZGYHCPSX4AK"
	code := authenticator.GenerateToken(secret)
	ok := authenticator.VerifyToken(secret, code)
	uri := authenticator.GenerateTotpUri("foo", "bar", secret)
	url := authenticator.GenerateQRCodeGoogleUrl("foo", "bar", secret)
	fmt.Printf("%v\n%s\n%s\n", ok, uri, url)

	a.Equal(ok, true)
}
