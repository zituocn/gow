package alipay

import (
	"net/url"
)

// TradePagePay TradePagePay
// https://docs.open.alipay.com/270/alipay.trade.page.pay
func (t *AliPay) TradePagePay(param *TradePagePay) (results *url.URL, err error) {
	p, err := t.URLValues(param)
	if err != nil {
		return nil, err
	}
	results, err = url.Parse(t.apiDomain + "?" + p.Encode())
	if err != nil {
		return nil, err
	}
	return results, err
}

// TradeAppPay TradeAppPay
// https://docs.open.alipay.com/api_1/alipay.trade.app.pay
func (t *AliPay) TradeAppPay(param *TradeAppPay) (results string, err error) {
	p, err := t.URLValues(param)
	if err != nil {
		return "", err
	}
	return p.Encode(), err
}

// TradeFastpayRefundQuery TradeFastpayRefundQuery
// https://docs.open.alipay.com/api_1/alipay.trade.fastpay.refund.query
func (t *AliPay) TradeFastpayRefundQuery(param *FastpayTradeRefundQuery) (results *FastpayTradeRefundQueryResponse, err error) {
	err = t.doRequest("POST", param, &results)
	return results, err
}

// TradeOrderSettle https://docs.open.alipay.com/api_1/alipay.trade.order.settle
func (t *AliPay) TradeOrderSettle(param *TradeOrderSettle) (results *TradeOrderSettleResponse, err error) {
	err = t.doRequest("POST", param, &results)
	return results, err
}

// TradeClose https://docs.open.alipay.com/api_1/alipay.trade.close/
func (t *AliPay) TradeClose(param TradeClose) (results *TradeCloseResponse, err error) {
	err = t.doRequest("POST", param, &results)
	return results, err
}

// TradeCancel https://docs.open.alipay.com/api_1/alipay.trade.cancel/
func (t *AliPay) TradeCancel(param TradeCancel) (results *TradeCancelResponse, err error) {
	err = t.doRequest("POST", param, &results)
	return results, err
}

// TradeRefund https://docs.open.alipay.com/api_1/alipay.trade.refund/
func (t *AliPay) TradeRefund(param TradeRefund) (results *TradeRefundResponse, err error) {
	err = t.doRequest("POST", param, &results)
	return results, err
}

// TradePreCreate https://docs.open.alipay.com/api_1/alipay.trade.precreate/
func (t *AliPay) TradePreCreate(param TradePreCreate) (results *TradePreCreateResponse, err error) {
	err = t.doRequest("POST", param, &results)
	return results, err
}

// TradeQuery 交易订单查询 https://docs.open.alipay.com/api_1/alipay.trade.query/
//		results.Msg=="Success" && results.AliPayTradeQuery.TradeStatus == "TRADE_SUCCESS" => 支付成功
func (t *AliPay) TradeQuery(param *TradeQuery) (results *TradeQueryResponse, err error) {
	err = t.doRequest("POST", param, &results)
	return results, err
}

// TradeCreate https://docs.open.alipay.com/api_1/alipay.trade.create/
func (t *AliPay) TradeCreate(param TradeCreate) (results *TradeCreateResponse, err error) {
	err = t.doRequest("POST", param, &results)
	return results, err
}

// TradePay https://docs.open.alipay.com/api_1/alipay.trade.pay/
func (t *AliPay) TradePay(param AliPayTradePay) (results *AliPayTradePayResponse, err error) {
	err = t.doRequest("POST", param, &results)
	return results, err
}

// TradeOrderInfoSync https://docs.open.alipay.com/api_1/alipay.trade.orderinfo.sync/
func (t *AliPay) TradeOrderInfoSync(param TradeOrderInfoSync) (results *TradeOrderInfoSyncResponse, err error) {
	err = t.doRequest("POST", param, &results)
	return results, err
}
