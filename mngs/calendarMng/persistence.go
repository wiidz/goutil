package calendarMng

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CalendarDay struct {
	Date       string `gorm:"size:8;uniqueIndex;not null" json:"date"` // 公历日期，格式：YYYYMMDD
	SimpleMark string `gorm:"size:64" json:"simple_mark"`              // 简单标记（法定节假日 > 节气 > 农历简写），用于显示在日历上

	WeekDay   int    `gorm:"not null" json:"week_day"`         // 星期几的数字，0-6（0为周日），常规周末都是休息日
	LunarYear string `gorm:"size:16" json:"lunar_year"`        // 农历年份（乙巳年）
	LunarDate string `gorm:"size:32" json:"lunar_date"`        // 农历日期（十一月十一）
	Shengxiao string `gorm:"size:8" json:"shengxiao"`          // 蛇
	Jieqi24   string `gorm:"size:16" json:"jieqi24,omitempty"` // 二十四节气

	Xingzuo string `gorm:"size:16" json:"xingzuo"` // 星座

	Festival string `gorm:"size:32" json:"festival,omitempty"` // 当天的节日名称（除了法定节日外，植树节、青年节等）

	Holiday         string `gorm:"size:32" json:"holiday,omitempty"`                // 假日名称（元旦节，春节，清明节，劳动节，端午节，国庆节、中秋节）
	IsRest          bool   `gorm:"not null;default:false" json:"is_rest"`           // 是否是休息日（包含法定节假日和周末）
	IsAdjustWorkday bool   `gorm:"not null;default:false" json:"is_adjust_workday"` // 是否是调休工作日

	Yi string `gorm:"type:text" json:"yi,omitempty"` // 宜
	Ji string `gorm:"type:text" json:"ji,omitempty"` // 忌
}

// CalendarDayLite 用于日历展示的简化数据
type CalendarDayLite struct {
	Date            string `json:"date"`               // YYYYMMDD
	SimpleMark      string `json:"simple_mark"`        // 显示标记：法定假日/节气/农历简写
	WeekDay         int    `json:"week_day"`           // 0=周日
	Festival        string `json:"festival,omitempty"` // 纪念日/节日
	IsRest          bool   `json:"is_rest"`            // 是否休息日
	IsAdjustWorkday bool   `json:"is_adjust_workday"`  // 是否调休

	// LunarDate       string `json:"lunar_date"`         // 农历日期（简写）
	// Shengxiao       string `json:"shengxiao"`          // 生肖
	// Jieqi24         string `json:"jieqi24,omitempty"`  // 节气
	// Holiday         string `json:"holiday,omitempty"`  // 假日名称
}

type CalendarDayRepo struct {
	db    *gorm.DB
	table string
}

// NewCalendarDayRepo 允许自定义表名，默认 "a_calendar"
func NewCalendarDayRepo(db *gorm.DB, tableName ...string) *CalendarDayRepo {
	t := "a_calendar"
	if len(tableName) > 0 && tableName[0] != "" {
		t = tableName[0]
	}
	return &CalendarDayRepo{db: db, table: t}
}

func (r *CalendarDayRepo) AutoMigrate(ctx context.Context) error {
	return r.db.WithContext(ctx).Table(r.table).AutoMigrate(&CalendarDay{})
}

// EnsureTable 确保表存在（默认表名 a_calendar)
func (r *CalendarDayRepo) EnsureTable(ctx context.Context) error {
	return r.AutoMigrate(ctx)
}

// Upsert 按 date 唯一键插入或更新
func (r *CalendarDayRepo) Upsert(ctx context.Context, day *CalendarDay) error {
	return r.db.WithContext(ctx).
		Table(r.table).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "date"}},
			UpdateAll: true,
		}).
		Create(day).Error
}

// BatchUpsert 批量插入/更新
func (r *CalendarDayRepo) BatchUpsert(ctx context.Context, days []*CalendarDay) error {
	if len(days) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Table(r.table).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "date"}},
			UpdateAll: true,
		}).
		Create(&days).Error
}

func (r *CalendarDayRepo) GetByDate(ctx context.Context, date string) (*CalendarDay, error) {
	var day CalendarDay
	err := r.db.WithContext(ctx).Table(r.table).Where("date = ?", date).First(&day).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &day, err
}

// ListByRange 按日期区间查询 [start,end]，按日期升序
func (r *CalendarDayRepo) ListByRange(ctx context.Context, startDate, endDate string) ([]CalendarDay, error) {
	var days []CalendarDay
	err := r.db.WithContext(ctx).
		Table(r.table).
		Where("date >= ? AND date <= ?", startDate, endDate).
		Order("date ASC").
		Find(&days).Error
	return days, err
}

func (r *CalendarDayRepo) DeleteByDate(ctx context.Context, date string) error {
	return r.db.WithContext(ctx).Table(r.table).Where("date = ?", date).Delete(&CalendarDay{}).Error
}

// SyncRangeFromAPI 调用 calendarMng 接口获取 [startDate,endDate]（YYYYMMDD）数据并落库
func (r *CalendarDayRepo) SyncRangeFromAPI(ctx context.Context, mng *CalendarMng, startDate, endDate string) error {
	start, err := time.Parse("20060102", startDate)
	if err != nil {
		return err
	}
	end, err := time.Parse("20060102", endDate)
	if err != nil {
		return err
	}
	if end.Before(start) {
		return errors.New("endDate before startDate")
	}

	var toSave []*CalendarDay
	for curr := start; !curr.After(end); curr = curr.AddDate(0, 0, 1) {
		dateStr := curr.Format("20060102")

		holidayDetail, _ := mng.GetHolidayDetail(dateStr, "1")
		almanac, _ := mng.GetAlmanac(dateStr)

		weekDay := int(curr.Weekday())
		if holidayDetail != nil && holidayDetail.WeekDay >= 0 {
			weekDay = holidayDetail.WeekDay
		}

		lunarYear := ""
		lunarDate := ""
		shengxiao := ""
		xingzuo := ""
		jieqi24 := ""
		yi := ""
		ji := ""

		if almanac != nil {
			lunarYear = almanac.Ganzhi
			lunarDate = almanac.Nongli
			shengxiao = almanac.Shengxiao
			xingzuo = almanac.Xingzuo
			jieqi24 = almanac.Jieqi24
			if almanac.Yi != "" {
				yi = strings.Join(splitClean(almanac.Yi), "、")
			}
			if almanac.Ji != "" {
				ji = strings.Join(splitClean(almanac.Ji), "、")
			}
		}

		holiday := ""
		holidayRemark := ""
		isRest := weekDay == 0 || weekDay == 6
		if holidayDetail != nil {
			holiday = holidayDetail.Holiday
			holidayRemark = holidayDetail.HolidayRemark
			isRest = holidayDetail.Type == "2" || holidayDetail.Type == "3" || isRest
		}
		isAdjust := strings.Contains(holidayRemark, "调休")

		festival := ""
		if almanac != nil && almanac.Jieri != "" {
			festival = almanac.Jieri
		}

		simpleMark := firstNonEmpty(
			nonWorkHoliday(holiday),
			festival,
			jieqi24,
			shortLunar(lunarDate),
		)

		day := &CalendarDay{
			Date:            dateStr,
			SimpleMark:      simpleMark,
			WeekDay:         weekDay,
			LunarYear:       lunarYear,
			LunarDate:       lunarDate,
			Shengxiao:       shengxiao,
			Jieqi24:         jieqi24,
			Xingzuo:         xingzuo,
			Festival:        festival,
			Holiday:         holiday,
			IsRest:          isRest,
			IsAdjustWorkday: isAdjust,
			Yi:              yi,
			Ji:              ji,
		}
		toSave = append(toSave, day)
	}

	return r.BatchUpsert(ctx, toSave)
}

// helper: 拆分宜/忌文本
func splitClean(s string) []string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '，' || r == '、' || r == ',' || r == ' ' || r == '　'
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// helper: 如果 holiday 不是“无/周末”则返回，否则空
func nonWorkHoliday(h string) string {
	if h == "" || h == "无" || h == "周末" {
		return ""
	}
	return h
}

// helper: 返回第一个非空字符串
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// helper: 从农历中取简写，如 “二零二五年十一月(大)初二 今日冬至” -> “十一月初二”
func shortLunar(nongli string) string {
	if nongli == "" {
		return ""
	}
	parts := strings.FieldsFunc(nongli, func(r rune) bool {
		return r == ' ' || r == '　'
	})
	base := parts[0]
	if idx := strings.LastIndex(base, "年"); idx >= 0 && idx+len("年") < len(base) {
		return strings.Trim(base[idx+len("年"):], "()")
	}
	return base
}
