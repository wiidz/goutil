package loggerHelper

import (
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestNewLoggerHelper(t *testing.T) {

	temp, _ := NewLoggerHelper(&Config{
		Filename: "./test.log",
		//Filename:  "",
		//AddCaller:     false,
		IsFullPath:      false,
		ShowFileAndLine: true,
		Json:            false,
		Level:           zapcore.InfoLevel,
		SyncToConsole:   true,
	})
	temp.Info("test")
}
