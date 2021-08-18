package util

import (
	"fmt"
	"strings"
	"time"
)

// IntDateTime returns int64 time to format
func IntDateTime(val int64) (ret string) {
	if val < 1 {
		return
	}
	ret = IntDateTimeFormat(val, "YYYY-MM-DD HH:mm")
	return
}

// IntDateTimeFormat DateTimeFormat
func IntDateTimeFormat(val int64, format string) (ret string) {
	if val < 1 {
		return
	}
	tm := time.Unix(val, 0)
	ret = TimeFormat(tm, format)
	return
}

// TimeFormat time.Time format
func TimeFormat(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}
