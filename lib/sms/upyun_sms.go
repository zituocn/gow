package sms

import (
	"fmt"
	"github.com/imroc/req"
	"time"
)

const (
	//api Url
	smsUrl = "https://sms-api.upyun.com/api/messages"
)

// UpYunSmsClient UpYunSmsClient
type UpYunSmsClient struct {
	Token string
}

// NewUpYunSmsClient 返回一个新的client
func NewUpYunSmsClient(token string) *UpYunSmsClient {
	return &UpYunSmsClient{
		Token: token,
	}
}

// SendVerifyCode 发送验证码短信
//
//	templateId = 模板ID
func (m *UpYunSmsClient) SendVerifyCode(mobile string, templateId int64, code string) (messageId string, err error) {
	if mobile == "" {
		err = fmt.Errorf("手机号格式不正确")
		return
	}
	body := new(UpYunSmsBody)
	body.Mobile = mobile
	body.TemplateID = templateId
	body.Vars = fmt.Sprintf("%v", code)
	messageId, err = m.send(body)
	if err != nil {
		return
	}
	return
}

// send 发送请求
func (m *UpYunSmsClient) send(body *UpYunSmsBody) (messageId string, err error) {
	req.SetTimeout(5 * time.Second)
	header := req.Header{
		"Authorization": m.Token,
		"Content-Type":  "application/json",
		"User-Agent":    "golang 1.10",
	}
	resp, err := req.Post(smsUrl, header, req.BodyJSON(&body))
	if err != nil {
		err = fmt.Errorf("向upyun发送请求错误:%v", err)
		return
	}
	if resp != nil && resp.Response().StatusCode != 200 {
		err = fmt.Errorf("向upyun发送请求，返回错误码:%v", resp.Response().StatusCode)
		return
	}
	ret := new(UpYunResult)
	err = resp.ToJSON(&ret)
	if err != nil {
		err = fmt.Errorf("发送返回信息失败:%v", err)
		return
	}

	if len(ret.MessageIDS) == 0 || ret.MessageIDS[0] == nil {
		err = fmt.Errorf("向upyun发送后，错误的返回:%v", ret)
		return
	}

	if len(ret.MessageIDS) == 1 {
		messageId = fmt.Sprintf("%d", ret.MessageIDS[0].MessageID)
	}

	return
}

// UpYunSmsBody 发送模型
type UpYunSmsBody struct {
	Mobile     string `json:"mobile"`
	TemplateID int64  `json:"template_id"`
	Vars       string `json:"vars"`
}

// UpYunResult  返回信息
type UpYunResult struct {
	MessageIDS []*UpYunMessage `json:"message_ids"`
}

// UpYunMessage UpYunMessage
type UpYunMessage struct {
	MessageID int64  `json:"message_id"`
	Mobile    string `json:"mobile"`
}
