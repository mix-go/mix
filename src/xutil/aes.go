package xutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

const (
	AESModeCBC = "CBC"
	AESModeCFB = "CFB"
	AESModeOFB = "OFB"
)

var ErrUnknownMode = errors.New("unknown mode")

// AESEncrypt AES encryption, supports CBC,CFB,OFB mode, PKCS7Padding padding.
// key128 := "abcdefghijklmnop"                 // 16 bytes = 128 bits
// key192 := "abcdefghijklmnopqrstuvwx"         // 24 bytes = 192 bits
// key256 := "abcdefghijklmnopabcdefghijklmnop" // 32 bytes = 256 bits
func AESEncrypt(plainText, mode, key, iv string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plainTextBytes := PKCS7Padding([]byte(plainText), block.BlockSize())
	cipherText := make([]byte, len(plainTextBytes))

	switch mode {
	case AESModeCBC:
		cbc := cipher.NewCBCEncrypter(block, []byte(iv))
		cbc.CryptBlocks(cipherText, plainTextBytes)
	case AESModeCFB:
		cfb := cipher.NewCFBEncrypter(block, []byte(iv))
		cfb.XORKeyStream(cipherText, plainTextBytes)
	case AESModeOFB:
		ofb := cipher.NewOFB(block, []byte(iv))
		ofb.XORKeyStream(cipherText, plainTextBytes)
	default:
		return "", ErrUnknownMode
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func AESDecrypt(cipherText, mode, key, iv string) (string, error) {
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	decryptedText := make([]byte, len(cipherTextBytes))

	switch mode {
	case AESModeCBC:
		cbc := cipher.NewCBCDecrypter(block, []byte(iv))
		cbc.CryptBlocks(decryptedText, cipherTextBytes)
	case AESModeCFB:
		cfb := cipher.NewCFBDecrypter(block, []byte(iv))
		cfb.XORKeyStream(decryptedText, cipherTextBytes)
	case AESModeOFB:
		ofb := cipher.NewOFB(block, []byte(iv))
		ofb.XORKeyStream(decryptedText, cipherTextBytes)
	default:
		return "", ErrUnknownMode
	}

	decryptedText = PKCS7UnPadding(decryptedText)

	return string(decryptedText), nil
}

// PKCS7Padding adds padding to the input plainTextByte according to PKCS7 mode
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding removes padding from the input cipherTextBytes according to PKCS7 mode
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	// subtract the number of padding bytes
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
