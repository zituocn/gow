/*

基于 context 的扩展方法

处理了：
	1. 自定义JSON输出格式；
	2. 常用的翻页处理；

sam
2021-06-07

*/

package gow

import "time"


// DataResponse data struct
type DataResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Time int    `json:"time"`
	Body *Body  `json:"body"`
}

// Body body struct
type Body struct {
	Pager *Pager      `json:"pager"`
	Data  interface{} `json:"data"`
}

// Pager pager struct
type Pager struct {
	Page      int64 `json:"page"`
	Limit     int64 `json:"-"`
	Offset    int64 `json:"-"`
	Count     int64 `json:"count"`
	PageCount int64 `json:"page_count"`
}

// DataPager middleware
//	r.Use(gow.DataPager())
// 	like : /v1/user/page?page=1&limit=20
func DataPager() HandlerFunc {
	return func(c *Context) {
		pager := &Pager{}

		pager.Page, _ = c.GetInt64("page", 1)
		if pager.Page < 1 {
			pager.Page = 1
		}

		pager.Limit, _ = c.GetInt64("limit", 10)
		if pager.Limit < 1 {
			pager.Limit = 1
		}

		pager.Offset = (pager.Page - 1) * pager.Limit
		c.Pager = pager

		c.Next()
	}

}

// ServerDataJSON response JSON format
//	like:  c.ServerDataJSON(401,1,"UnAuthorized")
func (c *Context) ServerDataJSON(statusCode int, args ...interface{}) {
	var (
		err   error
		pager *Pager
		data  interface{}
		msg   string
		code  int
	)
	for _, v := range args {
		switch vv := v.(type) {
		case int:
			code = vv
		case string:
			msg = vv
		case error:
			err = vv
		case *Pager:
			pager = vv
		default:
			data = vv
		}
	}
	if err != nil {
		debugPrint("[error] %s %s", c.Request.URL.String(), err.Error())
	}
	if code == 0 && msg == "" {
		msg = "success"
	}

	body := new(Body)

	if pager != nil {
		pager.PageCount = getPageCount(pager.Count, pager.Limit)
	} else {
		pager = &Pager{}
	}
	body.Pager = pager
	body.Data = data

	resp := &DataResponse{
		Code: code,
		Msg:  msg,
		Time: int(time.Now().Unix()),
		Body: body,
	}
	c.ServerJSON(statusCode, &resp)
	return
}

// DataJSON response JSON format
//	c.DataJSON(1,"lost param")
//	c.DataJSON()
func (c *Context) DataJSON(args ...interface{}) {
	c.ServerDataJSON(200, args...)
}

// getPageCount return pagerCount
func getPageCount(count, limit int64) (pageCount int64) {
	if count > 0 && limit > 0 {
		if count%limit == 0 {
			pageCount = count / limit
		} else {
			pageCount = (count / limit) + 1
		}
	}
	return pageCount
}
