package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var key []byte

func SetPasswors(pass string) {
	hasher := sha256.New()
	hasher.Write([]byte(pass))
	key = hasher.Sum(nil)
}

func Hmac(s string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("empty key")
	}
	h := hmac.New(sha256.New, key)
	h.Write([]byte(s))
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash), nil
}

func Crypt(s string) (string, error) {
	plainText := []byte(s)
	crypt, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(crypt, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.RawStdEncoding.EncodeToString(cipherText), err
}

func UnCrypt(s string) (string, error) {
	cipherText, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("ciphertext block size is too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
