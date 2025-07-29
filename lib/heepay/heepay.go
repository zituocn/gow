package heepay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/zituocn/gow/lib/util"
	"github.com/zituocn/logx"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	WX_MINI_TRADE_URL = "https://pay.heepay.com/Phone/SDK/PayInit.aspx"
	ALI_WAP_TRADE_URL = "https://pay.heepay.com/Payment/Index.aspx"
	ORDER_QUERY_URL   = "https://query.heepay.com/Payment/Query.aspx"
	REFUND_URL        = "https://pay.heepay.com/API/Payment/PaymentRefund.aspx"
	REFUND_QUERY_URL  = "https://Query.heepay.com/API/Payment/PaymentRefundQuery.aspx"
	SIGN_TYPE         = "md5"
	IS_PHONE          = "1"
	VERSION           = "1"
	RETURN_MODE       = "1"
)

type Heepay struct {
	AgentID   string
	PayKey    string
	RefundKey string
}

func NewHeepay(agentID, payKey, refundKey string) *Heepay {
	return &Heepay{
		AgentID:   agentID,
		PayKey:    payKey,
		RefundKey: refundKey,
	}
}

// WxMiniTrade 微信小程序下单
func (h *Heepay) WxMiniTrade(order *PaymentOrder) (tokenID string, err error) {
	err = h.ValidatePayAmt(order.PayAmt)
	if err != nil {
		return
	}

	signStr := fmt.Sprintf("version=%s&agent_id=%s&agent_bill_id=%s&agent_bill_time=%s&pay_type=%d&pay_amt=%s&notify_url=%s&user_ip=%s&key=%s",
		VERSION, h.AgentID, order.AgentBillID, order.AgentBillTime, order.PayType, order.PayAmt.String(), order.NotifyURL, order.UserIP, h.PayKey)
	signStr = strings.ReplaceAll(signStr, " ", "")
	sign := util.MD5(signStr)

	encoder := simplifiedchinese.GBK.NewEncoder()
	gbkReader := transform.NewReader(
		bytes.NewReader([]byte(order.MetaOption)),
		encoder,
	)
	gbkBytes, err := io.ReadAll(gbkReader)
	if err != nil {
		return
	}
	metaOption := base64.StdEncoding.EncodeToString(gbkBytes)

	params := url.Values{}
	params.Add("version", VERSION)
	params.Add("agent_id", h.AgentID)
	params.Add("agent_bill_id", order.AgentBillID)
	params.Add("agent_bill_time", order.AgentBillTime)
	params.Add("pay_type", strconv.Itoa(order.PayType))
	params.Add("pay_amt", order.PayAmt.String())
	params.Add("notify_url", order.NotifyURL)
	params.Add("return_url", order.ReturnURL)
	params.Add("user_ip", order.UserIP)
	params.Add("goods_name", order.GoodsName)
	params.Add("goods_num", strconv.Itoa(order.GoodsNum))
	params.Add("goods_note", order.GoodsNote)
	params.Add("remark", order.Remark)
	params.Add("meta_option", metaOption)
	params.Add("sign", sign)

	fullURL := WX_MINI_TRADE_URL + "?" + params.Encode()
	response, err := http.Get(fullURL)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	tokenID, err = h.parseResponse(body)
	return
}

func (h *Heepay) parseResponse(xmlData []byte) (string, error) {
	var token TokenResponse
	err := xml.Unmarshal(xmlData, &token)
	if err == nil && token.Value != "" {
		return token.Value, nil
	}

	var errorResp ErrorResponse
	err = xml.Unmarshal(xmlData, &errorResp)
	if err == nil && errorResp.Message != "" {
		return "", fmt.Errorf("parseResponse: 接口错误: %s", errorResp.Message)
	}

	return "", fmt.Errorf("parseResponse: 无法解析的XML格式: %s", string(xmlData))
}

// AliWapTrade payUrl: 前端访问的支付地址
func (h *Heepay) AliWapTrade(order *PaymentOrder) (payUrl string, err error) {
	err = h.ValidatePayAmt(order.PayAmt)
	if err != nil {
		return
	}

	signStr := fmt.Sprintf("version=%s&agent_id=%s&agent_bill_id=%s&agent_bill_time=%s&pay_type=%d&pay_amt=%s&notify_url=%s&return_url=%s&user_ip=%s&key=%s",
		VERSION, h.AgentID, order.AgentBillID, order.AgentBillTime, order.PayType, order.PayAmt.String(), order.NotifyURL, order.ReturnURL, order.UserIP, h.PayKey)
	if order.TimeStamp != 0 {
		signStr += fmt.Sprintf("&timestamp=%d", order.TimeStamp)
	}
	signStr = strings.ReplaceAll(signStr, " ", "")
	sign := util.MD5(signStr)

	var metaOption string
	if order.MetaOption != "" {
		encoder := simplifiedchinese.GBK.NewEncoder()
		gbkReader := transform.NewReader(
			bytes.NewReader([]byte(order.MetaOption)),
			encoder,
		)
		gbkBytes, err := io.ReadAll(gbkReader)
		if err != nil {
			return "", err
		}
		metaOption = base64.StdEncoding.EncodeToString(gbkBytes)
	}

	params := url.Values{}
	params.Add("version", VERSION)
	params.Add("is_phone", IS_PHONE)
	params.Add("agent_id", h.AgentID)
	params.Add("agent_bill_id", order.AgentBillID)
	params.Add("agent_bill_time", order.AgentBillTime)
	params.Add("pay_type", strconv.Itoa(order.PayType))
	params.Add("pay_amt", order.PayAmt.String())
	params.Add("notify_url", order.NotifyURL)
	params.Add("return_url", order.ReturnURL)
	params.Add("user_ip", order.UserIP)
	if order.GoodsName != "" {
		params.Add("goods_name", order.GoodsName)
	}
	if order.GoodsNum > 0 {
		params.Add("goods_num", strconv.Itoa(order.GoodsNum))
	}
	if order.GoodsNote != "" {
		params.Add("goods_note", order.GoodsNote)
	}
	if order.Remark != "" {
		params.Add("remark", order.Remark)
	}
	if order.TimeStamp > 0 {
		params.Add("timestamp", fmt.Sprintf("%d", order.TimeStamp))
	}
	if metaOption != "" {
		params.Add("meta_option", metaOption)
	}
	params.Add("expire_time", fmt.Sprintf("%d", order.ExpireTime))
	params.Add("sign_type", SIGN_TYPE)
	params.Add("sign", sign)

	payUrl = ALI_WAP_TRADE_URL + "?" + params.Encode()
	return
}

// ValidatePayAmt 验证金额
func (h *Heepay) ValidatePayAmt(payAmt decimal.Decimal) error {
	validate := validator.New()

	// 检查是否为两位小数
	validate.RegisterValidation("twoDecimalPlaces", func(fl validator.FieldLevel) bool {
		str := fl.Field().String()
		parts := strings.Split(str, ".")
		if len(parts) == 2 && len(parts[1]) > 2 {
			return false
		}
		return true
	})

	if err := validate.Struct(payAmt); err != nil {
		return err
	}

	// 检查小数点后是否超过两位
	if payAmt.Shift(2).Mod(decimal.NewFromInt(1)).GreaterThan(decimal.Zero) {
		return errors.New("金额必须精确到小数点后两位")
	}

	// 检查范围
	min := decimal.NewFromFloat(0.01)
	max := decimal.NewFromFloat(10000000.00)
	if payAmt.LessThan(min) || payAmt.GreaterThan(max) {
		return errors.New("金额必须在0.01到10000000.00之间")
	}

	return nil
}

func (h *Heepay) PayNotify(req *http.Request) (ret string, payDetail *PayNotifyResp, err error) {
	result := req.URL.Query().Get("result")
	payMessage := req.URL.Query().Get("pay_message")
	agentID := req.URL.Query().Get("agent_id")
	jnetBillNO := req.URL.Query().Get("jnet_bill_no")
	agentBillID := req.URL.Query().Get("agent_bill_id")
	payTypeStr := req.URL.Query().Get("pay_type")
	payAmtStr := req.URL.Query().Get("pay_amt")
	remark := req.URL.Query().Get("remark")
	payUser := req.URL.Query().Get("pay_user")
	tradeBillNO := req.URL.Query().Get("trade_bill_no")
	bankCardType := req.URL.Query().Get("bank_card_type")
	bankCardOwnerType := req.URL.Query().Get("bank_card_owner_type")
	dealTime := req.URL.Query().Get("deal_time")
	paramSign := req.URL.Query().Get("sign")
	ret = "error"
	verifySignStr := fmt.Sprintf("result=%s&agent_id=%s&jnet_bill_no=%s&agent_bill_id=%s&pay_type=%s&pay_amt=%s&remark=%s&key=%s",
		result, agentID, jnetBillNO, agentBillID, payTypeStr, payAmtStr, remark, h.PayKey)
	verifySignStr = strings.ReplaceAll(verifySignStr, " ", "")
	verifySign := util.MD5(verifySignStr)
	if paramSign != verifySign {
		err = errors.New("PayNotify: 签名验证错误")
		return
	}

	ret = "ok"
	payAmt, err := decimal.NewFromString(payAmtStr)
	if err != nil {
		return
	}

	payType, err := strconv.Atoi(payTypeStr)
	if err != nil {
		return
	}
	payDetail = &PayNotifyResp{
		Result:            result,
		PayMessage:        payMessage,
		AgentID:           agentID,
		JnetBillNO:        jnetBillNO,
		AgentBillID:       agentBillID,
		PayType:           payType,
		PayAmt:            payAmt,
		Remark:            remark,
		PayUser:           payUser,
		TradeBillNO:       tradeBillNO,
		BankCardType:      bankCardType,
		BankCardOwnerType: bankCardOwnerType,
		DealTime:          dealTime,
		Sign:              verifySign,
	}
	return
}

// OrderQuery 订单收款查询 agentBillID:商户系统内部的订单号 agentBillTime:提交单据的时间yyyyMMddHHmmss，该参数共计14位，当时不满14位时，在后面加0补足14位 remark:商户自定义，原样返回,可以为空。
func (h *Heepay) OrderQuery(agentBillID, agentBillTime, remark string) (queryInfo OrderQueryInfo, err error) {
	signStr := fmt.Sprintf("version=%s&agent_id=%s&agent_bill_id=%s&agent_bill_time=%s&return_mode=%s&key=%s",
		VERSION, h.AgentID, agentBillID, agentBillTime, RETURN_MODE, h.PayKey)
	signStr = strings.ReplaceAll(signStr, " ", "")
	sign := util.MD5(signStr)

	params := url.Values{}
	params.Add("version", VERSION)
	params.Add("agent_id", h.AgentID)
	params.Add("agent_bill_id", agentBillID)
	params.Add("agent_bill_time", agentBillTime)
	params.Add("remark", remark)
	params.Add("return_mode", RETURN_MODE)
	params.Add("sign", sign)

	fullURL := ORDER_QUERY_URL + "?" + params.Encode()
	response, err := http.Get(fullURL)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	queryInfo, err = h.parseQueryInfo(body)
	return
}

func (h *Heepay) parseQueryInfo(info []byte) (queryInfo OrderQueryInfo, err error) {
	data := make(map[string]interface{})
	pairs := strings.Split(string(info), "|")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		data[kv[0]] = kv[1]

		if kv[0] == "pay_type" {
			var value int
			value, err = strconv.Atoi(kv[1])
			if err != nil {
				return
			}
			data[kv[0]] = value
		}
		if kv[0] == "pay_amt" {
			var value decimal.Decimal
			value, err = decimal.NewFromString(kv[1])
			if err != nil {
				return
			}
			data[kv[0]] = value
		}

	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonData, &queryInfo)
	return
}

// PaymentRefund 退款
// agent_bill_id和refund_details为互斥参数，必须传其中一个参数。
// 默认使用agent_bill_id，当agent_bill_id为空使用refund_details。
// 单笔全额退款使用agent_bill_id参数，批量部分退款使用refund_details参数
// refund_details格式: 商户原支付单号,金额,商户退款单号; 金额传0或为空默认做全额退款，商户退款单号可为空，多条数据中间用竖线隔开，例：63548281250,0.01,5232112|6358281251,0,
func (h *Heepay) PaymentRefund(agentBillID, refundDetails, notifyURL string) (refundInfo RefundXMLResp, err error) {
	if agentBillID == "" && refundDetails == "" || agentBillID != "" && refundDetails != "" {
		err = errors.New("PaymentRefund: agent_bill_id refund_details参数错误")
		return
	}

	var signStr string
	params := url.Values{}
	if agentBillID != "" {
		params.Add("agent_bill_id", agentBillID)
		signStr = fmt.Sprintf("agent_bill_id=%s&agent_id=%s&key=%s&notify_url=%s&version=%s", agentBillID, h.AgentID, h.RefundKey, notifyURL, VERSION)
	} else if refundDetails != "" {
		params.Add("refund_details", refundDetails)
		signStr = fmt.Sprintf("agent_id=%s&key=%s&notify_url=%s&refund_details=%s&version=%s", h.AgentID, h.RefundKey, notifyURL, refundDetails, VERSION)
	}
	signStr = strings.ReplaceAll(signStr, " ", "")
	signStr = strings.ToLower(signStr)
	sign := util.MD5(signStr)
	params.Add("version", VERSION)
	params.Add("agent_id", h.AgentID)
	params.Add("notify_url", notifyURL)
	params.Add("sign", sign)
	params.Add("sign_type", SIGN_TYPE)

	fullURL := REFUND_URL + "?" + params.Encode()
	response, err := http.Get(fullURL)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	reader, err := charset.NewReader(bytes.NewReader(body), response.Header.Get("Content-Type"))
	if err != nil {
		logx.Errorf("PaymentRefund: 编码检测错误: %v\n", err)
		reader = bytes.NewReader(body) // 无法检测时使用原始内容
	}

	// 解析XML
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&refundInfo)
	if err != nil {
		logx.Errorf("PaymentRefund: XML解析错误: %v\n", err)
	}

	return
}

func (h *Heepay) RefundNotify(req *http.Request) (ret string, refundDetail *RefundNotifyResp, err error) {
	agentID := req.URL.Query().Get("agent_id")
	hyBillNO := req.URL.Query().Get("hy_bill_no")
	agentBillID := req.URL.Query().Get("agent_bill_id")
	agentRefundBillNO := req.URL.Query().Get("agent_refund_bill_no")
	refundAmt := req.URL.Query().Get("refund_amt")
	refundStatus := strings.ToUpper(req.URL.Query().Get("refund_status"))
	hyDealtime := req.URL.Query().Get("hy_deal_time")
	paramSign := req.URL.Query().Get("sign")
	ret = "error"
	verifySignStr := fmt.Sprintf("agent_id=%s&hy_bill_no=%s&agent_bill_id=%s&agent_refund_bill_no=%s&refund_amt=%s&refund_status=%s&hy_deal_time=%s&key=%s",
		agentID, hyBillNO, agentBillID, agentRefundBillNO, refundAmt, refundStatus, hyDealtime, h.RefundKey)
	verifySignStr = strings.ReplaceAll(verifySignStr, " ", "")
	verifySignStr = strings.ToLower(verifySignStr)
	verifySign := util.MD5(verifySignStr)
	if paramSign != verifySign {
		err = errors.New("RefundNotify: 签名验证错误")
		return
	}

	ret = "ok"
	refundDetail = &RefundNotifyResp{
		AgentID:           agentID,
		HyBillNO:          hyBillNO,
		AgentBillID:       agentBillID,
		AgentRefundBillNO: agentRefundBillNO,
		RefundAmt:         refundAmt,
		RefundStatus:      refundStatus,
		HyDealTime:        hyDealtime,
		Sign:              verifySign,
	}
	return
}

// PaymentRefundQuery 退款查询,agent_bill_id: 商户系统内部订单号,必填；agent_refund_bill_no: 商户系统内部退款单号，可为空，非必填
func (h *Heepay) PaymentRefundQuery(agentBillID, agentRefundBillNO string) (refundQueryInfo RefundQueryXMLResp, err error) {
	signStr := fmt.Sprintf("agent_bill_id=%s&agent_id=%s&key=%s&version=%s", agentBillID, h.AgentID, h.RefundKey, VERSION)
	signStr = strings.ReplaceAll(signStr, " ", "")
	signStr = strings.ToLower(signStr)
	sign := util.MD5(signStr)

	params := url.Values{}
	params.Add("version", VERSION)
	params.Add("agent_id", h.AgentID)
	params.Add("agent_bill_id", agentBillID)
	params.Add("sign", sign)
	params.Add("agent_refund_bill_no", agentRefundBillNO)

	fullURL := REFUND_QUERY_URL + "?" + params.Encode()
	response, err := http.Get(fullURL)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	reader, err := charset.NewReader(bytes.NewReader(body), response.Header.Get("Content-Type"))
	if err != nil {
		logx.Errorf("PaymentRefundQuery: 编码检测错误: %v\n", err)
		reader = bytes.NewReader(body) // 无法检测时使用原始内容
	}

	// 解析XML
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&refundQueryInfo)
	if err != nil {
		logx.Errorf("PaymentRefundQuery: XML解析错误: %v\n", err)
	}
	return
}
