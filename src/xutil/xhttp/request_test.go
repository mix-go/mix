package xhttp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequest(t *testing.T) {
	a := assert.New(t)

	url := "https://github.com/"
	resp, err := Request("GET", url)

	a.Equal(resp.StatusCode, 200)
	a.Nil(err)
}

func TestRequestPOST(t *testing.T) {
	a := assert.New(t)

	url := "https://github.com/"
	resp, err := Request("POST", url, WithBodyString("abc"), WithContentType("application/json"))

	a.Equal(resp.StatusCode, 404)
	a.Nil(err)
}

func TestRequestError(t *testing.T) {
	a := assert.New(t)

	url := "https://aaaaa.com/"
	resp, err := Request("GET", url)

	a.Nil(resp)
	a.NotNil(err)
}
