package googleauthenticator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyToken(t *testing.T) {
	a := assert.New(t)

	secret := GenerateSecret()
	code := GenerateToken(secret)
	ok := VerifyToken(secret, code)
	uri := GenerateTotpUri("foo", "bar", secret)
	url := GenerateQRCodeGoogleUrl("foo", "bar", secret)
	fmt.Printf("%v\n%s\n%s\n", ok, uri, url)

	a.Equal(ok, true)
}

func TestPHPVerifyToken(t *testing.T) {
	a := assert.New(t)

	secret := "OQB6ZZGYHCPSX4AK"
	code := GenerateToken(secret)
	ok := VerifyToken(secret, code)
	uri := GenerateTotpUri("foo", "bar", secret)
	url := GenerateQRCodeGoogleUrl("foo", "bar", secret)
	fmt.Printf("%v\n%s\n%s\n", ok, uri, url)

	a.Equal(ok, true)
}
