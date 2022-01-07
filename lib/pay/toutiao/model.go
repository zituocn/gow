package toutiao

// CreateOrderResp 预下单应答参数
type CreateOrderResp struct {
	ErrNo   int64      `json:"err_no"`
	ErrTips string     `json:"err_tips"`
	Data    *OrderInfo `json:"data"`
}

type OrderInfo struct {
	OrderId    string `json:"order_id"`
	OrderToken string `json:"order_token"`
}

// QueryOrderResp 查询订单应答参数
type QueryOrderResp struct {
	ErrNo       int64               `json:"err_no"`
	ErrTips     string              `json:"err_tips"`
	OutOrderNo  string              `json:"out_order_no"`
	OrderId     string              `json:"order_id"`
	PaymentInfo *QueryOrderRespData `json:"payment_info"`
}

type QueryOrderRespData struct {
	TotalFee         int64  `json:"total_fee"`
	OrderStatus      string `json:"order_status"`       //PROCESSING-处理中|SUCCESS-成功|FAIL-失败|TIMEOUT-超时
	PayTime          string `json:"pay_time"`           //支付时间
	Way              int64  `json:"way"`                //支付渠道 2-支付宝，1-微信，3-银行卡
	ChannelNo        string `json:"channel_no"`         //渠道单号
	ChannelGatewayNo string `json:"channel_gateway_no"` //渠道网关号
}

// NotifyReq 支付回调请求参数
type NotifyReq struct {
	Timestamp    int64  `json:"timestamp"`
	Nonce        string `json:"nonce"`         //随机数
	Msg          string `json:"msg"`           //订单信息的json字符串
	Type         string `json:"type"`          //回调类型标记，支付成功回调为"payment"
	MsgSignature string `json:"msg_signature"` //签名
}

// NotifyMsgData 回调信息包括 msg 信息为以下内容序列化得到的 json 字符串
type NotifyMsgData struct {
	Appid     string `json:"appid"`
	CpOrderno string `json:"cp_orderno"` //开发者传入订单号
	Way       string `json:"way"`        //way 字段中标识了支付渠道：2-支付宝，1-微信
	CpExtra   string `json:"cp_extra"`   //预下单时开发者传入字段
}

// NotifyReturn 回调响应参数
type NotifyReturn struct {
	ErrNo  int64  `json:"err_no"`
	ErrTip string `json:"err_tips"`
}

type CreateOrderReq struct {
	AppId       string `json:"app_id"`
	Sign        string `json:"sign"`
	OutOrderNo  string `json:"out_order_no"`
	TotalAmount int64  `json:"total_amount"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	ValidTime   int64  `json:"valid_time"`
	CpExtra     string `json:"cp_extra"`
	NotifyUrl   string `json:"notify_url"`
	DisableMsg  int64  `json:"disable_msg"`
}

//申请退款应答参数
type RefundResp struct {
	ErrNo    int64  `json:"err_no"`
	ErrTips  string `json:"err_tips"`
	RefundNo string `json:"refund_no"`
}

//退款回调信息包括 msg 信息为以下内容序列化得到的 json 字符串
type RefundNotifyMsgData struct {
	Appid        string `json:"appid"`
	CpRefundno   string `json:"cp_refundno"`   //开发者自定义的退款单号
	CpExtra      string `json:"cp_extra"`      //预下单时开发者传入字段
	Status       string `json:"status"`        //退款状态 PROCESSING-处理中|SUCCESS-成功|FAIL-失败
	RefundAmount int64  `json:"refund_amount"` //退款金额
}

//查询退款应答参数
type RefundQueryResp struct {
	ErrNo      int64            `json:"err_no"`
	ErrTips    string           `json:"err_tips"`
	RefundInfo *RefundQueryData `json:"refundInfo"`
}

type RefundQueryData struct {
	RefundNo     string `json:"refund_no"`
	RefundAmount int64  `json:"refund_amount"`
	RefundStatus string `json:"refund_status"` //退款状态，成功-SUCCESS；失败-FAIL
}
