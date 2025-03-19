package sms

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/sms"
	"github.com/qiniu/go-sdk/v7/sms/client"
	"github.com/qiniu/go-sdk/v7/sms/rpc"
	"github.com/zituocn/logx"
	"net/http"
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

func (s *QiNiuSmsClient) SendVerifyCode(signId, templateId, mobile, code string) (messageId string, err error) {
	messageId, err = s.sendSingle("", templateId, mobile, map[string]interface{}{"code": code})
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

func (s *QiNiuSmsClient) sendSingle(signID, templateId string, mobile string, params map[string]interface{}) (messageId string, err error) {
	//鉴权类
	mac := auth.New(s.AccessKey, s.SecretKey)

	manager := NewManager(mac)

	request := SingleMessageRequest{
		SignatureID: signID,
		TemplateID:  templateId,
		Mobile:      mobile,
		Parameters:  params,
	}

	response, err := manager.SendSingleMessage(request)
	if err != nil {
		logx.Errorf("[七牛云]请求发送短信出错:%v", err)
		err = fmt.Errorf("发送返回信息失败:%v", err)
		return
	}
	messageId = response.MessageId
	//logx.Errorf("messageId: %v", messageId)
	return
}

var (
	// Host 为 Qiniu SMS Server API 服务域名
	Host = "https://sms.qiniuapi.com"
)

// Manager 提供了 Qiniu SMS Server API 相关功能
type Manager struct {
	mac    *auth.Credentials
	client rpc.Client
}

// NewManager 用来构建一个新的 Manager
func NewManager(mac *auth.Credentials) (manager *Manager) {

	manager = &Manager{}

	if mac == nil {
		mac = auth.Default()
	}
	mac1 := &client.Mac{
		AccessKey: mac.AccessKey,
		SecretKey: mac.SecretKey,
	}

	transport := client.NewTransport(mac1, nil)
	manager.client = rpc.Client{Client: &http.Client{Transport: transport}}

	return
}

// MessagesRequest 短信消息
type MessagesRequest struct {
	SignatureID string                 `json:"signature_id"`
	TemplateID  string                 `json:"template_id"`
	Mobiles     []string               `json:"mobiles"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// MessagesResponse 发送短信响应
type MessagesResponse struct {
	JobID     string `json:"job_id,omitempty"`
	MessageId string `json:"message_id,omitempty"`
}

// SendMessage 发送短信 可单条 可多条
func (m *Manager) SendMessage(args MessagesRequest) (ret MessagesResponse, err error) {
	url := fmt.Sprintf("%s%s", Host, "/v1/message")
	err = m.client.CallWithJSON(&ret, url, args)
	return
}

type SingleMessageRequest struct {
	SignatureID string                 `json:"signature_id"`
	TemplateID  string                 `json:"template_id"`
	Mobile      string                 `json:"mobile"`
	Parameters  map[string]interface{} `json:"parameters"`
	Seq         string                 `json:"seq,omitempty"` //业务方提供的消息序列号，状态回调时携带回来
}

// SendSingleMessage 单条发送短信
func (m *Manager) SendSingleMessage(args SingleMessageRequest) (ret MessagesResponse, err error) {
	url := fmt.Sprintf("%s%s", Host, "/v1/message/single")
	err = m.client.CallWithJSON(&ret, url, args)
	return
}
