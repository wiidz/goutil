package gormZapLogger

import (
	"context"
	"errors"
	"fmt"
	"github.com/wiidz/goutil/helpers/loggerHelper"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

// GormZapLogger 此包为重写logger中的方法，以适用于gorm使用
// 是一个新的结构体，嵌入了 example 包中的 StructA 结构体
type GormZapLogger struct {
	loggerHelper.LoggerHelper // 嵌入 example 包中的 StructA 结构体

	Config                              gormLogger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *GormZapLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.Config.LogLevel = level
	return &newLogger
}

// Info print info
func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= gormLogger.Info {
		l.Sugar.Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= gormLogger.Warn {
		l.Sugar.Warnf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.Config.LogLevel >= gormLogger.Error {
		l.Sugar.Errorf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
//
//nolint:cyclop
func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.Config.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.Config.LogLevel >= gormLogger.Error && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.Config.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			//l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			l.Sugar.Infof(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			//l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			l.Sugar.Infof(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.Config.SlowThreshold && l.Config.SlowThreshold != 0 && l.Config.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.Config.SlowThreshold)
		if rows == -1 {
			l.Sugar.Infof(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Sugar.Infof(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.Config.LogLevel == gormLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Sugar.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Sugar.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
