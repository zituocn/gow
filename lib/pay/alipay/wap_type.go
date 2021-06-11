package alipay

import (
	"encoding/json"
)

// AliPayTradeWapPay 手机网站支付
// doc: https://docs.open.alipay.com/api_1/alipay.trade.wap.pay/
type AliPayTradeWapPay struct {
	TradePay
	QuitURL    string `json:"quit_url,omitempty"`
	AuthToken  string `json:"auth_token,omitempty"`  // 针对用户授权接口，获取用户相关数据时，用于标识用户授权关系
	TimeExpire string `json:"time_expire,omitempty"` // 绝对超时时间，格式为yyyy-MM-dd HH:mm。
}

//APIName APIName
func (t *AliPayTradeWapPay) APIName() string {
	return "alipay.trade.wap.pay"
}

//Params Params
func (t *AliPayTradeWapPay) Params() map[string]string {
	var m = make(map[string]string)
	m["notify_url"] = t.NotifyURL
	m["return_url"] = t.ReturnURL
	return m
}

//ExtJSONParamName ExtJSONParamName
func (t *AliPayTradeWapPay) ExtJSONParamName() string {
	return "biz_content"
}

//ExtJSONParamValue ExtJSONParamValue
func (t *AliPayTradeWapPay) ExtJSONParamValue() string {
	var bytes, err = json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(bytes)
}
