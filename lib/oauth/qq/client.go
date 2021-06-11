/*
qq oauth 2 登录方法
sam
2020-10-22
*/

package qq

import (
	"fmt"
	"github.com/imroc/req"
	"time"
)

const (

	// 获取token
	// https://wiki.connect.qq.com/%E4%BD%BF%E7%94%A8authorization_code%E8%8E%B7%E5%8F%96access_token
	accessTokenUrl = "https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&fmt=json"

	//获取qq用户信息
	// https://wiki.connect.qq.com/get_user_info
	userInfoUrl = "https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%v&openid=%s"

	//根据accessToken获取用户openId和unionid
	// https://wiki.connect.qq.com/%E8%8E%B7%E5%8F%96%E7%94%A8%E6%88%B7openid_oauth2-0
	openIdUrl = "https://graph.qq.com/oauth2.0/me?access_token=%v&unionid=1&fmt=json"
)

// Client QQ oauth client
type Client struct {
	AppId  string
	Secret string
}

// NewClient return a new qq oauth client
func NewClient(appId, secret string) *Client {
	return &Client{
		AppId:  appId,
		Secret: secret,
	}
}

// GetAccessTokenByCode 根据code获取accessToken
// 主要用到PC网站上
//	第二步
//	之前还有获取code的方法 https://wiki.connect.qq.com/%E4%BD%BF%E7%94%A8authorization_code%E8%8E%B7%E5%8F%96access_token
func (m *Client) GetAccessTokenByCode(code, redirectUrl string) (accessToken string, err error) {
	url := fmt.Sprintf(accessTokenUrl, m.AppId, m.Secret, code, redirectUrl)
	req.SetTimeout(5 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		err = fmt.Errorf("[QQ] 通讯错误：%v", err)
		return
	}

	accessData := new(AccessData)
	resp.ToJSON(&accessData)
	if accessData != nil && accessData.AccessToken != "" {
		accessToken = accessData.AccessToken
	} else {
		err = fmt.Errorf("[QQ] 返回错误:%v", resp.String())
		return
	}
	return
}

// GetOpenIdAndUnionIdByAccessToken 返回openid和unionid
//	主要用到pc网上
//	第三步
func (m *Client) GetOpenIdAndUnionIdByAccessToken(accessToken string) (openId string, unionId string, err error) {
	url := fmt.Sprintf(openIdUrl, accessToken)
	req.SetTimeout(5 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		err = fmt.Errorf("[QQ] 通讯错误：%v", err)
		return
	}

	openIdData := new(OpenIdData)
	resp.ToJSON(&openIdData)

	if openIdData != nil && openIdData.OpenId != "" {
		openId = openIdData.OpenId
		unionId = openIdData.UnionId
	} else {
		err = fmt.Errorf("[QQ] 返回错误:%v", resp.String())
		return
	}
	return
}

// GetQQUser return qq user
//	by openId
//	app端可直接使用此方法
//	pc网站调用时，这是第四步
func (m *Client) GetQQUser(openId, accessToken string) (user *QUser, err error) {
	url := fmt.Sprintf(userInfoUrl, accessToken, m.AppId, openId)
	req.SetTimeout(5 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		err = fmt.Errorf("[QQ] 通讯错误：%v", err)
		return
	}
	user = new(QUser)
	err = resp.ToJSON(&user)
	if err != nil {
		err = fmt.Errorf("[QQ] 反序列化错误：%v", err)
	}
	if user.Ret != 0 {
		err = fmt.Errorf("[QQ] 返回状态错误:%v", user)
		return
	}

	return
}
