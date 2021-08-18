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

// AESDecrypt AES CBC decrypt
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

// AesDecryptECB AES-256-ECB decrypt
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}

/*
private
*/

func generateKey(key []byte) (genKey []byte) {
	if len(key) == 16 || len(key) == 24 || len(key) == 32 {
		return key
	}
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}
