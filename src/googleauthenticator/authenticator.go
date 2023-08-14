package googleauthenticator

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const NONE = "none"

func GenerateSecret() string {
	formattedKey := encodeGoogleAuthKey(generateOtpKey())
	rx := regexp.MustCompile(`\W+`)
	secret := []byte(rx.ReplaceAllString(formattedKey, ""))
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)
}

// Generate a key
func generateOtpKey() []byte {
	// 20 cryptographically random binary bytes (160-bit key)
	key := make([]byte, 20)
	_, _ = rand.Read(key)
	return key
}

// Text-encode the key as base32 (in the style of Google Authenticator - same as Facebook, Microsoft, etc)
func encodeGoogleAuthKey(bin []byte) string {
	// 32 ascii characters without trailing '='s
	rx := regexp.MustCompile(`=`)
	base32 := rx.ReplaceAllString(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bin), "")
	base32 = strings.ToLower(base32)

	// lowercase with a space every 4 characters
	rx = regexp.MustCompile(`(\w{4})`)
	key := strings.TrimSpace(rx.ReplaceAllString(base32, "$1 "))

	return key
}

func VerifyToken(secret, passcode string) bool {
	return VerifyTokenCustom(secret, passcode, 30)
}

func VerifyTokenCustom(secret, passcode string, period int) bool {
	b, _ := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      NONE,
		AccountName: NONE,
		SecretSize:  uint(len(b)),
		Secret:      b,
	})
	rv, _ := totp.ValidateCustom(
		passcode,
		key.Secret(),
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    uint(period),
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	return rv
}

func GenerateToken(secret string) (passcode string) {
	b, _ := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      NONE,
		AccountName: NONE,
		SecretSize:  uint(len(b)),
		Secret:      b,
	})
	passcode, _ = totp.GenerateCode(key.Secret(), time.Now().UTC())
	return
}

func GenerateTotpUri(issuer, accountName, secret string) string {
	// Full OTPAUTH URI spec as explained at
	// https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	u := url.URL{}
	v := url.Values{}
	u.Scheme = "otpauth"
	u.Host = "totp"
	u.Path = fmt.Sprintf("%s:%s", issuer, accountName)
	v.Add("secret", secret)
	v.Add("issuer", issuer)
	v.Add("algorithm", "SHA1")
	v.Add("digits", strconv.Itoa(6))
	v.Add("period", strconv.Itoa(30))
	u.RawQuery = v.Encode()
	return u.String()
}

func GenerateQRCodeGoogleUrl(issuer, accountName, secret string) string {
	uri := GenerateTotpUri(issuer, accountName, secret)
	return fmt.Sprintf("https://www.google.com/chart?chs=200x200&chld=M|0&cht=qr&chl=%s", url.QueryEscape(uri))
}
