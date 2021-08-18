/*
基础通讯业务实现
sam

*/

package wepay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/zituocn/gow/lib/logy"
	"github.com/zituocn/gow/lib/util"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// Client 请求与返回的封装
type Client struct {
	*WxConfig
	signType           string
	httpConnectTimeout int
	httpReadTimeout    int
}

// SetHTTPConnectTimeoutMs ...
func (c *Client) SetHTTPConnectTimeoutMs(ms int) {
	c.httpConnectTimeout = ms
}

// SetHTTPReadTimeoutMs ...
func (c *Client) SetHTTPReadTimeoutMs(ms int) {
	c.httpReadTimeout = ms
}

//SetSignType SetSignType
func (c *Client) SetSignType(signType string) {
	c.signType = signType
}

// fillRequestData 向 params 中添加 appid、mch_id、nonce_str、sign_type、sign
func (c *Client) fillRequestData(params Params) Params {
	params["appid"] = c.AppId
	params["mch_id"] = c.MchId
	params["nonce_str"] = makeNonceStr(20)
	params["sign_type"] = c.signType
	params["sign"] = c.Sign(params)
	return params
}

// postWithoutCert https no cert post
func (c *Client) postWithoutCert(url string, params Params) (string, error) {
	h := &http.Client{}
	p := c.fillRequestData(params)
	response, err := h.Post(url, bodyType, strings.NewReader(MapToXML(p)))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// postWithCert https need cert post
func (c *Client) postWithCert(url string, params Params) (string, error) {
	if c.certData == nil {
		return "", errors.New("证书数据为空")
	}

	// 将pkcs12证书转成pem
	cert := pkcs12ToPem(c.certData, c.MchId)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	transport := &http.Transport{
		TLSClientConfig:    config,
		DisableCompression: true,
	}
	h := &http.Client{Transport: transport}
	p := c.fillRequestData(params)
	response, err := h.Post(url, bodyType, strings.NewReader(MapToXML(p)))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// generateSignedXML 生成带有签名的xml字符串
func (c *Client) generateSignedXML(params Params) string {
	sign := c.Sign(params)
	params.SetString(Sign, sign)
	return MapToXML(params)
}

// ValidSign 验证签名
func (c *Client) ValidSign(params Params) bool {
	if !params.ContainsKey(Sign) {
		return false
	}
	return params.GetString(Sign) == c.Sign(params)
}

// Sign 签名
func (c *Client) Sign(params Params) string {
	// 创建切片
	var keys = make([]string, 0, len(params))
	// 遍历签名参数
	for k := range params {
		if k != "sign" { // 排除sign字段
			keys = append(keys, k)
		}
	}

	// 由于切片的元素顺序是不固定，所以这里强制给切片元素加个顺序
	sort.Strings(keys)

	//创建字符缓冲
	var buf bytes.Buffer
	for _, k := range keys {
		if len(params.GetString(k)) > 0 {
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(params.GetString(k))
			buf.WriteString(`&`)
		}
	}
	// 加入apiKey作加密密钥
	buf.WriteString(`key=`)
	buf.WriteString(c.APIKey)

	var (
		dataMd5    []byte
		dataSha256 []byte
		str        string
	)

	//fmt.Println("签名前的Params:", buf.String())

	switch c.signType {
	case MD5:
		h := md5.New()
		h.Write(buf.Bytes())
		dataMd5 = h.Sum(nil)
		str = hex.EncodeToString(dataMd5)
	case HMACSHA256:
		h := hmac.New(sha256.New, []byte(c.APIKey))
		h.Write(buf.Bytes())
		dataSha256 = h.Sum(nil)
		str = hex.EncodeToString(dataSha256[:])
	}

	return strings.ToUpper(str)
}

// processResponseXML 处理 HTTPS API返回数据，转换成Map对象。return_code为SUCCESS时，验证签名。
func (c *Client) processResponseXML(xmlStr string, needCheckSign ...bool) (Params, error) {
	var returnCode string
	params := XMLToMap(xmlStr, "")
	fmt.Println("返回的参数:", params)
	if params.ContainsKey("return_code") {
		returnCode = params.GetString("return_code")
	} else {
		return nil, errors.New("返回的xml中不存在[return_code]")
	}
	if returnCode == Fail {
		return nil, fmt.Errorf("返回[Fail]:%v", params.GetString("return_msg"))
	} else if returnCode == Success {
		//默认需要验证签名
		if len(needCheckSign) == 0 || (len(needCheckSign) > 0 && needCheckSign[0]) {
			if c.ValidSign(params) { //验证签名
				return params, nil
			} else {
				return nil, errors.New("返回的xml信息签名错误")
			}
		} else {
			return params, nil
		}
	} else {
		return nil, errors.New("返回的[return_code]无法识别类型")
	}
}

//Notify 异步通知处理
func (c *Client) Notify(req *http.Request) (Params, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	// 写回 body 内容
	req.Body = ioutil.NopCloser(bytes.NewReader(data))
	return c.processResponseXML(string(data))
}

// RefundNotify 退款异步通知处理
func (c *Client) RefundNotify(req *http.Request) (Params, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	params, err := c.processResponseXML(string(data), false)
	if err != nil {
		return nil, err
	}
	//加密数据
	/*
		解密步骤如下：
		（1）对加密串A做base64解码，得到加密串B
		（2）对商户key做md5，得到32位小写key* (key设置路径：微信商户平台(pay.weixin.qq.com)-->账户设置-->API安全-->密钥设置)
		（3）用key*对加密串B做AES-256-ECB解密（PKCS7Padding）
	*/
	reqInfo := params.GetString("req_info")
	if reqInfo == "" {
		err = errors.New("获得加密信息失败")
		return nil, err
	}
	//对加密串做base64解码
	decodeInfo, err := base64.StdEncoding.DecodeString(reqInfo)
	if err != nil {
		logy.Errorf(fmt.Sprintf("对加密串做base64解码出错:%v", err))
	}

	//对商户key做md5
	key := util.MD5(c.APIKey)

	//AES-256-ECB解密
	decryptData := util.AesDecryptECB(decodeInfo, []byte(key))

	/*
		解密出来的xml示例：
		<root>
		<out_refund_no><![CDATA[WX****************]]></out_refund_no>
		<out_trade_no><![CDATA[WX****************]]></out_trade_no>
		<refund_account><![CDATA[********]]></refund_account>
		<refund_fee><![CDATA[1]]></refund_fee>
		<refund_id><![CDATA[********]]></refund_id>
		<refund_recv_accout><![CDATA[支付用户零钱]]></refund_recv_accout>
		<refund_request_source><![CDATA[API]]></refund_request_source>
		<refund_status><![CDATA[SUCCESS]]></refund_status>
		<settlement_refund_fee><![CDATA[1]]></settlement_refund_fee>
		<settlement_total_fee><![CDATA[1]]></settlement_total_fee>
		<success_time><![CDATA[2021-08-12 17:49:54]]></success_time>
		<total_fee><![CDATA[1]]></total_fee>
		<transaction_id><![CDATA[*********************]]></transaction_id>
		</root>
	*/

	//将xml转换成map：根据解密出来的xml可以看出，xml解析成map的key应该是root，XMLToMap方法默认的key是xml
	retParams := XMLToMap(string(decryptData), "root")
	retParams.SetString("return_code", params.GetString("return_code"))

	// 写回 body 内容
	req.Body = ioutil.NopCloser(bytes.NewReader(data))
	return retParams, nil
}

// UnifiedOrder 统一下单
func (c *Client) UnifiedOrder(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxUnifiedOrderUrl
	} else {
		url = UnifiedOrderUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(xmlStr)
}

// MicroPay 刷卡支付
func (c *Client) MicroPay(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxMicroPayUrl
	} else {
		url = MicroPayUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(xmlStr)
}

// Refund 退款
func (c *Client) Refund(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxRefundUrl
	} else {
		url = RefundUrl
	}
	xmlStr, err := c.postWithCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(xmlStr)
}

// OrderQuery 订单查询
func (c *Client) OrderQuery(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxOrderQueryUrl
	} else {
		url = OrderQueryUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(xmlStr)
}

// RefundQuery 退款查询
func (c *Client) RefundQuery(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxRefundQueryUrl
	} else {
		url = RefundQueryUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(xmlStr)
}

// Reverse 撤销订单
func (c *Client) Reverse(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxReverseUrl
	} else {
		url = ReverseUrl
	}
	xmlStr, err := c.postWithCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(xmlStr)
}

// CloseOrder 关闭订单
func (c *Client) CloseOrder(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxCloseOrderUrl
	} else {
		url = CloseOrderUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(xmlStr)
}

// DownloadBill 对账单下载
func (c *Client) DownloadBill(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxDownloadBillUrl
	} else {
		url = DownloadBillUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)

	var p Params

	// 如果出现错误，返回XML数据
	if strings.Index(xmlStr, "<") == 0 {
		p = XMLToMap(xmlStr, "")
		return p, err
	}

	// 正常返回csv数据
	p.SetString("return_code", Success)
	p.SetString("return_msg", "ok")
	p.SetString("data", xmlStr)
	return p, err
}

// Report 交易保障
func (c *Client) Report(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxReportUrl
	} else {
		url = ReportUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(xmlStr)
}

// ShortURL 转换短链接
func (c *Client) ShortURL(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxShortUrl
	} else {
		url = ShortUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(xmlStr)
}

// AuthCodeToOpenid 授权码查询OPENID接口
func (c *Client) AuthCodeToOpenid(params Params) (Params, error) {
	var url string
	if c.isSandbox {
		url = SandboxAuthCodeToOpenidUrl
	} else {
		url = AuthCodeToOpenidUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXML(xmlStr)
}
