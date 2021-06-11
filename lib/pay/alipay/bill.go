package alipay

// BillDownloadURLQuery https://docs.open.alipay.com/api_15/alipay.data.dataservice.bill.downloadurl.query
func (m *AliPay) BillDownloadURLQuery(param BillDownloadURLQuery) (results *BillDownloadURLQueryResponse, err error) {
	err = m.doRequest("POST", param, &results)
	return results, err
}
