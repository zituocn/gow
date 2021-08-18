package toutiao

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req"
	"github.com/zituocn/gow/lib/util"
	"time"
)

//Client 头条oauth2 client
type Client struct {
	AppId  string
	Secret string
}

const (
	bodyType          = "application/x-www-form-urlencoded"
	getAccessTokenUrl = "https://developer.toutiao.com/api/apps/token?appid=%s&secret=%s&grant_type=client_credential"
	code2SessionUrl   = "https://developer.toutiao.com/api/apps/jscode2session?appid=%s&secret=%s&code=%s"
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

// CodeToSession 通过code换取session_key 和 openId
func (c *Client) CodeToSession(code string) (sessionData *SessionData, err error) {
	url := fmt.Sprintf(code2SessionUrl, c.AppId, c.Secret, code)
	req.SetTimeout(10 * time.Second)
	resp, err := req.Get(url)
	if err != nil {
		return
	}
	sessionData = new(SessionData)
	err = resp.ToJSON(&sessionData)
	if err != nil {
		return
	}
	if sessionData.ErrCode != 0 {
		err = fmt.Errorf("[toutiao]code换取sessionData出错 %v", sessionData.ErrMsg)
		return
	}
	return
}

// Decrypt 解密敏感信息
func (c *Client) Decrypt(sessionKey, encryptedData string, iv string) (*TTAppletUserInfo, error) {
	decrypted := new(TTAppletUserInfo)
	aesPlantText, err := util.AppletDecrypt(sessionKey, encryptedData, iv)
	if err != nil {
		return decrypted, err
	}
	err = json.Unmarshal(aesPlantText, &decrypted)
	if err != nil {
		return decrypted, err
	}
	if decrypted.Watermark.AppID != c.AppId {
		return decrypted, errors.New("appId is not match")
	}
	return decrypted, nil
}
