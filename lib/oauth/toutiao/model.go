package toutiao

//AccessToken model
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Error       int64  `json:"error"`
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	Message     string `json:"message"`
}

//SessionData model
type SessionData struct {
	SessionKey      string `json:"session_key"`
	Openid          string `json:"openid"`
	AnonymousOpenid string `json:"anonymous_openid"`
	Unionid         string `json:"unionid"`
	Error           int64  `json:"error"`
	ErrCode         int64  `json:"errcode"`
	ErrMsg          string `json:"errmsg"`
	Message         string `json:"message"`
}

// TTAppletUserInfo 敏感用户信息
type TTAppletUserInfo struct {
	OpenID    string `json:"openId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`//0: 未知；1:男性；2:女性
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	Watermark struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
}