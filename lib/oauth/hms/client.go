package hms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	bodyType          = "application/x-www-form-urlencoded"
	getAccessTokenUrl = "https://login.cloud.huawei.com/oauth2/v2/token"
)

// Client hms client
type Client struct {
	AppId     string
	AppSecret string
}

// NewClient NewClient
//	传入appId和appSecret
func NewClient(appId, secret string) *Client {
	return &Client{
		AppId:     appId,
		AppSecret: secret,
	}
}

// CodeToAccessToken 根据code值换行accessToken
func (c *Client) CodeToAccessToken(code string) (accessData *QuickAppAccessToken, err error) {
	if code == "" {
		err = fmt.Errorf("[hms]code为空")
		return
	}
	accessData = new(QuickAppAccessToken)
	client := &http.Client{}
	param := url.Values{}
	param.Add("grant_type", "authorization_code")
	param.Add("code", code)
	param.Add("client_id", c.AppId)
	param.Add("client_secret", c.AppSecret)
	param.Add("redirect_uri", "hms://redirect_url")
	buf := strings.NewReader(param.Encode())
	req, err := http.NewRequest("POST", getAccessTokenUrl, buf)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)

	response, err := client.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return
	}
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	json.Unmarshal(resp, &accessData)
	if accessData.Error != "" {
		err = fmt.Errorf("[hms]获取accessToken出错 %v", accessData.ErrDescription)
		return
	}
	return
}

// RefreshTokenToAccessToken 用refreshToken换accessToken
func (c *Client) RefreshTokenToAccessToken(refreshToken string) (accessData *QuickAppAccessToken, err error) {
	if refreshToken == "" {
		err = fmt.Errorf("[hms]refreshToken为空")
		return
	}
	accessData = new(QuickAppAccessToken)
	client := &http.Client{}
	param := url.Values{}
	param.Add("grant_type", "refresh_token")
	param.Add("refresh_token", refreshToken)
	param.Add("client_id", c.AppId)
	param.Add("client_secret", c.AppSecret)
	buf := strings.NewReader(param.Encode())
	req, err := http.NewRequest("POST", getAccessTokenUrl, buf)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)

	response, err := client.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return
	}
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	json.Unmarshal(resp, &accessData)
	if accessData.Error != "" {
		err = fmt.Errorf("[hms]获取accessToken出错 %v", accessData.ErrDescription)
		return
	}
	return
}
