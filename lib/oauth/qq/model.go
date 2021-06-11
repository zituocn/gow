package qq


// AccessData access token response
type AccessData struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// OpenIdData openId response
type OpenIdData struct {
	ClientId string `json:"client_id"`
	OpenId   string `json:"openid"`
	UnionId  string `json:"unionid"`
}

// QUser QQ返回的用户
type QUser struct {
	Ret        int64  `json:"ret"`
	Msg        string `json:"msg"`
	NickName   string `json:"nickname"`
	FigureUrl1 string `json:"figureurl_qq_1"` //QQ头像1 40*40
	FigureUrl2 string `json:"figureurl_qq_2"` //QQ头像2 100*100
}
