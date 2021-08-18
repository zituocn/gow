/*

微信小程序/头条小程序
解密小程序敏感信息:
对称解密使用的算法为 AES-128-CBC，数据采用PKCS#7填充。
对称解密的目标密文为 Base64_Decode(encryptedData)。
对称解密秘钥 aeskey = Base64_Decode(session_key), aeskey 是16字节。
对称解密算法初始向量 为Base64_Decode(iv)，其中iv由数据接口返回。

*/

package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
)

// AppletDecrypt 解密小程序敏感信息
func AppletDecrypt(sessionKey, encryptedData string, iv string) ([]byte, error) {
	sessionKey = strings.Replace(strings.TrimSpace(sessionKey), " ", "+", -1)
	if len(sessionKey) != 24 {
		return nil, errors.New("sessionKey length is error")
	}
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, errors.New("DecodeString err")
	}
	iv = strings.Replace(strings.TrimSpace(iv), " ", "+", -1)
	if len(iv) != 24 {
		return nil, errors.New("iv length is error")
	}
	aesIv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, errors.New("DecodeString err")
	}
	encryptedData = strings.Replace(strings.TrimSpace(encryptedData), " ", "+", -1)
	aesCipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, errors.New("DecodeString err")
	}
	aesPlantText := make([]byte, len(aesCipherText))

	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, errors.New("NewCipher err")
	}
	mode := cipher.NewCBCDecrypter(aesBlock, aesIv)
	mode.CryptBlocks(aesPlantText, aesCipherText)
	aesPlantText = PKCS7UnPadding(aesPlantText)

	re := regexp.MustCompile(`[^\{]*(\{.*\})[^\}]*`)
	aesPlantText = []byte(re.ReplaceAllString(string(aesPlantText), "$1"))

	return aesPlantText, nil
}

// PKCS7UnPadding return unpadding []Byte plantText
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	if length > 0 {
		unPadding := int(plantText[length-1])
		return plantText[:(length - unPadding)]
	}
	return plantText
}
