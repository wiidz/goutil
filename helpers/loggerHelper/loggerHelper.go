package loggerHelper

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"runtime"
)

// Config 配置项
type Config struct {
	Filename      string        // 输出的log文件路径
	Level         zapcore.Level // 输出限制等级
	Json          bool          // 是否以json格式输出
	SyncToConsole bool          // 是否同步到控制台（仅文本时生效）

	IsFullPath      bool // 输出的复杂度，true的时候只输出文件名，false就输出复杂一点
	AddCaller       bool // 是否输出调用文件和行数，目前这个caller只会输出loggerHelper调用，所以可以不要，我们已经默认输出行数了
	ShowFileAndLine bool // 我们用这个来控制是否输出行数
}

type LoggerHelper struct {
	Config *Config

	// Example 函数用于创建一个示例日志记录器，用于演示和测试目的。它不会输出到任何地方，只是简单地把日志消息打印到控制台。这个函数主要用于演示 zap 日志库的使用方法，以及帮助用户理解日志的格式和输出内容。
	Example *zap.Logger

	// Normal Production 创建一个适用于生产环境的日志记录器。它会使用默认的日志配置，并且性能较高，适用于在生产环境中记录日志。该日志记录器会包含文件名和行号等额外的调用信息，方便排查问题。
	Normal *zap.Logger // 常规

	// Sugar 在性能很好但不是很关键的上下文中，使用SugaredLogger。它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录。
	// 在每一微秒和每一次内存分配都很重要的上下文中，使用Logger。它甚至比SugaredLogger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。
	// 然后使用SugaredLogger以printf格式记录语句
	Sugar *zap.SugaredLogger // 短小精悍

	consoleSugar *zap.SugaredLogger // 主要用于文本log的时候，同步输出到控制台
}

// NewLoggerHelper 返回日志榜首
func NewLoggerHelper(config *Config) (helper *LoggerHelper, err error) {
	helper = &LoggerHelper{
		Config: config,
	}

	var isFileLogger bool
	helper.Normal, isFileLogger, err = getLogger(config.Filename, config)
	if err != nil {
		return
	}

	//【2】如果是文件logger，同步做一个输出到控制台
	if config.SyncToConsole && isFileLogger {
		var tempLogger *zap.Logger
		tempLogger, _, err = getLogger("", config)
		if err != nil {
			return
		}
		helper.consoleSugar = tempLogger.Sugar()
	}

	helper.Sugar = helper.Normal.Sugar()
	return
}

// log 输出信息
// 在 Zap 中，Panic 级别的日志消息确实比 Fatal 级别更严重。这可能有些让人困惑，因为通常我们认为 Panic 表示程序遇到了无法恢复的严重错误，会立即终止程序的执行，而 Fatal 表示程序遇到了严重错误，但还有可能继续执行。
// Zap 设计 Panic 级别的日志消息比 Fatal 更严重，是因为 Panic 级别的日志消息触发了程序的 panic 行为，这可能会导致程序在运行时出现崩溃，因此被视为更为严重的错误。
// 在 Zap 中，Panic 级别的日志消息通常用于标识一些非常严重的、无法恢复的错误，即使这些错误可能并不需要导致程序立即终止。而 Fatal 级别的日志消息则表示程序遇到了严重错误，但仍有可能继续执行，因此其严重性略低于 Panic。
func (helper *LoggerHelper) log(sugar *zap.SugaredLogger, level zapcore.Level, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2) // 获取调用方法的文件名和行号
	if !helper.Config.IsFullPath {
		file = filepath.Base(file)
	}
	msg := fmt.Sprint(args...)
	switch level {
	// 以下也代表了level层级
	case zapcore.PanicLevel:
		sugar.Panicw(msg, "file", file, "line", line) // 会导致程序退出（用于生产环境），当发生严重错误但程序仍然有可能继续执行时，可以选择使用 Panic 级别。记录 Panic 级别的日志消息会导致程序 panic，并终止程序的执行，但与 Fatal 不同的是，Panic 可以提供更多的上下文信息和调试信息，有助于排查问题。
	case zapcore.FatalLevel:
		sugar.Fatalw(msg, "file", file, "line", line) // 会导致程序退出（用于生产环境），当发生严重错误并且程序无法继续执行时，可以选择使用 Fatal 级别。记录 Fatal 级别的日志消息将导致程序立即退出，这在生产环境中可以帮助及时发现并解决一些无法修复的严重问题。
	case zapcore.DPanicLevel:
		sugar.DPanicw(msg, "file", file, "line", line) // 不会退出（用户开发环境）
	case zapcore.ErrorLevel:
		sugar.Errorw(msg, "file", file, "line", line)
	case zapcore.WarnLevel:
		sugar.Warnw(msg, "file", file, "line", line)
	case zapcore.InfoLevel:
		sugar.Infow(msg, "file", file, "line", line)
	case zapcore.DebugLevel:
		sugar.Debugw(msg, "file", file, "line", line)
	}
}

// Info 简单方法，用sugar输出
func (helper *LoggerHelper) Info(args ...interface{}) {
	if helper.Config.ShowFileAndLine {
		// 用补充数据输出行数
		helper.log(helper.Sugar, zapcore.InfoLevel, args...)
		if helper.consoleSugar != nil {
			helper.log(helper.consoleSugar, zapcore.InfoLevel, args...)
		}
	} else {
		// 不用补充数据输出行数
		helper.Sugar.Info(args)
		if helper.consoleSugar != nil {
			helper.consoleSugar.Info(args)
		}
	}
}
func (helper *LoggerHelper) Error(args ...interface{}) {
	if helper.Config.ShowFileAndLine {
		// 用补充数据输出行数
		helper.log(helper.Sugar, zapcore.ErrorLevel, args...)
		if helper.consoleSugar != nil {
			helper.log(helper.consoleSugar, zapcore.ErrorLevel, args...)
		}
	} else {
		// 不用补充数据输出行数
		helper.Sugar.Error(args)
		if helper.consoleSugar != nil {
			helper.consoleSugar.Error(args)
		}
	}
}
func (helper *LoggerHelper) Debug(args ...interface{}) {
	if helper.Config.ShowFileAndLine {
		// 用补充数据输出行数
		helper.log(helper.Sugar, zapcore.DebugLevel, args...)
		if helper.consoleSugar != nil {
			helper.log(helper.consoleSugar, zapcore.DebugLevel, args...)
		}
	} else {
		// 不用补充数据输出行数
		helper.Sugar.Debug(args)
		if helper.consoleSugar != nil {
			helper.consoleSugar.Debug(args)
		}
	}
}
func (helper *LoggerHelper) Fatal(args ...interface{}) {
	if helper.Config.ShowFileAndLine {
		// 用补充数据输出行数
		helper.log(helper.Sugar, zapcore.FatalLevel, args...)
		if helper.consoleSugar != nil {
			helper.log(helper.consoleSugar, zapcore.FatalLevel, args...)
		}
	} else {
		// 不用补充数据输出行数
		helper.Sugar.Fatal(args)
		if helper.consoleSugar != nil {
			helper.consoleSugar.Fatal(args)
		}
	}
}
func (helper *LoggerHelper) Panic(args ...interface{}) {
	if helper.Config.ShowFileAndLine {
		// 用补充数据输出行数
		helper.log(helper.Sugar, zapcore.PanicLevel, args...)
		if helper.consoleSugar != nil {
			helper.log(helper.consoleSugar, zapcore.PanicLevel, args...)
		}
	} else {
		// 不用补充数据输出行数
		helper.Sugar.Panic(args)
		if helper.consoleSugar != nil {
			helper.consoleSugar.Panic(args)
		}
	}
}
func (helper *LoggerHelper) DPanic(args ...interface{}) {
	if helper.Config.ShowFileAndLine {
		// 用补充数据输出行数
		helper.log(helper.Sugar, zapcore.DPanicLevel, args...)
		if helper.consoleSugar != nil {
			helper.log(helper.consoleSugar, zapcore.DPanicLevel, args...)
		}
	} else {
		// 不用补充数据输出行数
		helper.Sugar.DPanic(args)
		if helper.consoleSugar != nil {
			helper.consoleSugar.DPanic(args)
		}
	}
}
func (helper *LoggerHelper) Warn(args ...interface{}) {
	if helper.Config.ShowFileAndLine {
		// 用补充数据输出行数
		helper.log(helper.Sugar, zapcore.WarnLevel, args...)
		if helper.consoleSugar != nil {
			helper.log(helper.consoleSugar, zapcore.WarnLevel, args...)
		}
	} else {
		// 不用补充数据输出行数
		helper.Sugar.Warn(args)
		if helper.consoleSugar != nil {
			helper.consoleSugar.Warn(args)
		}
	}
}

// getLogger 获取logger
func getLogger(fileName string, config *Config) (logger *zap.Logger, isFileLogger bool, err error) {
	var core zapcore.Core
	if fileName != "" {
		// 输出到文本
		encoder := getEncoder(config.Json, false, config.IsFullPath) // 文本中强制不要颜色
		writeSyncer := getLogWriter(fileName)
		core = getCore(encoder, writeSyncer, config.Level)
		isFileLogger = true
	} else {
		// 输出到控制台(json输出到控制台可读性太差了，强制不要)
		var encoder zapcore.Encoder
		if config.Json {
			//【1】如果是json输出就不要颜色
			encoder = getEncoder(false, false, config.IsFullPath) // 控制台输出强制有颜色，json其实没啥意义，先强制不要
		} else {
			//【2】如果是控制台就强制要颜色
			encoder = getEncoder(false, true, config.IsFullPath) // 控制台输出强制有颜色，json其实没啥意义，先强制不要
		}

		//consoleDebugging := zapcore.Lock(zapcore.AddSync(os.Stdout))
		//core = getCore(encoder, consoleDebugging, level) // 直接使用 os.Stdout 输出到控制台
		core = getCore(encoder, os.Stdout, config.Level) // 直接使用 os.Stdout 输出到控制台
	}

	if config.AddCaller {
		logger = zap.New(core, zap.AddCaller())
	} else {
		logger = zap.New(core)
	}
	return
}

func getEncoder(json, color, isFullPath bool) (encoder zapcore.Encoder) {

	var encodeLevel zapcore.LevelEncoder
	if color {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encodeLevel = zapcore.CapitalLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,                    // 可以选择其他编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // 可以选择其他编码器
		EncodeDuration: zapcore.SecondsDurationEncoder, // 可以选择其他编码器
		//EncodeCaller:   zapcore.FullCallerEncoder, // 注意这里设置了以后，外面还是要加addCaller才行
	}

	// 注意：由于我们使用了log方法，不需要输出原本的路径了，因为不管咋输出都是loggerHelper，但是当外面使用helper.Sugar输出就有效
	if isFullPath {
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	} else {
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	if json {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // json格式
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return
}

func getLogWriter(filename string) zapcore.WriteSyncer {
	l := &lumberjack.Logger{
		Filename:   filename, // 指定新的文件名
		MaxSize:    10,       // 每个日志文件的大小限制，单位MB
		MaxBackups: 3,        // 保留旧日志文件的个数
		MaxAge:     28,       // 文件最多保存多少天
		Compress:   true,     // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(l)
}

func getCore(encoder zapcore.Encoder, write zapcore.WriteSyncer, zapLevel zapcore.Level) (core zapcore.Core) {
	level := zap.NewAtomicLevelAt(zapLevel)
	core = zapcore.NewCore(encoder, write, level)
	return
}
