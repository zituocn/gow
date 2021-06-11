package wepay

import (
	"crypto/tls"
	"encoding/pem"
	"github.com/clbanning/mxj"
	"golang.org/x/crypto/pkcs12"
	"log"
	"math/rand"
	"strings"
)

var (
	//RandChar 随机字串
	RandChar = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	//RandCharLen 随机字串长度
	RandCharLen int32 = 36
	//TradeNumStr 商户订单号前缀
	TradeNumStr = "WX"
)

//XMLToMap XMLToMap
// 使用mxj第三方库，转换xml2map
func XMLToMap(xmlStr string) Params {
	mv, _ := mxj.NewMapXml([]byte(xmlStr))
	params, _ := mv["xml"].(map[string]interface{})
	// if !ok {
	// 	fmt.Println("error:", ok)
	// }
	return params
}

// MapToXML MapToXML
// 使用mxj第三方库，转换map2xml
func MapToXML(params Params) string {
	mv := mxj.Map(params)
	data, _ := mv.Xml()
	return string(data)
}

// makeNonceStr 生成随机字符串
func makeNonceStr(n int) string {
	sb := new(strings.Builder)
	for i := 0; i < n; i++ {
		sb.WriteByte(RandChar[rand.Int31n(RandCharLen)])
	}
	return sb.String()
}

// 将Pkcs12转成Pem
func pkcs12ToPem(p12 []byte, password string) tls.Certificate {

	blocks, err := pkcs12.ToPEM(p12, password)

	// 从恐慌恢复
	defer func() {
		if x := recover(); x != nil {
			log.Print(x)
		}
	}()

	if err != nil {
		panic(err)
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		panic(err)
	}
	return cert
}
