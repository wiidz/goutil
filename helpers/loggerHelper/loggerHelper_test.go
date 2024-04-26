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
		ShowFileAndLine: true,
		Json:            false,
		Level:           zapcore.DebugLevel,
		SyncToConsole:   true,
	})
	temp.Info("test")
}
