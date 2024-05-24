package wepay

// 常量
const (
	bodyType                    = "application/xml; charset=utf-8"
	Fail                        = "FAIL"
	Success                     = "SUCCESS"
	HMACSHA256                  = "HMAC-SHA256"
	MD5                         = "MD5"
	Sign                        = "sign"
	MicroPayUrl                 = "https://api.mch.weixin.qq.com/pay/micropay"
	UnifiedOrderUrl             = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	OrderQueryUrl               = "https://api.mch.weixin.qq.com/pay/orderquery"
	ReverseUrl                  = "https://api.mch.weixin.qq.com/secapi/pay/reverse"
	CloseOrderUrl               = "https://api.mch.weixin.qq.com/pay/closeorder"
	RefundUrl                   = "https://api.mch.weixin.qq.com/secapi/pay/refund"
	RefundQueryUrl              = "https://api.mch.weixin.qq.com/pay/refundquery"
	DownloadBillUrl             = "https://api.mch.weixin.qq.com/pay/downloadbill"
	ReportUrl                   = "https://api.mch.weixin.qq.com/payitil/report"
	ShortUrl                    = "https://api.mch.weixin.qq.com/tools/shorturl"
	AuthCodeToOpenidUrl         = "https://api.mch.weixin.qq.com/tools/authcodetoopenid"
	ProfitSharingUrl            = "https://api.mch.weixin.qq.com/secapi/pay/profitsharing"       //请求单次分账
	MultiProfitSharingUrl       = "https://api.mch.weixin.qq.com/secapi/pay/multiprofitsharing"  //请求多次分账
	ProfitSharingFinishUrl      = "https://api.mch.weixin.qq.com/secapi/pay/profitsharingfinish" //完结分账
	ProfitSharingQueryUrl       = "https://api.mch.weixin.qq.com/pay/profitsharingquery"         //查询分账结果
	ProfitSharingReturnUrl      = "https://api.mch.weixin.qq.com/secapi/pay/profitsharingreturn" //分账回退
	ProfitSharingReturnQueryUrl = "https://api.mch.weixin.qq.com/pay/profitsharingreturnquery"   //分账回退结果查询
	SandboxMicroPayUrl          = "https://api.mch.weixin.qq.com/sandboxnew/pay/micropay"
	SandboxUnifiedOrderUrl      = "https://api.mch.weixin.qq.com/sandboxnew/pay/unifiedorder"
	SandboxOrderQueryUrl        = "https://api.mch.weixin.qq.com/sandboxnew/pay/orderquery"
	SandboxReverseUrl           = "https://api.mch.weixin.qq.com/sandboxnew/secapi/pay/reverse"
	SandboxCloseOrderUrl        = "https://api.mch.weixin.qq.com/sandboxnew/pay/closeorder"
	SandboxRefundUrl            = "https://api.mch.weixin.qq.com/sandboxnew/secapi/pay/refund"
	SandboxRefundQueryUrl       = "https://api.mch.weixin.qq.com/sandboxnew/pay/refundquery"
	SandboxDownloadBillUrl      = "https://api.mch.weixin.qq.com/sandboxnew/pay/downloadbill"
	SandboxReportUrl            = "https://api.mch.weixin.qq.com/sandboxnew/payitil/report"
	SandboxShortUrl             = "https://api.mch.weixin.qq.com/sandboxnew/tools/shorturl"
	SandboxAuthCodeToOpenidUrl  = "https://api.mch.weixin.qq.com/sandboxnew/tools/authcodetoopenid"
)
