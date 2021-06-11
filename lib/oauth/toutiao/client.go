package toutiao

import (
	"fmt"
	"github.com/imroc/req"
	"time"
)

//Client Client
type Client struct {
	AppId  string
	Secret string
}

const (
	bodyType          = "application/x-www-form-urlencoded"
	getAccessTokenUrl = "https://developer.toutiao.com/api/apps/token?appid=%s&secret=%s&grant_type=client_credential"
)

//NewClient NewClient
//传入appId和appSecret
func NewClient(appId, secret string) *Client {
	return &Client{
		AppId:  appId,
		Secret: secret,
	}
}

//GetAccessToken 获得头条的接口登录凭证
func (c *Client) GetAccessToken() (accessData *AccessToken, err error) {
	url := fmt.Sprintf(getAccessTokenUrl, c.AppId, c.Secret)
	req.SetTimeout(10 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		return
	}
	accessData = new(AccessToken)
	err = resp.ToJSON(&accessData)
	if err != nil {
		return
	}
	if accessData.ErrCode != 0 {
		err = fmt.Errorf("[toutiao]获取accessToken出错 %v", accessData.ErrMsg)
		return
	}
	return
}
