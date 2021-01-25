package util

import (
	"time"
)

func Strtime2Int(datetime string)(int) {
	//日期转化为时间戳
	timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	loc, _ := time.LoadLocation("Local")    //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix()    //转化为时间戳 类型是int64
	return (int)(timestamp)
}