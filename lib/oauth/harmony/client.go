package harmony

import (
	"encoding/json"
	"fmt"
	"github.com/zituocn/logx"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	reqBodyType       = "application/x-www-form-urlencoded"
	respBodyType      = "application/json;charset=UTF-8"
	getAccessTokenUrl = "https://oauth-login.cloud.huawei.com/oauth2/v3/token"
	getAccountInfoUrl = "https://account.cloud.huawei.com/rest.php?nsp_svc=GOpen.User.getInfo"
)

// Client harmony client
type Client struct {
	ClientId     string
	ClientSecret string
}

// NewClient NewClient
//
//	传入appId和appSecret
func NewClient(clientId, clientSecret string) *Client {
	return &Client{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
}

// CodeToAccessToken 根据code值换行accessToken
func (c *Client) CodeToAccessToken(code string) (accessData *HarmonyAccessToken, err error) {
	if code == "" {
		err = fmt.Errorf("[hms]code为空")
		return
	}
	accessData = new(HarmonyAccessToken)

	param := url.Values{}
	param.Add("grant_type", "authorization_code")
	param.Add("client_id", c.ClientId)
	param.Add("client_secret", c.ClientSecret)
	param.Add("code", code)

	reqData := param.Encode()

	payload := strings.NewReader(reqData)

	client := &http.Client{}

	req, err := http.NewRequest("POST", getAccessTokenUrl, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	if err != nil {
		return
	}
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	json.Unmarshal(resp, &accessData)

	logx.Errorf("accessToken : %v", accessData.AccessToken)

	logx.Errorf("accessData.Error : %v , subError: %v, errMsg: %v", accessData.Error, accessData.SubError, accessData.ErrDescription)

	//if accessData.Error != 200 {
	//	err = fmt.Errorf("[harmony]通过code获取accessToken出错 %v", accessData.ErrDescription)
	//	return
	//}
	return
}

// RefreshTokenToAccessToken 用refreshToken换accessToken
func (c *Client) RefreshTokenToAccessToken(refreshToken string) (accessData *HarmonyAccessToken, err error) {
	if refreshToken == "" {
		err = fmt.Errorf("[hms]refreshToken为空")
		return
	}
	accessData = new(HarmonyAccessToken)
	client := &http.Client{}
	param := url.Values{}
	param.Add("grant_type", "refresh_token")
	param.Add("client_id", c.ClientId)
	param.Add("client_secret", c.ClientSecret)
	param.Add("refresh_token", refreshToken)
	buf := strings.NewReader(param.Encode())
	req, err := http.NewRequest("POST", getAccessTokenUrl, buf)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", reqBodyType)

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

	//if accessData.Error != 200 {
	//	err = fmt.Errorf("[harmony]通过refreshToken获取accessToken出错 %v", accessData.ErrDescription)
	//	return
	//}
	return
}

func (c *Client) GetAccountInfo(accessToken string) (accountInfo *HarmonyAccountInfo, err error) {
	client := &http.Client{}
	param := url.Values{}
	param.Add("access_token", accessToken)
	buf := strings.NewReader(param.Encode())
	req, err := http.NewRequest("POST", getAccountInfoUrl, buf)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", reqBodyType)

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
	json.Unmarshal(resp, &accountInfo)

	//在处理成功时不回返回
	status := response.Header.Get("NSP_STATUS")
	logx.Errorf("NSP_STATUS: %v", status)

	accountInfo.ErrorCode = status

	logx.Errorf("error: %v", accountInfo.Error)

	if accountInfo.Error != "" {
		err = fmt.Errorf("[harmony]获取账号信息出错 %v", accountInfo.Error)
		return
	}
	return
}
