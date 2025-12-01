package xutil

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hash(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
