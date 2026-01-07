package appMng

import (
	"github.com/wiidz/goutil/helpers/gormZapLogger"
	"github.com/wiidz/goutil/helpers/loggerHelper"
)

type AppLogger interface {
	Build() error
	GetDefault() *loggerHelper.LoggerHelper
	GetGorm() *gormZapLogger.GormZapLogger
}

type testLogger struct {
	Client  *loggerHelper.LoggerHelper
	Console *loggerHelper.LoggerHelper
	Admin   *loggerHelper.LoggerHelper
	Gorm    *gormZapLogger.GormZapLogger
}

func (l *testLogger) Build() error {
	return nil
}

func (l *testLogger) GetDefault() *loggerHelper.LoggerHelper {
	return l.Client
}

func (l *testLogger) GetGorm() *gormZapLogger.GormZapLogger {
	return l.Gorm
}
