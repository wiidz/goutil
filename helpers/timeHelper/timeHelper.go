package timeHelper

import (
	"fmt"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"time"
)

const (
	HyphenTimeStr = "2006-01-02 15:04:05"
	SlashTimeStr  = "2006/01/02 15:04:05"
	HyphenDateStr = "2006-01-02"
	SlashDateStr  = "2006/01/02"
	ISO8601TimeStr = "2006-01-02T15:04:05Z"
)

var weekdayCn = [][]string{
	{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"},
	{"周日", "周一", "周二", "周三", "周四", "周五", "周六"},
}

// FromTodayToTomorrowTimeStamp : 返回今天凌晨和明天凌晨的时间戳
func FromTodayToTomorrowTimeStamp() (today, tomorrow int64) {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tm2 := tm1.AddDate(0, 0, 1)
	return tm1.Unix(), tm2.Unix()
}

// LastDayOfTimeStamp : 获取本日最后一天的时间戳
func LastDayOfTimeStamp(d time.Time) int64 {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTimeStamp(d).Unix()
}

// GetZeroTimeStamp  获取某一天的0点时间
func GetZeroTimeStamp(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// GetLastTimeStamp : 获取某一天的最后时间
func GetLastTimeStamp(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location())
}

// BeautyTimeStamp : 美化时间
func BeautyTimeStamp(timeStamp, currentTime int64) string {
	if currentTime == 0 {
		currentTime = time.Now().Unix()
	}

	span := currentTime - timeStamp
	if span < 60 {

		return "刚刚"

	} else if span < 3600 {

		tmp := span / 60

		return typeHelper.Int64ToStr(tmp) + "分钟前"

	} else if span < 24*3600 {

		tmp := span / 3600

		return typeHelper.Int64ToStr(tmp) + "小时前"

	} else if span < (7 * 24 * 3600) {

		tmp := span / 24 * 3600

		return typeHelper.Int64ToStr(tmp) + "天前"

	} else {

		tm := time.Unix(timeStamp, 0)

		return tm.Format(HyphenTimeStr)

	}

}

// GetISO8601 获取iso8601格式的时间
func GetISO8601(date int64) string {
	var formattedDate = time.Unix(date, 0).Format(ISO8601TimeStr)
	return formattedDate
}

type MyJsonTime time.Time

// ParseFromDate 从DateStr 专程MyJsonTime
func ParseFromDate(dateStr string) MyJsonTime {
	local, _ := time.LoadLocation("Local")
	temp, _ := time.ParseInLocation(HyphenDateStr, dateStr, local)
	return MyJsonTime(temp)
}

// GetTimePoint 实现它的json序列化方法 注意测试
func (tm MyJsonTime) GetTimePoint() *time.Time {
	//temp,_:= time.Parse("2006-01-02 15:04:05",tm.GetHyphenDateStrStr())
	local, _ := time.LoadLocation("Local")
	temp, _ := time.ParseInLocation("2006-01-02 15:04:05", tm.GetHyphenDateStrStr(), local)

	return &temp
}

// MarshalJSON 实现它的json序列化方法
func (tm MyJsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(tm).Format(HyphenTimeStr))
	return []byte(stamp), nil
}

// GetHyphenDateStrStr 实现它的json序列化方法
func (tm MyJsonTime) GetHyphenDateStrStr() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(HyphenTimeStr)
}

// GetSlashTimeStr 获取 斜杠 日期+时间 字符串
func (tm MyJsonTime) GetSlashTimeStr() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(SlashTimeStr)
}

// GetDateStr 获取短横线 日期 字符串
func (tm MyJsonTime) GetDateStr() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(HyphenDateStr)
}

// GetSlashDateStr 获取 斜杠 日期 字符串
func (tm MyJsonTime) GetSlashDateStr() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(SlashDateStr)
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

func (tm MyJsonTime) Add(duration time.Duration) MyJsonTime {
	temp := time.Time(tm)
	temp = temp.Add(duration)
	return MyJsonTime(temp)
}

func (tm MyJsonTime) Day() int {
	temp := time.Time(tm)
	return temp.Day()
}

func (tm MyJsonTime) Hour() int {
	temp := time.Time(tm)
	return temp.Hour()
}

func (tm MyJsonTime) Minute() int {
	temp := time.Time(tm)
	return temp.Minute()
}

func (tm MyJsonTime) Second() int {
	temp := time.Time(tm)
	return temp.Second()
}

// WeekdayStrCn 返回中文星期
func (tm MyJsonTime) WeekdayStrCn(style int) string {
	temp := time.Time(tm)
	return weekdayCn[style][int(temp.Weekday())]
}

// WeekdayStrEn 返回英文星期
func (tm MyJsonTime) WeekdayStrEn() string {
	temp := time.Time(tm)
	return temp.Weekday().String()
}

func (tm MyJsonTime) WeekdayInt() int {
	temp := time.Time(tm)
	return int(temp.Weekday())
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

// GetCST8Now 获取东八区现在的时间
func GetCST8Now() time.Time {
	return time.Now().UTC().Add(8 * time.Hour)
}

// GetFirstDateOfWeek 获取本周周一的日期
func GetFirstDateOfWeek(target time.Time) (weekStartDate time.Time) {

	offset := int(time.Monday - target.Weekday())
	if offset > 0 {
		offset = -6
	}

	return time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}

// GetLastDateOfWeek 获取本周周日的日期
func GetLastDateOfWeek(target time.Time) (weekStartDate time.Time) {

	offset := int(time.Saturday - target.Weekday())
	if offset > 6 {
		offset = -1
	}

	return time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}

// GetMonthDayStr 获得当前月的初始和结束日期
func GetMonthDayStr(target time.Time) (string, string) {

	currentYear, currentMonth, _ := target.Date()
	currentLocation := target.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	f := firstOfMonth.Unix()
	l := lastOfMonth.Unix()
	return time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format(HyphenDateStr) + " 23:59:59"
}

// GetWeekDayStr 获得当前周的初始和结束日期
func GetWeekDayStr(target time.Time) (string, string) {

	offset := int(time.Monday - target.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	lastOffset := int(time.Saturday - target.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if lastOffset == 6 {
		lastOffset = -1
	}

	firstOfWeek := time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	lastOfWeeK := time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, lastOffset+1)
	f := firstOfWeek.Unix()
	l := lastOfWeeK.Unix()
	return time.Unix(f, 0).Format(HyphenDateStr) + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
}

// GetQuarterDayStr 获得当前季度的初始和结束日期
func GetQuarterDayStr(target time.Time) (string, string) {
	year := target.Format("2006")
	month := int(target.Month())
	var firstOfQuarter string
	var lastOfQuarter string
	if month >= 1 && month <= 3 {
		//1月1号
		firstOfQuarter = year + "-01-01 00:00:00"
		lastOfQuarter = year + "-03-31 23:59:59"
	} else if month >= 4 && month <= 6 {
		firstOfQuarter = year + "-04-01 00:00:00"
		lastOfQuarter = year + "-06-30 23:59:59"
	} else if month >= 7 && month <= 9 {
		firstOfQuarter = year + "-07-01 00:00:00"
		lastOfQuarter = year + "-09-30 23:59:59"
	} else {
		firstOfQuarter = year + "-10-01 00:00:00"
		lastOfQuarter = year + "-12-31 23:59:59"
	}
	return firstOfQuarter, lastOfQuarter
}

// GetBetweenDateStrs 根据开始日期和结束日期计算出时间段内所有日期
// 参数为日期格式，如：2020-01-01
func GetBetweenDateStrs(startDate, endDate string) []string {
	d := []string{}
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(startDate) {
		timeFormatTpl = timeFormatTpl[0:len(startDate)]
	}
	date, err := time.Parse(timeFormatTpl, startDate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, endDate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}
