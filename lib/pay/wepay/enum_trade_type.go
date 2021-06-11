package wepay

//TradeType 支付场景类型
type TradeType string

const (
	// TradeTypeJSAPI 公众号支付
	TradeTypeJSAPI TradeType = "JSAPI"

	// TradeTypeNative 扫码支付
	TradeTypeNative TradeType = "NATIVE"

	// TradeTypeApp App支付
	TradeTypeApp TradeType = "APP"

	//TradeTypeH5 非微信浏览器的H5唤起支付
	TradeTypeH5 TradeType = "MWEB"

	//TradeTypeApplet 小程序
	TradeTypeApplet TradeType = ""
)
