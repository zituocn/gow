package wepay

import "io/ioutil"

// WxConfig 配置参数
type WxConfig struct {
	AppId           string //传入的appID
	MchId           string //分配的mchID
	APIKey          string //分配的apiKey
	ServerIP        string //服务器IP
	certData        []byte //证书
	isSandbox       bool   //是否沙箱
	NotifyURL       string //异步通知地址
	RefundNotifyUrl string //订单退款异步通知地址
	OrderTime       int    //订单有效分钟数
	IsProfitSharing bool   //是否分账 默认为不分账
}

// NewWxConfig 一个新的配置信息
// 也可以自己组装
func NewWxConfig(appId, mchId, apiKey, serverIP string, isSandbox bool, notifyUrl string, orderTime int) *WxConfig {
	return &WxConfig{
		AppId:     appId,
		MchId:     mchId,
		APIKey:    apiKey,
		ServerIP:  serverIP,
		isSandbox: isSandbox,
		NotifyURL: notifyUrl,
		OrderTime: orderTime,
	}
}

// SetCertData 设置证书
func (m *WxConfig) SetCertData(certPath string) (err error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return
	}
	m.certData = certData
	return
}

// SetRefundNotifyUrl 设置退款异步通知地址
func (m *WxConfig) SetRefundNotifyUrl(url string) {
	m.RefundNotifyUrl = url
}

func (m *WxConfig) SetIsProfitSharing(needProfitSharing bool) {
	if needProfitSharing {
		m.IsProfitSharing = true
	}
}
