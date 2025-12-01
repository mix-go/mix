package xutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAESEncryptCBC(t *testing.T) {
	a := assert.New(t)

	plain := "Hello, playground"
	key128 := "abcdefghijklmnop"                 // 16 bytes = 128 bits
	key192 := "abcdefghijklmnopqrstuvwx"         // 24 bytes = 192 bits
	key256 := "abcdefghijklmnopabcdefghijklmnop" // 32 bytes = 256 bits
	ivInfo := "abcdefghijklmnop"

	cipherText128, _ := AESEncrypt(plain, AESModeCBC, key128, ivInfo)
	fmt.Printf("Cipher text for AES-128-CBC: %s\n", cipherText128)
	cipherText192, _ := AESEncrypt(plain, AESModeCBC, key192, ivInfo)
	fmt.Printf("Cipher text for AES-192-CBC: %s\n", cipherText192)
	cipherText256, _ := AESEncrypt(plain, AESModeCBC, key256, ivInfo)
	fmt.Printf("Cipher text for AES-256-CBC: %s\n", cipherText256)

	decryptedText128, _ := AESDecrypt(cipherText128, AESModeCBC, key128, ivInfo)
	fmt.Printf("Decrypted text for AES-128-CBC: %s\n", decryptedText128)
	decryptedText192, _ := AESDecrypt(cipherText192, AESModeCBC, key192, ivInfo)
	fmt.Printf("Decrypted text for AES-192-CBC: %s\n", decryptedText192)
	decryptedText256, _ := AESDecrypt(cipherText256, AESModeCBC, key256, ivInfo)
	fmt.Printf("Decrypted text for AES-256-CBC: %s\n", decryptedText256)

	a.Equal(plain, decryptedText128)
	a.Equal(plain, decryptedText192)
	a.Equal(plain, decryptedText256)
}

func TestAESEncryptCFB(t *testing.T) {
	a := assert.New(t)

	plain := "Hello, playground"
	key128 := "abcdefghijklmnop"                 // 16 bytes = 128 bits
	key192 := "abcdefghijklmnopqrstuvwx"         // 24 bytes = 192 bits
	key256 := "abcdefghijklmnopabcdefghijklmnop" // 32 bytes = 256 bits
	ivInfo := "abcdefghijklmnop"

	cipherText128, _ := AESEncrypt(plain, AESModeCFB, key128, ivInfo)
	fmt.Printf("Cipher text for AES-128-CBC: %s\n", cipherText128)
	cipherText192, _ := AESEncrypt(plain, AESModeCFB, key192, ivInfo)
	fmt.Printf("Cipher text for AES-192-CBC: %s\n", cipherText192)
	cipherText256, _ := AESEncrypt(plain, AESModeCFB, key256, ivInfo)
	fmt.Printf("Cipher text for AES-256-CBC: %s\n", cipherText256)

	decryptedText128, _ := AESDecrypt(cipherText128, AESModeCFB, key128, ivInfo)
	fmt.Printf("Decrypted text for AES-128-CBC: %s\n", decryptedText128)
	decryptedText192, _ := AESDecrypt(cipherText192, AESModeCFB, key192, ivInfo)
	fmt.Printf("Decrypted text for AES-192-CBC: %s\n", decryptedText192)
	decryptedText256, _ := AESDecrypt(cipherText256, AESModeCFB, key256, ivInfo)
	fmt.Printf("Decrypted text for AES-256-CBC: %s\n", decryptedText256)

	a.Equal(plain, decryptedText128)
	a.Equal(plain, decryptedText192)
	a.Equal(plain, decryptedText256)
}

func TestAESEncryptOFB(t *testing.T) {
	a := assert.New(t)

	plain := "Hello, playground"
	key128 := "abcdefghijklmnop"                 // 16 bytes = 128 bits
	key192 := "abcdefghijklmnopqrstuvwx"         // 24 bytes = 192 bits
	key256 := "abcdefghijklmnopabcdefghijklmnop" // 32 bytes = 256 bits
	ivInfo := "abcdefghijklmnop"

	cipherText128, _ := AESEncrypt(plain, AESModeOFB, key128, ivInfo)
	fmt.Printf("Cipher text for AES-128-CBC: %s\n", cipherText128)
	cipherText192, _ := AESEncrypt(plain, AESModeOFB, key192, ivInfo)
	fmt.Printf("Cipher text for AES-192-CBC: %s\n", cipherText192)
	cipherText256, _ := AESEncrypt(plain, AESModeOFB, key256, ivInfo)
	fmt.Printf("Cipher text for AES-256-CBC: %s\n", cipherText256)

	decryptedText128, _ := AESDecrypt(cipherText128, AESModeOFB, key128, ivInfo)
	fmt.Printf("Decrypted text for AES-128-CBC: %s\n", decryptedText128)
	decryptedText192, _ := AESDecrypt(cipherText192, AESModeOFB, key192, ivInfo)
	fmt.Printf("Decrypted text for AES-192-CBC: %s\n", decryptedText192)
	decryptedText256, _ := AESDecrypt(cipherText256, AESModeOFB, key256, ivInfo)
	fmt.Printf("Decrypted text for AES-256-CBC: %s\n", decryptedText256)

	a.Equal(plain, decryptedText128)
	a.Equal(plain, decryptedText192)
	a.Equal(plain, decryptedText256)
}
