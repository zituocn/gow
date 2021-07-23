package toutiao

// AccessToken model
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Error       int64  `json:"error"`
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	Message     string `json:"message"`
}
