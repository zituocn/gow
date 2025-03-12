package sms

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

// TencentSmsClient tencent sms client
type TencentSmsClient struct {
	SDKAppID  string
	SecretId  string
	SecretKey string
}

// NewTencentSmsClient return TencentSmsClient
func NewTencentSmsClient(sdkAppId, secretId, secretKey string) *TencentSmsClient {
	return &TencentSmsClient{
		SDKAppID:  sdkAppId,
		SecretId:  secretId,
		SecretKey: secretKey,
	}
}

// SendVerifyCode 验证码短信
func (m *TencentSmsClient) SendVerifyCode(sign, templateID, phone, code string) (err error) {
	return m.send(sign, templateID, []string{fmt.Sprintf("+86%s", phone)}, []string{code})
}

// SendMarket 营销类短信 批量号码 发送相同的内容
func (m *TencentSmsClient) SendMarket(sign, templateID string, phone, templateParam []string) (err error) {
	var phoneData []string
	for _, item := range phone {
		phoneData = append(phoneData, fmt.Sprintf("+86%s", item))
	}
	return m.send(sign, templateID, phoneData, templateParam)
}

// 参考文档：https://cloud.tencent.com/document/product/382/43199
func (m *TencentSmsClient) send(sign, templateID string, phone, templateParam []string) (err error) {
	credential := common.NewCredential(m.SecretId, m.SecretKey)
	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()
	/* SDK默认使用POST方法。
	 * 如果您一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 60 // 请求超时时间，单位为秒(默认60秒)
	/* 指定接入地域域名，默认就近地域接入域名为 sms.tencentcloudapi.com ，也支持指定地域域名访问，例如广州地域的域名为 sms.ap-guangzhou.tencentcloudapi.com */
	//cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	cpf.HttpProfile.Endpoint = "sms.ap-guangzhou.tencentcloudapi.com"
	/* SDK默认用TC3-HMAC-SHA256进行签名，非必要请不要修改这个字段 */
	//cpf.SignMethod = "HmacSHA1"

	/* 实例化要请求产品(以sms为例)的client对象
	 * 第二个参数是地域信息，可以直接填写字符串ap-guangzhou，支持的地域列表参考 https://cloud.tencent.com/document/api/382/52071#.E5.9C.B0.E5.9F.9F.E5.88.97.E8.A1.A8 */
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	/* 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	 * 您可以直接查询SDK源码确定接口有哪些属性可以设置
	 * 属性可能是基本类型，也可能引用了另一个数据结构
	 * 推荐使用IDE进行开发，可以方便的跳转查阅各个接口和数据结构的文档说明 */
	request := sms.NewSendSmsRequest()
	/* 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666 */
	request.SmsSdkAppId = common.StringPtr(m.SDKAppID)
	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名 */
	request.SignName = common.StringPtr(sign)
	/* 模板 ID: 必须填写已审核通过的模板 ID */
	request.TemplateId = common.StringPtr(templateID)
	/* 模板参数: 模板参数的个数需要与 TemplateId 对应模板的变量个数保持一致，若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs(templateParam)
	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs(phone)
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.SendSms(request)

	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		err = fmt.Errorf("调用[腾讯云API]发送短信出错：%v", err)
		return
	}

	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	//if err != nil {
	//	panic(err)
	//}

	//b, _ := json.Marshal(response.Response)
	// 打印返回的json字符串
	//logx.Errorf("%s", b)

	// {"SendStatusSet":[{"SerialNo":"9331:261873439117417597983741023","PhoneNumber":"+8615095910236","Fee":1,"SessionContext":"","Code":"Ok","Message":"send success","IsoCode":"CN"}],"RequestId":"46f1c9ae-a917-4fb0-bd52-2b609415890e"}

	if response.Response != nil && len(response.Response.SendStatusSet) > 0 {
		//暂时1次只发送1个号码
		//logx.Errorf("code:%v ,,,, message:%v", *response.Response.SendStatusSet[0].Code, *response.Response.SendStatusSet[0].Message)
		if *response.Response.SendStatusSet[0].Code != "Ok" {
			err = fmt.Errorf("[腾讯云]发送短信失败：%v", *response.Response.SendStatusSet[0].Message)
			return
		}
	}
	return
}
