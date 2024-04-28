package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

const (
	AESKey    = "OSC@GiteeCodeAES"
	AESPrefix = "{AES}"
)

func Encrypt(key []byte, text string) (string, error) {
	if len(key) == 0 {
		key = []byte(AESKey)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(key []byte, text string) (string, error) {
	if len(key) == 0 {
		key = []byte(AESKey)
	}
	ciphertext, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func GetPlainPassword(key []byte, text string) string {
	if !strings.HasPrefix(text, AESPrefix) {
		return text
	}
	p, err := Decrypt(key, text)
	if err != nil {
		fmt.Printf("decrypt password fail\n")
		return text
	} else {
		return p
	}
}
