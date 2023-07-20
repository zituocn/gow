/*
供调用的封装程度更高的func
sam
参见测试代码
*/

package wepay

import (
	"errors"
	"fmt"
	"github.com/zituocn/gow/lib/util"
	"github.com/zituocn/logx"

	"net/http"
	"strings"
	"time"
)

// WxAPI 调用实体
type WxAPI struct {
	Client *Client
}

// NewWxAPI init
// account基本的商户信息
// notifyURL 异步通知地址
// endMinute 订单有效期
func NewWxAPI(wxConfig *WxConfig) *WxAPI {
	return &WxAPI{
		Client: &Client{
			WxConfig:           wxConfig, //配置信息
			signType:           MD5,      //验证方式
			httpConnectTimeout: 2000,     //连接时间
			httpReadTimeout:    1000,     //读时间
		},
	}
}

// UnifiedOrder 统一下单接口
func (m *WxAPI) UnifiedOrder(body string, outTradeNo string, totalFee int, openID string, clientIP string, tradeType TradeType) (Params, error) {
	params := make(Params)
	params.SetString("body", body)
	params.SetString("out_trade_no", outTradeNo)
	params.SetInt64("total_fee", int64(totalFee))
	params.SetString("openid", openID)
	params.SetString("trade_type", string(tradeType))
	params.SetString("notify_url", m.Client.NotifyURL)

	//订单的有效期，开始和结束时间
	now := time.Now()
	params.SetString("time_start", util.TimeFormat(now, "YYYYMMDDHHmmss"))
	params.SetString("time_expire", util.TimeFormat(now.Add(time.Minute*time.Duration(m.Client.OrderTime)), "YYYYMMDDHHmmss"))

	if clientIP != "" {
		params.SetString("spbill_create_ip", clientIP)
	} else if tradeType == TradeTypeNative {
		params.SetString("spbill_create_ip", m.Client.ServerIP) //服务器IP
	}
	return m.Client.UnifiedOrder(params)
}

// AppTrade APP下单
func (m *WxAPI) AppTrade(body, outTradeNo string, totalFee int, clientIP string) (ret *AppPayResp, err error) {
	params, err := m.UnifiedOrder(body, outTradeNo, totalFee, "", clientIP, TradeTypeApp)
	if err != nil {
		return
	}

	prepayID := strings.TrimSpace(params.GetString("prepay_id"))
	if len(prepayID) == 0 {
		err = fmt.Errorf("返回prepay_id失败")
		return
	}

	//时间戳
	timestamp := time.Now().Unix()
	//随机值
	nonceStr := makeNonceStr(20)
	pg := "Sign=WXPay"
	p := make(Params)
	p.SetString("appid", m.Client.AppId)
	p.SetString("partnerid", m.Client.MchId)
	p.SetString("prepayid", prepayID)
	p.SetString("package", pg)
	p.SetString("noncestr", nonceStr)
	p.SetString("timestamp", fmt.Sprintf("%d", timestamp))
	//计算并返回签名
	sign := m.Client.Sign(p)

	ret = &AppPayResp{
		AppID:     m.Client.AppId,
		PartnerID: m.Client.MchId,
		PrepayID:  prepayID,
		Package:   pg,
		NonceStr:  nonceStr,
		Timestamp: fmt.Sprintf("%d", timestamp),
		Sign:      sign,
	}
	return
}

// NativeTrade 扫码支付下单
// 使用时把返回的code_url生成二维码，供前台用户扫码支付
func (m *WxAPI) NativeTrade(body, outTradeNo string, totalFee int) (ret string, err error) {
	params, err := m.UnifiedOrder(body, outTradeNo, totalFee, "", "", TradeTypeNative)
	if err != nil {
		return
	}
	ret = params.GetString("code_url")
	if len(ret) == 0 {
		err = fmt.Errorf("返回的code_url为空")
		return
	}
	return
}

// H5Trade H5下单
func (m *WxAPI) H5Trade(body, outTradeNo string, totalFee int, clientIP string) (ret string, err error) {
	params, err := m.UnifiedOrder(body, outTradeNo, totalFee, "", clientIP, TradeTypeH5)
	if err != nil {
		return
	}
	ret = params.GetString("mweb_url")
	if len(ret) == 0 {
		err = fmt.Errorf("返回的mweb_url为空")
		return
	}
	return
}

// JSAPITrade 公众号支付
// 需要传入公众号对应用户的openID
func (m *WxAPI) JSAPITrade(body, outTradeNo string, totalFee int, openID, clientIP string) (ret string, err error) {
	params, err := m.UnifiedOrder(body, outTradeNo, totalFee, openID, clientIP, TradeTypeJSAPI)
	if err != nil {
		return
	}

	errCode := params.GetString("err_code")
	if errCode == "PARAM_ERROR" {
		err = fmt.Errorf(params.GetString("err_code_desc"))
		return
	}

	ret = strings.TrimSpace(params.GetString("prepay_id"))
	if len(ret) == 0 {
		err = fmt.Errorf("返回prepay_id失败")
		return
	}
	return
}

// AppletTrade 小程序支付
// 需要传入公众号对应用户的openID
func (m *WxAPI) AppletTrade(body, outTradeNo string, totalFee int, openID, clientIP string) (ret *AppletPayResp, err error) {
	params, err := m.UnifiedOrder(body, outTradeNo, totalFee, openID, clientIP, TradeTypeJSAPI)
	if err != nil {
		return
	}

	errCode := params.GetString("err_code")
	if errCode == "PARAM_ERROR" {
		err = fmt.Errorf(params.GetString("err_code_desc"))
		return
	}

	prepayID := strings.TrimSpace(params.GetString("prepay_id"))
	if len(prepayID) == 0 {
		err = fmt.Errorf("返回prepay_id失败")
		return
	}

	//时间戳
	timestamp := time.Now().Unix()
	//随机值
	nonceStr := makeNonceStr(20)
	pg := fmt.Sprintf("prepay_id=%v", prepayID)
	p := make(Params)
	p.SetString("appId", m.Client.AppId)
	p.SetString("package", pg)
	p.SetString("nonceStr", nonceStr)
	p.SetString("timeStamp", fmt.Sprintf("%d", timestamp))
	p.SetString("signType", m.Client.signType)
	//计算并返回签名
	sign := m.Client.Sign(p)
	ret = &AppletPayResp{
		Timestamp: fmt.Sprintf("%d", timestamp),
		NonceStr:  nonceStr,
		Package:   pg,
		SignType:  m.Client.signType,
		Sign:      sign,
	}
	return
}

// Notify 异步通知
// 返回异步通知状态信息
// 调用方拿到返回值后，需要根据 outTradeNo tradeNo openID等值，做进一点检验，如果检验失败，设置ret.ReturnCode="FAIL"；
// 成功时，需要回写返回值到本地
// 最后以xml方式输出 ret
func (m *WxAPI) Notify(req *http.Request) (ret *NotifyRet, outTradeNo, tradeNo, openID string, err error) {
	params, err := m.Client.Notify(req)
	if err != nil {
		return
	}

	ret = new(NotifyRet)
	if params["return_code"] == "SUCCESS" {
		ret.ReturnCode = "SUCCESS"
		ret.ReturnMsg = "OK"
	}

	//返回给调用方做校验证和回写使用；
	outTradeNo = params.GetString("out_trade_no")
	tradeNo = params.GetString("transaction_id")
	openID = params.GetString("openid")

	return
}

// OrderQuery 订单查询
//
//	返回是否成功，和错误信息
func (m *WxAPI) OrderQuery(transactionID, outTradeNo string) (state bool, tradeNo string, err error) {
	if transactionID == "" && outTradeNo == "" {
		err = fmt.Errorf("[transactionID]与[outTradeNo]不能同时为空")
		return
	}
	params := make(Params)
	params.SetString("transaction_id", transactionID)
	params.SetString("out_trade_no", outTradeNo)
	params, err = m.Client.OrderQuery(params)
	if err != nil {
		return
	}

	//向前台隐藏几个关键信息
	delete(params, "appid")
	delete(params, "mch_id")
	delete(params, "sign")

	//返回有可能存在的tradeNo
	tradeNo = params.GetString("transaction_id")

	if params.GetString("return_code") == "SUCCESS" && params.GetString("trade_state") == "SUCCESS" {
		state = true
	} else if params.GetString("trade_state") == "NOTPAY" { //未支付
		err = fmt.Errorf(params.GetString("trade_state_desc"))
	} else if params.GetString("trade_state") == "CLOSED" { //订单已关闭
		err = fmt.Errorf(params.GetString("trade_state_desc"))
	} else {
		err = fmt.Errorf(params.GetString("err_code_des")) //其他错误
	}

	return
}

/*
退款相关SDK
2021/8/18
*/

// Refund 申请退款
func (m *WxAPI) Refund(transactionID, outTradeNo, outRefundNo, refundDesc string, totalFee, refundFee int) (refundId string, err error) {
	if outTradeNo == "" && transactionID == "" {
		err = fmt.Errorf("[outTradeNo]和[transactionID]不能同时为空")
		return
	}
	if outRefundNo == "" {
		err = fmt.Errorf("[outRefundNo]不能为空")
		return
	}
	params := make(Params)
	params.SetString("transaction_id", transactionID)
	params.SetString("out_trade_no", outTradeNo)
	params.SetString("out_refund_no", outRefundNo)           //商户退款单号
	params.SetInt64("total_fee", int64(totalFee))            //订单金额
	params.SetInt64("refund_fee", int64(refundFee))          //退款金额
	params.SetString("refund_desc", refundDesc)              //退款原因
	params.SetString("notify_url", m.Client.RefundNotifyUrl) //退款结果通知url
	params, err = m.Client.Refund(params)
	if err != nil {
		return
	}
	if params.GetString("return_code") == "SUCCESS" {
		refundId = params.GetString("refund_id")
		if params.GetString("err_code") != "" {
			logx.Errorf(fmt.Sprintf("申请退款失败：错误码：%v,错误描述：%v", params.GetString("err_code"), params.GetString("err_code_des")))
			err = errors.New(params.GetString("err_code_des"))
		}
	} else {
		logx.Errorf("申请退款提交业务失败：%v", params.GetString("return_msg"))
	}
	return
}

// RefundQuery 退款查询
func (m *WxAPI) RefundQuery(transactionID, outTradeNo, outRefundNo, refundId string) (status string, recvName string, err error) {
	if transactionID == "" && outTradeNo == "" && outRefundNo == "" && refundId == "" {
		err = fmt.Errorf("[transactionID]和[outTradeNo]和[outRefundNo]和[refundId]不能同时为空")
		return
	}
	params := make(Params)
	params.SetString("transaction_id", transactionID)
	params.SetString("out_trade_no", outTradeNo)
	params.SetString("out_refund_no", outRefundNo) //商户退款单号
	params.SetString("refund_id", refundId)        //微信生成的退款单号
	params, err = m.Client.RefundQuery(params)
	if err != nil {
		return
	}
	if params.GetString("return_code") == "SUCCESS" {
		//校验退款单号是否匹配
		var tag bool
		if transactionID != "" && transactionID == params.GetString("transaction_id") {
			tag = true
		}
		if outTradeNo != "" && outTradeNo == params.GetString("out_trade_no") {
			tag = true
		}
		//_0 这类字段都是加了偏移量offset，由于此处的业务没有一个订单下多笔退款单的情况，都将offset默认为0
		if outRefundNo != "" && outTradeNo == params.GetString("out_refund_no_0") {
			tag = true
		}
		if refundId != "" && refundId == params.GetString("refund_id_0") {
			tag = true
		}

		logy.Debugf("refund_id_0:%v", params.GetString("refund_id_0"))
		logy.Debugf("refund_id_1:%v", params.GetString("refund_id_1"))
		logy.Debugf("refund_id_2:%v", params.GetString("refund_id_2"))

		if tag {
			//退款状态
			/*
				SUCCESS—退款成功
				REFUNDCLOSE—退款关闭，指商户发起退款失败的情况
				PROCESSING—退款处理中
				CHANGE—退款异常，退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，可前往商户平台（pay.weixin.qq.com）-交易中心，手动处理此笔退款。$n为下标，从0开始编号。
			*/
			status = params.GetString("refund_status_0")
			recvName = params.GetString("refund_recv_accout_0") //退款入账账户

			if status == "REFUNDCLOSE" {
				logx.Errorf("退款已关闭，商户发起退款失败,错误码：%v，错误描述：%v", params.GetString("err_code"), params.GetString("err_code_des"))
			} else if status == "CHANGE" {
				logx.Errorf("退款异常，退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，错误码：%v，错误描述：%v", params.GetString("err_code"), params.GetString("err_code_des"))
			}
		} else {
			logx.Errorf("单号不匹配")
		}
	} else {
		logx.Errorf("查询退款业务失败：%v", params.GetString("return_msg"))
	}
	return
}

// RefundNotify 退款异步通知
//
//	返回异步通知状态信息
//	以xml方式输出 ret
func (m *WxAPI) RefundNotify(req *http.Request) (ret *NotifyRet, refundData *WXPayRefundSuccessNotifyResp, err error) {
	params, err := m.Client.RefundNotify(req)
	if err != nil {
		return
	}
	ret = new(NotifyRet)
	if params["return_code"] == "SUCCESS" {
		ret.ReturnCode = "SUCCESS"
		ret.ReturnMsg = "OK"
	}
	//返回给调用方做校验证和回写使用；
	//recvName = params.GetString("refund_recv_accout")
	//refundStatus = params.GetString("refund_status")
	//outTradeNo = params.GetString("out_trade_no")
	//outRefundNo = params.GetString("out_refund_no")
	refundData = &WXPayRefundSuccessNotifyResp{
		TransactionId:     params.GetString("transaction_id"),
		OutTradeNo:        params.GetString("out_trade_no"),
		RefundId:          params.GetString("refund_id"),
		OutRefundNo:       params.GetString("out_refund_no"),
		RefundStatus:      params.GetString("refund_status"),
		RefundRecvAccount: params.GetString("refund_recv_accout"),
		RefundFee:         params.GetInt64("refund_fee"),
		SuccessTime:       params.GetString("success_time"),
	}
	return
}
