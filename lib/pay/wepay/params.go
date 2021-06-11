package wepay

import "strconv"

//Params Params
type Params map[string]interface{}

// SetString SetString
func (p Params) SetString(k, s string) Params {
	p[k] = s
	return p
}

//GetString GetString
func (p Params) GetString(k string) string {
	s, _ := p[k]
	str, _ := s.(string)
	return str
}

//SetInt64 SetInt64
func (p Params) SetInt64(k string, i int64) Params {
	p[k] = strconv.FormatInt(i, 10)
	return p
}

//GetInt64 GetInt64
func (p Params) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}

// ContainsKey 判断key是否存在
func (p Params) ContainsKey(key string) bool {
	_, ok := p[key]
	return ok
}
