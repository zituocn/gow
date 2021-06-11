package sms

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/zituocn/gow/lib/logy"
	"github.com/zituocn/gow/lib/util"
)

// AliSmsClient aliyun sms client
type AliSmsClient struct {
	AccessKeyID     string
	AccessKeySecret string
	HTTPDebugEnable bool
}

// NewAliSmsClient return AliSmsClient
func NewAliSmsClient(accessKeyID, accessKeySecret string) *AliSmsClient {
	return &AliSmsClient{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		HTTPDebugEnable: false,
	}
}

// SetDebug SetDebug
func (m *AliSmsClient) SetDebug(enabled bool) {
	m.HTTPDebugEnable = enabled
}

// SendVerifyCode 验证码短信
func (m *AliSmsClient) SendVerifyCode(sign, templateID, phone, code string) (err error) {
	templateParam := fmt.Sprintf(`{"code":"%v"}`, code)
	return m.send(phone, sign, templateID, templateParam)
}

//SendMarket 营销短信
func (m *AliSmsClient) SendMarket(sign, templateId, phone string) (err error) {
	return m.send(phone, sign, templateId, "")
}

// SendNotice 发送通知短信
func (m *AliSmsClient) SendNotice(sign, templateId, phone string, templateParam string) (err error) {
	return m.send(phone, sign, templateId, templateParam)
}

//{"Message":"OK","RequestId":"8295BA5B-0536-463D-BC2D-49A96CAAC37E","BizId":"354401592887345639^0","Code":"OK"}

// send
func (m *AliSmsClient) send(phoneNumbers, sign, templateId, templateParam string) (err error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", m.AccessKeyID, m.AccessKeySecret)
	if err != nil {
		return
	}
	uuid, _ := util.GetUUID()
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phoneNumbers
	request.SignName = sign
	request.TemplateCode = templateId
	request.TemplateParam = templateParam
	request.OutId = uuid
	response, err := client.SendSms(request)
	if err != nil {
		return
	}
	if m.HTTPDebugEnable {
		logy.Debug("response:", response)
	}
	if response.Message != "OK" && response.Code != "OK" {
		err = fmt.Errorf("发送失败：%v", response.Message)
		return
	}
	return

}
