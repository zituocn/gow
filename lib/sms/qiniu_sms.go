package sms

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/sms"
	"github.com/zituocn/logx"
)

type QiNiuSmsClient struct {
	AccessKey string
	SecretKey string
}

func NewQiNiuSmsClient(accessKey, secretKey string) *QiNiuSmsClient {
	return &QiNiuSmsClient{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

func (s *QiNiuSmsClient) SendVerifyCode(signId, templateId, mobile, code string) (err error) {
	_, err = s.send("", templateId, []string{mobile}, map[string]interface{}{"code": code})
	return
}

func (s *QiNiuSmsClient) send(signID, templateId string, mobiles []string, params map[string]interface{}) (jobId string, err error) {
	//鉴权类
	mac := auth.New(s.AccessKey, s.SecretKey)

	manager := sms.NewManager(mac)

	request := sms.MessagesRequest{
		SignatureID: signID,
		TemplateID:  templateId,
		Mobiles:     mobiles,
		Parameters:  params,
	}

	response, err := manager.SendMessage(request)
	if err != nil {
		logx.Errorf("[七牛云]请求发送短信出错:%v", err)
		err = fmt.Errorf("发送返回信息失败:%v", err)
		return
	}
	jobId = response.JobID
	//logx.Errorf("jobId: %v", jobId)
	return
}
