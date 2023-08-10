package xcrypt

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

	cipherText128, _ := AESEncrypt(plain, ModeCBC, key128, ivInfo)
	fmt.Printf("Cipher text for AES-128-CBC: %s\n", cipherText128)
	cipherText192, _ := AESEncrypt(plain, ModeCBC, key192, ivInfo)
	fmt.Printf("Cipher text for AES-192-CBC: %s\n", cipherText192)
	cipherText256, _ := AESEncrypt(plain, ModeCBC, key256, ivInfo)
	fmt.Printf("Cipher text for AES-256-CBC: %s\n", cipherText256)

	decryptedText128, _ := AESDecrypt(cipherText128, ModeCBC, key128, ivInfo)
	fmt.Printf("Decrypted text for AES-128-CBC: %s\n", decryptedText128)
	decryptedText192, _ := AESDecrypt(cipherText192, ModeCBC, key192, ivInfo)
	fmt.Printf("Decrypted text for AES-192-CBC: %s\n", decryptedText192)
	decryptedText256, _ := AESDecrypt(cipherText256, ModeCBC, key256, ivInfo)
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

	cipherText128, _ := AESEncrypt(plain, ModeCFB, key128, ivInfo)
	fmt.Printf("Cipher text for AES-128-CBC: %s\n", cipherText128)
	cipherText192, _ := AESEncrypt(plain, ModeCFB, key192, ivInfo)
	fmt.Printf("Cipher text for AES-192-CBC: %s\n", cipherText192)
	cipherText256, _ := AESEncrypt(plain, ModeCFB, key256, ivInfo)
	fmt.Printf("Cipher text for AES-256-CBC: %s\n", cipherText256)

	decryptedText128, _ := AESDecrypt(cipherText128, ModeCFB, key128, ivInfo)
	fmt.Printf("Decrypted text for AES-128-CBC: %s\n", decryptedText128)
	decryptedText192, _ := AESDecrypt(cipherText192, ModeCFB, key192, ivInfo)
	fmt.Printf("Decrypted text for AES-192-CBC: %s\n", decryptedText192)
	decryptedText256, _ := AESDecrypt(cipherText256, ModeCFB, key256, ivInfo)
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

	cipherText128, _ := AESEncrypt(plain, ModeOFB, key128, ivInfo)
	fmt.Printf("Cipher text for AES-128-CBC: %s\n", cipherText128)
	cipherText192, _ := AESEncrypt(plain, ModeOFB, key192, ivInfo)
	fmt.Printf("Cipher text for AES-192-CBC: %s\n", cipherText192)
	cipherText256, _ := AESEncrypt(plain, ModeOFB, key256, ivInfo)
	fmt.Printf("Cipher text for AES-256-CBC: %s\n", cipherText256)

	decryptedText128, _ := AESDecrypt(cipherText128, ModeOFB, key128, ivInfo)
	fmt.Printf("Decrypted text for AES-128-CBC: %s\n", decryptedText128)
	decryptedText192, _ := AESDecrypt(cipherText192, ModeOFB, key192, ivInfo)
	fmt.Printf("Decrypted text for AES-192-CBC: %s\n", decryptedText192)
	decryptedText256, _ := AESDecrypt(cipherText256, ModeOFB, key256, ivInfo)
	fmt.Printf("Decrypted text for AES-256-CBC: %s\n", decryptedText256)

	a.Equal(plain, decryptedText128)
	a.Equal(plain, decryptedText192)
	a.Equal(plain, decryptedText256)
}
