package wepay

import "encoding/xml"

// AppPayResp App支付时，返回的结构体
// 包括prepayID和Sign等其他信息
type AppPayResp struct {
	AppID     string `json:"appid"`     //appid
	PartnerID string `json:"partnerid"` //商户ID
	PrepayID  string `json:"prepayid"`  //prepayid
	Package   string `json:"package"`   //package
	NonceStr  string `json:"noncestr"`  //随机字串
	Timestamp string `json:"timestamp"` //时间
	Sign      string `json:"sign"`      //签名
}

// NotifyRet 异步通知的返回值
// 返回
type NotifyRet struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

// AppletPayResp 微信小程序支付时，返回的结构体
// 包括prepayID和Sign等其他信息
type AppletPayResp struct {
	Timestamp string `json:"timeStamp"` //时间
	NonceStr  string `json:"nonceStr"`  //随机字串
	Package   string `json:"package"`   //package
	SignType  string `json:"signType"`
	Sign      string `json:"paySign"` //签名
}

// 微信退款成功通知 返回的信息
type WXPayRefundSuccessNotifyResp struct {
	TransactionId     string `json:"transaction_id"`
	OutTradeNo        string `json:"out_trade_no"`
	RefundId          string `json:"refund_id"`
	OutRefundNo       string `json:"out_refund_no"`
	RefundStatus      string `json:"refund_status"`
	RefundRecvAccount string `json:"refund_recv_account"`
	RefundFee         int64  `json:"refund_fee"`   //退款金额 单位为分
	SuccessTime       string `json:"success_time"` //资金退款至用户账号的时间，格式2017-12-15 09:46:01
}

// ProfitSharingReceiverReq 请求分账接收方信息
type ProfitSharingReceiverReq struct {
	Type        string `json:"type"`           //分账接收方类型-必传
	Account     string `json:"account"`        //分账接收方账号-必传
	Amount      int64  `json:"amount"`         //分账金额-必传，单位分
	Description string `json:"description"`    //分账描述-必传
	Name        string `json:"name,omitempty"` //分账个人接收方姓名-选填
}

// ProfitSharingReceiverData 请求分账相应接收方信息
type ProfitSharingReceiverData struct {
	Type    string `json:"type"`    //分账接收方类型
	Account string `json:"account"` //接收方账号
	//开发文档（https://pay.weixin.qq.com/wiki/doc/api/allocation.php?chapter=27_1&index=1）的字段是receiver_account和receiver_mchid，经测试，返回的字段是account且没有receiver_mchid
	//ReceiverAccount string `json:"receiver_account"`
	//ReceiveMchId string `json:"receiver_mchid"`
	Amount      int64  `json:"amount"`      //分账金额
	Description string `json:"description"` //分账描述
	FailReason  string `json:"fail_reason"` //分账失败原因
	/*
		分账失败原因，当分账结果result为CLOSED（已关闭）时，返回该字段
		枚举值：
		1、ACCOUNT_ABNORMAL : 分账接收账户异常
		2、NO_RELATION : 分账关系已解除
		3、RECEIVER_HIGH_RISK : 高风险接收方
		4、RECEIVER_REAL_NAME_NOT_VERIFIED : 接收方未实名
		5、NO_AUTH : 分账权限已解除
		6、RECEIVER_RECEIPT_LIMIT : 接收方已达收款限额
		7、PAYER_ACCOUNT_ABNORMAL : 分出方账户异常
	*/
	DetailId   string `json:"detail_id"`   //分账明细单号
	FinishTime string `json:"finish_time"` //完成时间
	Result     string `json:"result"`      //分账结果：PENDING：待分账，SUCCESS：分账成功，CLOSED：分账失败已关闭
}

// ProfitSharingQueryReceiverData 查询分账结果 接收方信息
type ProfitSharingQueryReceiverData struct {
	Type        string `json:"type"`        //分账接收方类型
	Account     string `json:"account"`     //接收方账号-【查询分账结果】接口返回
	Amount      int64  `json:"amount"`      //分账金额
	Description string `json:"description"` //分账描述
	FailReason  string `json:"fail_reason"` //分账失败原因
	/*
		分账失败原因，当分账结果result为CLOSED（已关闭）时，返回该字段
		枚举值：
		1、ACCOUNT_ABNORMAL : 分账接收账户异常
		2、NO_RELATION : 分账关系已解除
		3、RECEIVER_HIGH_RISK : 高风险接收方
		4、RECEIVER_REAL_NAME_NOT_VERIFIED : 接收方未实名
		5、NO_AUTH : 分账权限已解除
		6、RECEIVER_RECEIPT_LIMIT : 超出用户月收款限额
		7、PAYER_ACCOUNT_ABNORMAL : 分出方账户异常
		8、INVALID_REQUEST: 描述参数设置失败
	*/
	DetailId   string `json:"detail_id"`   //分账明细单号
	FinishTime string `json:"finish_time"` //完成时间
	Result     string `json:"result"`      //分账结果：PENDING：待分账，SUCCESS：分账成功，CLOSED：分账失败已关闭
}

type ProfitSharingResp struct {
	TransactionId string                       `json:"transaction_id"`
	OutOrderNo    string                       `json:"out_order_no"`
	OrderId       string                       `json:"order_id"`
	Status        string                       `json:"status"` //PROCESSING：处理中  FINISHED：处理完成
	Receivers     []*ProfitSharingReceiverData `json:"receivers"`
}

type ProfitSharingQueryResp struct {
	TransactionId string                            `json:"transaction_id"`
	OutOrderNo    string                            `json:"out_order_no"`
	OrderId       string                            `json:"order_id"`
	Status        string                            `json:"status"`
	Receivers     []*ProfitSharingQueryReceiverData `json:"receivers"`
}

type ProfitSharingReturnRet struct {
	OrderId           string `json:"order_id"`
	OutOrderNo        string `json:"out_order_no"`
	OutReturnNo       string `json:"out_return_no"`
	ReturnNo          string `json:"return_no"`
	ReturnAccountType string `json:"return_account_type"`
	ReturnAccount     string `json:"return_account"`
	ReturnAmount      int64  `json:"return_amount"`
	Description       string `json:"description"`
	Result            string `json:"result"`
	/*
		回退结果：
		PROCESSING:处理中
		SUCCESS:已成功
		FAILED: 已失败
		如果返回为处理中，请勿变更商户回退单号，使用相同的参数再次发起分账回退，否则会出现资金风险
		在处理中状态的回退单如果5天没有成功，会因为超时被设置为已失败
	*/
	FailReason string `json:"fail_reason"`
	FinishTime string `json:"finish_time"`
}
