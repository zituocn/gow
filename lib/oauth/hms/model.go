package hms

// QuickAppAccessToken 华为accesstoken应答参数
type QuickAppAccessToken struct {
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	Scope          string `json:"scope"`
	TokenType      int64  `json:"token_type"`
	Error          string `json:"error"`
	ErrDescription string `json:"error_description"`
}

