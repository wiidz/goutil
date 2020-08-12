package goutil

import (
	"time"
)

type TimeHelper struct{}

/**
 * @func: FromTodayToTomorrowTimeStamp  返回今天凌晨和明天凌晨的时间戳
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TimeHelper) FromTodayToTomorrowTimeStamp() (today, tomorrow int64) {
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
func (*TimeHelper) LastDayOfTimeStamp(d time.Time) int64 {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTimeStamp(d).Unix()
}

/**
 * @func: GetZeroTimeStamp  获取某一天的0点时间
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TimeHelper) GetZeroTimeStamp(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

/**
 * @func: BeautyTimeStamp 美化时间
 * @author Wiidz
 * @date   2019-11-16
 */

func (*TimeHelper) BeautyTimeStamp(timeStamp, currentTime int64) string {
	if currentTime == 0 {
		currentTime = time.Now().Unix()
	}

	span := currentTime - timeStamp
	var typeHelper TypeHelper
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
func (*TimeHelper) GetISO8601(date int64) string {
	var formattedDate = time.Unix(date, 0).Format("2006-01-02T15:04:05Z")
	return formattedDate
}
