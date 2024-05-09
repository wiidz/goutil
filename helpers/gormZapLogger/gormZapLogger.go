package gormZapLogger

import (
	"context"
	"errors"
	"fmt"
	"github.com/wiidz/goutil/helpers/loggerHelper"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"log"
	"time"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// GormZapLogger 此包为重写logger中的方法，以适用于gorm使用
// 是一个新的结构体，嵌入了 example 包中的 StructA 结构体
type GormZapLogger struct {
	*loggerHelper.LoggerHelper // 嵌入 example 包中的 StructA 结构体

	GormConfig                          *gormLogger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewGormZapLogger(config *loggerHelper.Config, gormConfig *gormLogger.Config) (helper *GormZapLogger, err error) {
	helper = &GormZapLogger{
		GormConfig: gormConfig,
	}
	helper.LoggerHelper, err = loggerHelper.NewLoggerHelper(config)

	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if gormConfig.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	//helper.Writer = writer
	helper.Config = config
	helper.infoStr = infoStr
	helper.warnStr = warnStr
	helper.errStr = errStr
	helper.traceStr = traceStr
	helper.traceWarnStr = traceWarnStr
	helper.traceErrStr = traceErrStr

	return
}

// LogMode log mode
func (l *GormZapLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.GormConfig.LogLevel = level
	return &newLogger
}

// Info print info
func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {

	log.Println("msg", msg)
	log.Println("data", data)

	if l.GormConfig.LogLevel >= gormLogger.Info {
		l.LoggerHelper.Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.GormConfig.LogLevel >= gormLogger.Warn {
		l.LoggerHelper.Warnf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.GormConfig.LogLevel >= gormLogger.Error {
		l.LoggerHelper.Errorf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
//
//nolint:cyclop
func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.GormConfig.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.GormConfig.LogLevel >= gormLogger.Error && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.GormConfig.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.LoggerHelper.Infof(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.LoggerHelper.Infof(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.GormConfig.SlowThreshold && l.GormConfig.SlowThreshold != 0 && l.GormConfig.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.GormConfig.SlowThreshold)
		if rows == -1 {
			l.LoggerHelper.Infof(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.LoggerHelper.Infof(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.GormConfig.LogLevel == gormLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.LoggerHelper.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.LoggerHelper.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
