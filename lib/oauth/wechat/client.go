package wechat

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/imroc/req"
	"io"
	"strings"
	"time"
)

const (
	//全局token
	globalTokenUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	//code换token
	codeToAccessTokenUrl = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	//刷新token
	refreshTokenUrl = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	//用户信息
	userInfoUrl = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s"
	//ticket
	ticketUrl = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%v&type=jsapi"
)

// Client webchat client
type Client struct {
	AppId  string
	Secret string
	ApiKey string //仅支付时需要
}

// NewClient NewClient
// 传入appId和appSecret
func NewClient(appId, secret string) *Client {
	return &Client{
		AppId:  appId,
		Secret: secret,
	}
}

// SetApiKey 公众号支付时，必须设置此值
func (c *Client) SetApiKey(apiKey string) {
	c.ApiKey = apiKey
}

// CodeToAccessToken 根据code值换行accessToken
// 第一步
func (c *Client) CodeToAccessToken(code string) (accessData *AccessData, err error) {
	if code == "" {
		err = fmt.Errorf("[wechat]code为空")
		return
	}
	url := fmt.Sprintf(codeToAccessTokenUrl, c.AppId, c.Secret, code)
	req.SetTimeout(10 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		return
	}
	accessData = new(AccessData)
	err = resp.ToJSON(&accessData)
	if err != nil {
		return
	}
	if accessData.ErrCode != 0 {
		err = fmt.Errorf("[wechat]获取code出错 %v", accessData.ErrMsg)
		return
	}

	return
}

// GetWxUser 根据accessToken和openid获取用户的基本信息
// 第二步
func (c *Client) GetWxUser(accessToken, openid string) (wxUser *WxUser, err error) {
	if accessToken == "" || openid == "" {
		err = fmt.Errorf("[wechat]accessToken或openido为空")
		return
	}
	url := fmt.Sprintf(userInfoUrl, accessToken, openid)
	req.SetTimeout(10 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		return
	}

	wxUser = new(WxUser)
	err = resp.ToJSON(&wxUser)
	if err != nil {
		return
	}

	if wxUser.ErrCode != 0 {
		err = fmt.Errorf("[wechat]获取用户信息失败 %v", wxUser.ErrMsg)
		return
	}

	return
}

// GetMpGlobalToken 获取微信公众号的全局token
// expiresIn:过期时间
// 此值返回后，需要存放起来，过期时间通常为：7200
func (c *Client) GetMpGlobalToken() (mpTokenData *MpTokenData, err error) {
	url := fmt.Sprintf(globalTokenUrl, c.AppId, c.Secret)
	req.SetTimeout(10 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		return
	}
	mpTokenData = new(MpTokenData)
	err = resp.ToJSON(&mpTokenData)
	if err != nil {
		return
	}
	if mpTokenData.ErrCode != 0 {
		err = fmt.Errorf("[wechat]获取全局token失败 %v", mpTokenData.ErrMsg)
		return
	}
	if mpTokenData.AccessToken == "" || mpTokenData.ExpiresIn == 0 {
		err = fmt.Errorf("[wechat]未获取到全局token")
		return
	}

	return
}

// GetJSAPITicket 根据全局token，获取ticket
// 返回的ticket需要存放起来，过期时间通常为:7200
func (c *Client) GetJSAPITicket(mpGlobalToken string) (ticketData *TicketData, err error) {
	if mpGlobalToken == "" {
		err = fmt.Errorf("[wechat]全局token为空")
		return
	}
	url := fmt.Sprintf(ticketUrl, &mpGlobalToken)
	req.SetTimeout(10 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		return
	}
	ticketData = new(TicketData)
	err = resp.ToJSON(&ticketData)
	if err != nil {
		return
	}
	if ticketData.ErrCode != 0 {
		err = fmt.Errorf("[wechat]获取全局ticket失败 %v", ticketData.ErrMsg)
		return
	}
	if ticketData.Ticket == "" || ticketData.ExpiresIn == 0 {
		err = fmt.Errorf("[wechat]未获取到全局ticket")
		return
	}
	return
}

// GetJSAPIConfig 获取微信公众号js sdk api 的配置信息
// 输出到H5页面
// ticket为暂存的ticket值
// url为H5页面的完整地址
func (c *Client) GetJSAPIConfig(ticket, url string) (data *JSWXConfig, err error) {
	if ticket == "" {
		err = fmt.Errorf("[wechat]ticket空")
		return
	}

	timeStamp := time.Now().Unix()
	//随机字串
	nonceStr := fmt.Sprintf("%d", timeStamp)

	//获取支付签名：sign
	paySignStr := fmt.Sprintf("jsapi_ticket=%v&noncestr=%v&timestamp=%v&url=%v", ticket, nonceStr, timeStamp, url)
	h := sha1.New()
	io.WriteString(h, paySignStr)
	sign := fmt.Sprintf("%x", h.Sum(nil))

	//组装输出配置信息
	data = &JSWXConfig{
		Debug:     false,
		AppID:     c.AppId,
		TimeStamp: timeStamp,
		NonceStr:  nonceStr,
		Signature: sign,
		JsAPIList: []string{
			"updateAppMessageShareData",
			"updateTimelineShareData",
			"onMenuShareWeibo",
			"onMenuShareTimeline",
			"onMenuShareAppMessage",
			"onMenuShareQQ",
			"onMenuShareQZone",
		},
	}

	return
}

// GetPayConfig 获取公众号支付时的参数
//
//	传入payParam为支付下单接口返回的信息
func (c *Client) GetPayConfig(payParam string) (payConfig *PayConfig, err error) {
	if payParam == "" {
		err = fmt.Errorf("[wechat]payParam为空")
		return
	}
	timeStamp := time.Now().Unix()
	//随机字串
	nonceStr := fmt.Sprintf("%d", timeStamp)
	packageStr := fmt.Sprintf("prepay_id=%v", payParam)
	urlStr := fmt.Sprintf("appId=%v&nonceStr=%v&package=%v&signType=%v&timeStamp=%v&key=%v", c.AppId, nonceStr, packageStr, "MD5", timeStamp, c.ApiKey)
	sign := sign(urlStr)
	payConfig = &PayConfig{
		AppID:     c.AppId,
		TimeStamp: timeStamp,
		NonceStr:  nonceStr,
		Package:   packageStr,
		SignType:  "MD5",
		PaySign:   sign,
	}

	return
}

// sign pay时的签名方法
func sign(signStr string) (sign string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStr))
	cipherStr := md5Ctx.Sum(nil)
	sign = strings.ToUpper(hex.EncodeToString(cipherStr))
	return
}
