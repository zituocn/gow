package alipay

import (
	"net/http"
	"net/url"
	"strings"
)

// TradeWapPay https://docs.open.alipay.com/api_1/alipay.trade.wap.pay/
// TradeWapPay
func (m *AliPay) TradeWapPay(param *AliPayTradeWapPay) (url *url.URL, err error) {
	p, err := m.URLValues(param)
	if err != nil {
		return nil, err
	}
	var buf = strings.NewReader(p.Encode())

	req, err := http.NewRequest("POST", m.apiDomain, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", K_CONTENT_TYPE_FORM)

	rep, err := m.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err != nil {
		return nil, err
	}
	url = rep.Request.URL
	return url, err
}
