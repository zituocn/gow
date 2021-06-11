package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

const (
	//iv 值
	//16位长度
	iv = "0000000000000000"
)

// AESEncrypt AES CBC encrypt
func AESEncrypt(plainText string, key string) string {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlainText := PKCS5Padding([]byte(plainText), aes.BlockSize, len(plainText))
	block, _ := aes.NewCipher(bKey)
	cipherText := make([]byte, len(bPlainText))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(cipherText, bPlainText)
	//注意：此处使用hex编码
	return hex.EncodeToString(cipherText)
}

// AESDecrypt 解码
func AESDecrypt(decodeStr string, key string) (string, error) {
	bKey := []byte(key)
	bIV := []byte(iv)
	//注意：此处使用使用hex解码
	decodeBytes, err := hex.DecodeString(decodeStr)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, bIV)
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData = PKCS5UnPadding(origData)
	return string(origData), nil
}

// PKCS5Padding PKCS5Padding
func PKCS5Padding(cipherText []byte, blockSize int, after int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// PKCS5UnPadding PKCS5UnPadding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
