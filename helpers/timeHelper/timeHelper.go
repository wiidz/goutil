package timeHelper

import (
	"fmt"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"time"
)

const (
	TimeStr        = "15:04:05"
	HyphenTimeStr  = "2006-01-02 15:04:05"
	SlashTimeStr   = "2006/01/02 15:04:05"
	HyphenDateStr  = "2006-01-02"
	SlashDateStr   = "2006/01/02"
	ISO8601TimeStr = "2006-01-02T15:04:05Z"
	PureNumber     = "20060102150405"
	PureNumberDate = "20060102"

	HyphenTimeStrNoYear = "01-02 15:04:05"
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

// ParseFromHyphenDateStr 从DateStr 转成 MyJsonTime "2006-01-02"
func ParseFromHyphenDateStr(hyphenDateStr string) (newTime MyJsonTime, err error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		return
	}

	temp, err := time.ParseInLocation(HyphenDateStr, hyphenDateStr, local)
	if err != nil {
		return
	}

	newTime = MyJsonTime(temp)
	return
}

// ParseFromPrueNumberDateTimeStr 从DateStr 转成 MyJsonTime "20060102"
func ParseFromPrueNumberDateTimeStr(timeStr string) (newTime MyJsonTime, err error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		return
	}

	temp, err := time.ParseInLocation(PureNumberDate, timeStr, local)
	if err != nil {
		return
	}

	newTime = MyJsonTime(temp)
	return
}

// ParseFromHyphenTimeStr 从DateStr 转成 MyJsonTime "2006-01-02 15:04:05"
func ParseFromHyphenTimeStr(hyphenTimeStr string) (newTime MyJsonTime, err error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		return
	}

	temp, err := time.ParseInLocation(HyphenTimeStr, hyphenTimeStr, local)
	if err != nil {
		return
	}

	newTime = MyJsonTime(temp)
	return
}

// ParseFromSlashDateStr 从DateStr 转成 MyJsonTime "2006/01/02"
func ParseFromSlashDateStr(slashDateStr string) (newTime MyJsonTime, err error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		return
	}

	temp, err := time.ParseInLocation(SlashDateStr, slashDateStr, local)
	if err != nil {
		return
	}

	newTime = MyJsonTime(temp)
	return
}

// ParseFromSlashTimeStr 从DateTimeStr 转成 MyJsonTime "2006/01/02 15:04:05"
func ParseFromSlashTimeStr(slashTimeStr string) (newTime MyJsonTime, err error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		return
	}
	temp, err := time.ParseInLocation(SlashTimeStr, slashTimeStr, local)
	if err != nil {
		return
	}

	newTime = MyJsonTime(temp)
	return

}

// ParseFromHyphenDateStrWithLocation 从DateStr 转成 MyJsonTime "2006-01-02"
func ParseFromHyphenDateStrWithLocation(hyphenDateStr string, location *time.Location) (jsonTime MyJsonTime, err error) {
	temp, err := time.ParseInLocation(HyphenDateStr, hyphenDateStr, location)
	jsonTime = MyJsonTime(temp)
	return
}

// ParseFromHyphenTimeStrWithLocation 从DateStr 转成 MyJsonTime "2006-01-02 15:04:05"
func ParseFromHyphenTimeStrWithLocation(hyphenTimeStr string, location *time.Location) (jsonTime MyJsonTime, err error) {
	temp, err := time.ParseInLocation(HyphenTimeStr, hyphenTimeStr, location)
	jsonTime = MyJsonTime(temp)
	return
}

// ParseFromSlashDateStrWithLocation 从DateStr 转成 MyJsonTime "2006/01/02"
func ParseFromSlashDateStrWithLocation(slashDateStr string, location *time.Location) (jsonTime MyJsonTime, err error) {
	temp, err := time.ParseInLocation(SlashDateStr, slashDateStr, location)
	jsonTime = MyJsonTime(temp)
	return
}

// ParseFromSlashTimeStrWithLocation 从DateTimeStr 转成 MyJsonTime "2006/01/02 15:04:05"
func ParseFromSlashTimeStrWithLocation(slashTimeStr string, location *time.Location) (jsonTime MyJsonTime, err error) {
	temp, err := time.ParseInLocation(SlashTimeStr, slashTimeStr, location)
	jsonTime = MyJsonTime(temp)
	return
}

// GetHyphenDateStr 获取 短横线 日期 字符串
func (tm MyJsonTime) GetHyphenDateStr() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(HyphenDateStr)
}

// GetHyphenTimeStr 获取 短横线 日期+时间 字符串
func (tm MyJsonTime) GetHyphenTimeStr() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(HyphenTimeStr)
}

// GetHyphenTimeStrNoYear 如果今年是今年，去掉今年（哲学）
func (tm MyJsonTime) GetHyphenTimeStrNoYear(yearNow int) string {
	if tm.IsNull() {
		return ""
	}

	if tm.Year() == yearNow {
		return time.Time(tm).Format(HyphenTimeStrNoYear)
	} else {
		return time.Time(tm).Format(HyphenTimeStr)
	}
}

// GetTimeStr 获取 时间 字符串
func (tm MyJsonTime) GetTimeStr() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(TimeStr)
}

// GetPureNumberStr 获取 纯数字 年月日时分秒
func (tm MyJsonTime) GetPureNumberStr() string {
	if tm.IsNull() {
		return ""
	}
	return time.Time(tm).Format(PureNumber)
}

// GetTimePoint 获取时间指针类型
func (tm MyJsonTime) GetTimePoint() *time.Time {
	//temp,_:= time.Parse("2006-01-02 15:04:05",tm.GetHyphenDateStrStr())
	local, _ := time.LoadLocation("Local")
	temp, _ := time.ParseInLocation(HyphenTimeStr, tm.GetHyphenDateStr(), local)

	return &temp
}

// MarshalJSON 实现它的json序列化方法
func (tm MyJsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(tm).Format(HyphenTimeStr))
	return []byte(stamp), nil
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

// After 在目标时间之后
func (tm MyJsonTime) After(target MyJsonTime) bool {
	temp := time.Time(tm)
	targetTemp := time.Time(target)
	return temp.After(targetTemp)
}
func (tm MyJsonTime) AfterTime(target time.Time) bool {
	temp := time.Time(tm)
	return temp.After(target)
}
func (tm MyJsonTime) AfterOrEqual(target MyJsonTime) bool {
	temp := time.Time(tm)
	targetTemp := time.Time(target)
	return temp.After(targetTemp) || temp.Equal(targetTemp)
}
func (tm MyJsonTime) AfterOrEqualTime(target time.Time) bool {
	temp := time.Time(tm)
	return temp.After(target) || temp.Equal(target)
}

// Before 在目标时间之前
func (tm MyJsonTime) Before(target MyJsonTime) bool {
	temp := time.Time(tm)
	targetTemp := time.Time(target)
	return temp.Before(targetTemp)
}
func (tm MyJsonTime) BeforeTime(target time.Time) bool {
	temp := time.Time(tm)
	return temp.Before(target)
}
func (tm MyJsonTime) BeforeOrEqual(target MyJsonTime) bool {
	temp := time.Time(tm)
	targetTemp := time.Time(target)
	return temp.Before(targetTemp) || temp.Equal(targetTemp)
}
func (tm MyJsonTime) BeforeOrEqualTime(target time.Time) bool {
	temp := time.Time(tm)
	return temp.Before(target) || temp.Equal(target)
}

// Equal 与目标时间相同
func (tm MyJsonTime) Equal(target MyJsonTime) bool {
	temp := time.Time(tm)
	targetTemp := time.Time(target)
	return temp.Equal(targetTemp)
}
func (tm MyJsonTime) EqualTime(target time.Time) bool {
	temp := time.Time(tm)
	return temp.Equal(target)
}

// Between 在两个时间中间，注意前面那个一定要小一点
func (tm MyJsonTime) Between(from, to MyJsonTime) bool {
	temp := time.Time(tm)
	fromTemp := time.Time(from)
	toTemp := time.Time(to)
	return temp.After(fromTemp) && temp.Before(toTemp)
}
func (tm MyJsonTime) BetweenTime(from, to time.Time) bool {
	temp := time.Time(tm)
	return temp.After(from) && temp.Before(to)
}
func (tm MyJsonTime) BetweenOrEqual(from, to MyJsonTime) bool {
	temp := time.Time(tm)
	fromTemp := time.Time(from)
	toTemp := time.Time(to)
	return (temp.After(fromTemp) || temp.Equal(fromTemp)) && (temp.Before(toTemp) || temp.Equal(toTemp))
}

func (tm MyJsonTime) BetweenOrEqualTime(from, to time.Time) bool {
	temp := time.Time(tm)
	return (temp.After(from) || temp.Equal(from)) && (temp.Before(to) || temp.Equal(to))
}

func (tm MyJsonTime) Unix() int64 {
	temp := time.Time(tm)
	return temp.Unix()
}

func (tm MyJsonTime) UnixNano() int64 {
	temp := time.Time(tm)
	return temp.UnixNano()
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

// GetFirstTimeOfMonth 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func (tm MyJsonTime) GetFirstTimeOfMonth() *MyJsonTime {
	newTime := tm.AddDate(0, 0, -tm.Day()+1)
	res := newTime.GetZeroTime()
	return &res
}

// GetLastTimeOfMonth 获取传入的时间所在月份的最后天，即某月最后一天的23点59分59秒
func (tm MyJsonTime) GetLastTimeOfMonth() *MyJsonTime {
	newTime := tm.AddDate(0, 1, -tm.Day())
	res := newTime.GetLastTime()
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

// SubIntDays 获取与目标时间相差的整数天
func (tm MyJsonTime) SubIntDays(targetTime time.Time) int {
	days := int(tm.Format2Time().Sub(targetTime).Hours() / 24)
	return days
}

// SubNowIntDays 获取与今天相差的整数天
func (tm MyJsonTime) SubNowIntDays() int {
	targetTime := time.Now()
	days := int(tm.Format2Time().Sub(targetTime).Hours() / 24)
	return days
}

// GetCST8Now 获取东八区现在的时间
func GetCST8Now() time.Time {
	return time.Now().UTC().Add(8 * time.Hour)
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

// GetFirstDateOfWeek 获取本周周一的日期
func (tm MyJsonTime) GetFirstDateOfWeek() (weekStartDate MyJsonTime) {

	temp := tm.Format2Time()
	offset := int(time.Monday - temp.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate = MyJsonTime(time.Date(temp.Year(), temp.Month(), temp.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset))
	return
}

// GetLastDateOfWeek 获取本周周日的日期
func (tm MyJsonTime) GetLastDateOfWeek() (weekStartDate MyJsonTime) {

	// 我们国内一周的开始是周一，国外是周日
	// go里面周日=0，我们改成7

	temp := tm.Format2Time()
	nowWeekDay := int(temp.Weekday())
	if nowWeekDay == 0 {
		nowWeekDay = 7
	}
	offset := int(7 - temp.Weekday())

	weekStartDate = MyJsonTime(time.Date(temp.Year(), temp.Month(), temp.Day(), 23, 59, 59, 1e9-1, time.Local).AddDate(0, 0, offset))
	return
}

func NowHyphenDateStr() string {
	currentTime := time.Now()
	return currentTime.Format(HyphenDateStr)
}

func NowHyphenTimeStr() string {
	currentTime := time.Now()
	return currentTime.Format(HyphenTimeStr)
}

// NowHyphenTimeStrByDays 获取距离今天n天的日期
func NowHyphenTimeStrByDays(diffDays int) string {
	currentTime := time.Now()
	formattedTime := currentTime.AddDate(0, 0, diffDays)
	return formattedTime.Format("2006-01-02")
}
