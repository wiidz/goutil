package appMng

import "github.com/wiidz/goutil/helpers/loggerHelper"

type AppLogger interface {
	Build() error
	GetDefault() *loggerHelper.LoggerHelper
}

type testLogger struct {
	Client  *loggerHelper.LoggerHelper
	Console *loggerHelper.LoggerHelper
	Admin   *loggerHelper.LoggerHelper
}

func (l *testLogger) Build() error {
	return nil
}

func (l *testLogger) GetDefault() *loggerHelper.LoggerHelper {
	return l.Client
}
