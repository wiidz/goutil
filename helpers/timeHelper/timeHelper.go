package timeHelper

import (
	"fmt"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"time"
)

/**
 * @func: FromTodayToTomorrowTimeStamp  返回今天凌晨和明天凌晨的时间戳
 * @author Wiidz
 * @date   2019-11-16
 */
func  FromTodayToTomorrowTimeStamp() (today, tomorrow int64) {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tm2 := tm1.AddDate(0, 0, 1)
	return tm1.Unix(), tm2.Unix()
}

//
/**
 * @func: LastDayOfTimeStamp  获取本日最后一天的时间戳
 * @author Wiidz
 * @date   2019-11-16
 */
func  LastDayOfTimeStamp(d time.Time) int64 {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTimeStamp(d).Unix()
}

/**
 * @func: GetZeroTimeStamp  获取某一天的0点时间
 * @author Wiidz
 * @date   2019-11-16
 */
func  GetZeroTimeStamp(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

/**
 * @func: BeautyTimeStamp 美化时间
 * @author Wiidz
 * @date   2019-11-16
 */

func  BeautyTimeStamp(timeStamp, currentTime int64) string {
	if currentTime == 0 {
		currentTime = time.Now().Unix()
	}

	span := currentTime - timeStamp
	if span < 60 {

		return "刚刚"

	} else if span < 3600 {

		tmp := int64(span / 60)

		return typeHelper.Int64ToStr(tmp) + "分钟前"

	} else if span < 24*3600 {

		tmp := int64(span / 3600)

		return typeHelper.Int64ToStr(tmp) + "小时前"

	} else if span < (7 * 24 * 3600) {

		tmp := int64(span / 24 * 3600)

		return typeHelper.Int64ToStr(tmp) + "天前"

	} else {

		tm := time.Unix(timeStamp, 0)

		return tm.Format("2006-01-02 03:04:05")

	}

}

/**
 * @func: GetISO8601 获取iso8601格式的时间
 * @author Wiidz
 * @date   2019-11-16
 */
func  GetISO8601(date int64) string {
	var formattedDate = time.Unix(date, 0).Format("2006-01-02T15:04:05Z")
	return formattedDate
}


const (
	dateTimeFormat      = "2006-01-02 15:04:05"
	dateTimeFormatSlash = "2006/01/02 15:04:05"
	dateTime            = "2006-01-02"
	dateTimeSlash       = "2006/01/02"
)

type MyJsonTime time.Time


// GetTimePoint 实现它的json序列化方法 注意测试
func (tm MyJsonTime) GetTimePoint() *time.Time {
	//temp,_:= time.Parse("2006-01-02 15:04:05",tm.GetDateTimeStr())
	local, _ := time.LoadLocation("Local")
	temp, _ := time.ParseInLocation("2006-01-02 15:04:05",tm.GetDateTimeStr(),local)

	return  &temp
}

// MarshalJSON 实现它的json序列化方法
func (tm MyJsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(tm).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// GetDateTimeStr 实现它的json序列化方法
func (tm MyJsonTime) GetDateTimeStr() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(dateTimeFormat)
}

// GetDateTimeStrSlash 实现它的json序列化方法
func (tm MyJsonTime) GetDateTimeStrSlash() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(dateTimeFormatSlash)
}

// GetDateStr 获取string格式
func (tm MyJsonTime) GetDateStr() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(dateTime)
}
// GetDateStrSlash 实现它的json序列化方法
func (tm MyJsonTime) GetDateStrSlash() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(dateTimeSlash)
}
// IsNull 判断是否为空
func (tm MyJsonTime) IsNull() bool {
	return time.Time(tm) == time.Time{}
}

func (tm MyJsonTime) AddDate(years int, months int, days int) MyJsonTime {
	temp := time.Time(tm)
	temp = temp.AddDate(years, months, days)
	return MyJsonTime(temp)
}

func (tm MyJsonTime) Day() int {
	temp := time.Time(tm)
	return temp.Day()
}

func (tm MyJsonTime) After(target time.Time) bool {
	temp := time.Time(tm)
	return temp.After(target)
}
func (tm MyJsonTime) Before(target time.Time) bool {
	temp := time.Time(tm)
	return temp.Before(target)
}

func (tm MyJsonTime) Equal(target time.Time) bool {
	temp := time.Time(tm)
	return temp.Equal(target)
}
func (tm MyJsonTime) Year() int {
	return time.Time(tm).Year()
}

func (tm MyJsonTime) Month() time.Month {
	return time.Time(tm).Month()
}

func (tm MyJsonTime) Location() *time.Location {
	return time.Time(tm).Location()
}

// GetFirstDateOfMonth 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func (tm MyJsonTime) GetFirstDateOfMonth() *MyJsonTime {
	newTime := tm.AddDate(0, 0, -tm.Day()+1)
	res := newTime.GetZeroTime()
	return &res
}

// GetZeroTime 获取某一天的0点时间
func (tm MyJsonTime) GetZeroTime() MyJsonTime {
	return MyJsonTime(time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location()))
}


// GetLastTime 获取某一天的最后时间
func (tm MyJsonTime) GetLastTime() MyJsonTime {
	return MyJsonTime(time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location()))
}


func (tm MyJsonTime) Format2Time() time.Time {
	return time.Time(tm)
}