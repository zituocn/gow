package wepay

import "encoding/xml"

//AppPayResp App支付时，返回的结构体
//包括prepayID和Sign等其他信息
type AppPayResp struct {
	AppID     string `json:"appid"`     //appid
	PartnerID string `json:"partnerid"` //商户ID
	PrepayID  string `json:"prepayid"`  //prepayid
	Package   string `json:"package"`   //package
	NonceStr  string `json:"noncestr"`  //随机字串
	Timestamp string `json:"timestamp"` //时间
	Sign      string `json:"sign"`      //签名
}

//NotifyRet 异步通知的返回值
//返回
type NotifyRet struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

//AppletPayResp 微信小程序支付时，返回的结构体
//包括prepayID和Sign等其他信息
type AppletPayResp struct {
	Timestamp string `json:"timeStamp"` //时间
	NonceStr  string `json:"nonceStr"`  //随机字串
	Package   string `json:"package"`   //package
	SignType  string `json:"signType"`
	Sign      string `json:"paySign"` //签名
}
