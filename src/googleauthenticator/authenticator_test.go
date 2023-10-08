package googleauthenticator_test

import (
	"fmt"
	"github.com/mix-go/googleauthenticator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyToken(t *testing.T) {
	a := assert.New(t)

	secret := googleauthenticator.GenerateSecret()
	code := googleauthenticator.GenerateToken(secret)
	ok := googleauthenticator.VerifyToken(secret, code)
	uri := googleauthenticator.GenerateTotpUri("foo", "bar", secret)
	url := googleauthenticator.GenerateQRCodeGoogleUrl("foo", "bar", secret)
	fmt.Printf("%v\n%s\n%s\n", ok, uri, url)

	a.Equal(ok, true)
}

func TestPHPVerifyToken(t *testing.T) {
	a := assert.New(t)

	secret := "OQB6ZZGYHCPSX4AK"
	code := googleauthenticator.GenerateToken(secret)
	ok := googleauthenticator.VerifyToken(secret, code)
	uri := googleauthenticator.GenerateTotpUri("foo", "bar", secret)
	url := googleauthenticator.GenerateQRCodeGoogleUrl("foo", "bar", secret)
	fmt.Printf("%v\n%s\n%s\n", ok, uri, url)

	a.Equal(ok, true)
}
