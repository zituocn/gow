package heepay

import (
	"encoding/xml"

	"github.com/shopspring/decimal"
)

type PaymentOrder struct {
	AgentBillID   string          // 必填 商户系统内部的订单号,长度最长50字符
	AgentBillTime string          // 必填 提交单据的时间yyyyMMddHHmmss,该参数共计14位,当时不满14位时,在后面加0补足14位
	PayType       int             // 必填 支付类型 微信：30 支付宝：22
	PayAmt        decimal.Decimal // 必填 订单总金额 不可为空,取值范围(0.01到10000000.00),单位:元,小数点后保留两位
	NotifyURL     string          // 必填 异步通知地址
	ReturnURL     string          // 必填 同步通知地址
	UserIP        string          // 必填 用户所在客户端的真实ip,其中的"."替换为"_"
	GoodsName     string          // 必填 商品名称,不能为空
	Remark        string          // 必填 商户自定义,原样返回,可以为空
	GoodsNum      int             // 选填 产品数量
	GoodsNote     string          // 选填 支付说明
	TimeStamp     int64           // 选填 时间戳,北京时间1970/1/1 0点到现在的毫秒值,订单在+-1min内有效,超过时间订单不能提交,如果传此参数，此参数也需要参与签名，参数加在key后面
	ExpireTime    int             // 选填 订单过期相对时间，单位分钟，最低1分钟，最高4320分钟(服务器收到请求开始计时)
	MetaOption    string          // 微信小程序:必填,{"s":"微信小程序","n":"","id":""} n:网站名称或者产品名称 id:网址,参与签名;支付宝wap:选填,{"is_guarantee":"1"} is_guarantee=1代表分润单
}

// TokenResponse 成功响应结构
type TokenResponse struct {
	XMLName xml.Name `xml:"token_id"`
	Value   string   `xml:",chardata"`
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	XMLName xml.Name `xml:"error"`
	Message string   `xml:",chardata"`
}

type PayNotifyResp struct {
	Result            string          // 支付结果, 1=成功 其它为未知
	PayMessage        string          // 支付结果信息,支付成功时为空
	AgentID           string          // 商户编号
	JnetBillNO        string          // 商户退款单号(商户退款时没有退款单号则为空)
	AgentBillID       string          // 商户系统内部的订单号
	PayType           int             // 支付类型 微信：30 支付宝：22
	PayAmt            decimal.Decimal // 订单实际支付金额
	Remark            string          // 商户自定义,原样返回。
	PayUser           string          // 支付人信息。不返回值说明此支付类型不支持
	TradeBillNO       string          // 上游通道支付单号
	BankCardType      string          // 银行卡类型(0=储蓄卡 1=信用卡 -1=未知)
	BankCardOwnerType string          // 银行卡所有者类型(-1=未知 0=个人账户 1=公司账户)
	DealTime          string          // 订单支付成功时间
	Sign              string          // 签名结果
}

type OrderQueryInfo struct {
	AgentID            string          `json:"agent_id"`             // 商户编号 必填
	AgentBillID        string          `json:"agent_bill_id"`        // 商户系统内部的订单号 必填
	JnetBillNO         string          `json:"jnet_bill_no"`         // 汇付宝交易号(订单号) 必填
	PayType            int             `json:"pay_type"`             // 支付类型 微信：30 支付宝：22 必填
	Result             string          `json:"result"`               // 支付结果 1=成功、0=处理中、-1=失败 必填
	PayAmt             decimal.Decimal `json:"pay_amt"`              // 订单实际支付金额(注意:此金额是用户的实付金额) 必填
	PayMessage         string          `json:"pay_message"`          // 交易结果描述,为空时说明订单正常 必填
	Remark             string          `json:"remark"`               // 商户自定义,原样返回。请求接口传入的值 必填
	Sign               string          `json:"sign"`                 // 签名结果 必填
	DealTime           string          `json:"deal_time"`            // 订单支付时间
	PayUser            string          `json:"pay_user"`             // 支付人信息。不返回值说明此支付类型不支持
	ThirdBillNO        string          `json:"third_bill_no"`        // 汇付宝提交给上游的订单号
	TradeBillNO        string          `json:"trade_bill_no"`        // 上游通道支付单号
	BankCardType       string          `json:"bank_card_type"`       // 银行卡类型(0=储蓄卡 1=信用卡 -1=未知)
	BankCardOwnerType  string          `json:"bank_card_owner_type"` // 银行卡所有者类型(-1=未知 0=个人账户 1=公司账户)
	DetailErrorCode    string          `json:"detail_error_code"`    // 失败状态码。E105 银行限制错误,E113 用户卡信息异常,E114 余额不足,E115 用户信息异常,E116 协议签约状态异常,E111 订单不存在
	DetailErrorMessage string          `json:"detail_error_message"` // detail_error_code返回值时,返回具体报错信息
	RefAgentId         string          `json:"ref_agent_id"`         // 二级商户号(集团商户模式返回)
	GuaranteeType      string          `json:"guarantee_type"`       // 订单类型(0=单笔结算、一般订单；1=担保分润 、担保支付 、单笔分润结算。)
}

type RefundXMLResp struct {
	XMLName xml.Name `xml:"root"`
	RetCode string   `xml:"ret_code"` // 0000表示操作成功
	RetMsg  string   `xml:"ret_msg"`
	AgentID string   `xml:"agent_id,omitempty"`
	Sign    string   `xml:"sign"`
}

// 退款通知 返回的信息
type RefundNotifyResp struct {
	AgentID           string // 商户编号
	HyBillNO          string // 汇付宝交易号(订单号)
	AgentBillID       string // 商户支付单据号
	AgentRefundBillNO string // 商户退款单号(商户退款时没有退款单号则为空)
	RefundAmt         string // 退款金额
	RefundStatus      string // 退款状态(SUCCESS=成功,FAIL=失败,REFUNDING=退款中)
	HyDealTime        string // 退款时间(格式:yyyyMMddHHmmss)
	Sign              string // 签名结果
}

type RefundQueryXMLResp struct {
	XMLName    xml.Name `xml:"root"`
	RetCode    string   `xml:"ret_code"` // 0000表示查询成功
	RetMsg     string   `xml:"ret_msg"`
	AgentID    string   `xml:"agent_id"`
	DetailData string   `xml:"detail_data"` // 为-1001时表示无此退款单据;每条明细后跟换行符分隔;以英文逗号分隔的字符串: 汇付宝退款单号,商户退款单号(可为空),退款金额,退款状态(-1=失败；0=处理中,1=退款中,2=退款成功),退款时间,记录内码,商户手续费,子商户手续费(手续费是单据成功时,才会返回)
	Sign       string   `xml:"sign"`
}
