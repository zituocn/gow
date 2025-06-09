package harmony

type HarmonyAccessToken struct {
	TokenType      string `json:"token_type"` //固定字符串 Bearer
	AccessToken    string `json:"access_token"`
	Scope          string `json:"scope"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	IdToken        string `json:"id_token"`
	Error          int    `json:"error"`
	SubError       int    `json:"sub_error"`
	ErrDescription string `json:"error_description"`
}

type HarmonyAccountInfo struct {
	OpenID            string `json:"openID"`
	UnionID           string `json:"unionID"`
	LoginMobileNumber string `json:"loginMobileNumber"`
	LoginMobileValid  string `json:"loginMobileValid"`
	PurePhoneNumber   string `json:"purePhoneNumber"`
	PhoneCountryCode  string `json:"phoneCountryCode"`
	Error             string `json:"error"`
	ErrorCode         string `json:"errorCode"` //在http响应头NSP_STATUS字段获取
}
