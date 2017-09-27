package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

func makeSha1(s []byte) []byte {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return bs
}

func GenerateSignature(key []byte, iv []byte, message []byte) string {
	sha1_message := makeSha1(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	content := PKCS5Padding(sha1_message, block.BlockSize())
	crypted := make([]byte, len(content))
	cbc.CryptBlocks(crypted, content)
	return base64.StdEncoding.EncodeToString(crypted)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
