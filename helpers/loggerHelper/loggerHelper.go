package loggerHelper

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 配置项
type Config struct {
	Filename      string              // 输出的log文件路径
	Level         zapcore.Level       // 输出限制等级
	Json          bool                // 是否以json格式输出
	SyncToConsole bool                // 是否同步到控制台（仅文本时生效）
	EncodeTime    zapcore.TimeEncoder // 时间格式 如 zapcore.ISO8601TimeEncoder

	IsFullPath      bool // 是否输出完整路径：true=完整路径，false=仅文件名
	AddCaller       bool // 是否输出调用文件和行数，目前这个caller只会输出loggerHelper调用，所以可以不要，我们已经默认输出行数了
	ShowFileAndLine bool // 我们用这个来控制是否输出行数

	// 轮转配置（针对文件输出）
	MaxSize    int  // 单个文件最大MB，默认10
	MaxBackups int  // 保留个数，默认3
	MaxAge     int  // 保留天数，默认28
	Compress   bool // 是否压缩，默认true
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

	// 统一通过 zapcore.NewTee 输出到多个目标
}

// MyTimeEncoder 自定义的时间encoder
func MyTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05"))
}

// NewLoggerHelper 返回日志榜首
func NewLoggerHelper(config *Config) (helper *LoggerHelper, err error) {
	helper = &LoggerHelper{
		Config: config,
	}

	// 默认时间编码
	if helper.Config.EncodeTime == nil {
		helper.Config.EncodeTime = MyTimeEncoder
	}

	// 构建 cores（文件 + 控制台）
	var cores []zapcore.Core

	if config.Filename != "" {
		fileEncoder := getEncoder(config, false, config.IsFullPath)
		fileWrite := getLogWriter(config.Filename, config)
		cores = append(cores, getCore(fileEncoder, fileWrite, config.Level))
	}

	if config.Filename == "" || config.SyncToConsole {
		// 控制台强制使用 ConsoleEncoder 且彩色、非JSON
		consCfg := *config
		consCfg.Json = false
		consoleEncoder := getEncoder(&consCfg, true, consCfg.IsFullPath)
		cores = append(cores, getCore(consoleEncoder, zapcore.AddSync(os.Stdout), consCfg.Level))
	}

	// 合并核心
	var core zapcore.Core
	if len(cores) == 1 {
		core = cores[0]
	} else {
		core = zapcore.NewTee(cores...)
	}

	// Caller 配置：当需要显示文件行号时开启；Skip(1) 指向真实调用者
	if config.AddCaller || config.ShowFileAndLine {
		helper.Normal = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		helper.Normal = zap.New(core)
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

// log 输出信息
// 在 Zap 中，Panic 级别的日志消息确实比 Fatal 级别更严重。这可能有些让人困惑，因为通常我们认为 Panic 表示程序遇到了无法恢复的严重错误，会立即终止程序的执行，而 Fatal 表示程序遇到了严重错误，但还有可能继续执行。
// Zap 设计 Panic 级别的日志消息比 Fatal 更严重，是因为 Panic 级别的日志消息触发了程序的 panic 行为，这可能会导致程序在运行时出现崩溃，因此被视为更为严重的错误。
// 在 Zap 中，Panic 级别的日志消息通常用于标识一些非常严重的、无法恢复的错误，即使这些错误可能并不需要导致程序立即终止。而 Fatal 级别的日志消息则表示程序遇到了严重错误，但仍有可能继续执行，因此其严重性略低于 Panic。
func (helper *LoggerHelper) logf(sugar *zap.SugaredLogger, level zapcore.Level, template string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2) // 获取调用方法的文件名和行号
	if !helper.Config.IsFullPath {
		file = filepath.Base(file)
	}
	msg := fmt.Sprintf(template, args...)
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
func (helper *LoggerHelper) Info(args ...interface{})   { helper.Sugar.Info(args...) }
func (helper *LoggerHelper) Error(args ...interface{})  { helper.Sugar.Error(args...) }
func (helper *LoggerHelper) Debug(args ...interface{})  { helper.Sugar.Debug(args...) }
func (helper *LoggerHelper) Fatal(args ...interface{})  { helper.Sugar.Fatal(args...) }
func (helper *LoggerHelper) Panic(args ...interface{})  { helper.Sugar.Panic(args...) }
func (helper *LoggerHelper) DPanic(args ...interface{}) { helper.Sugar.DPanic(args...) }
func (helper *LoggerHelper) Warn(args ...interface{})   { helper.Sugar.Warn(args...) }

func (helper *LoggerHelper) Infof(template string, args ...interface{}) {
	helper.Sugar.Infof(template, args...)
}
func (helper *LoggerHelper) Errorf(template string, args ...interface{}) {
	helper.Sugar.Errorf(template, args...)
}
func (helper *LoggerHelper) Debugf(template string, args ...interface{}) {
	helper.Sugar.Debugf(template, args...)
}
func (helper *LoggerHelper) Fatalf(template string, args ...interface{}) {
	helper.Sugar.Fatalf(template, args...)
}
func (helper *LoggerHelper) Panicf(template string, args ...interface{}) {
	helper.Sugar.Panicf(template, args...)
}
func (helper *LoggerHelper) DPanicf(template string, args ...interface{}) {
	helper.Sugar.DPanicf(template, args...)
}
func (helper *LoggerHelper) Warnf(template string, args ...interface{}) {
	helper.Sugar.Warnf(template, args...)
}

// GetLogger 获取logger
// gormZapLogger 也要用
func GetLogger(fileName string, config *Config) (logger *zap.Logger, isFileLogger bool, err error) {
	var core zapcore.Core
	if fileName != "" {
		// 输出到文本
		encoder := getEncoder(config, false, config.IsFullPath) // 文件是否 JSON 取决于 config.Json
		writeSyncer := getLogWriter(fileName, config)
		core = getCore(encoder, writeSyncer, config.Level)
		isFileLogger = true
	} else {
		// 输出到控制台(json输出到控制台可读性太差了，强制不要)
		consCfg := *config
		consCfg.Json = false
		encoder := getEncoder(&consCfg, true, consCfg.IsFullPath) // 控制台强制彩色文本

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

func getEncoder(config *Config, color, isFullPath bool) (encoder zapcore.Encoder) {

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
		EncodeTime:     config.EncodeTime,              // 可以选择其他编码器
		EncodeDuration: zapcore.SecondsDurationEncoder, // 可以选择其他编码器
		//EncodeCaller:   zapcore.FullCallerEncoder, // 注意这里设置了以后，外面还是要加addCaller才行，但是我们自定义了输出方法，所以不需要了
	}

	// 注意：由于我们使用了log方法，不需要输出原本的路径了，因为不管咋输出都是loggerHelper，但是当外面使用helper.Sugar输出就有效
	if isFullPath {
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	} else {
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	if config.Json {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // json格式
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return
}

func getLogWriter(filename string, cfg *Config) zapcore.WriteSyncer {
	maxSize := cfg.MaxSize
	if maxSize <= 0 {
		maxSize = 10
	}
	maxBackups := cfg.MaxBackups
	if maxBackups <= 0 {
		maxBackups = 3
	}
	maxAge := cfg.MaxAge
	if maxAge <= 0 {
		maxAge = 28
	}
	compress := cfg.Compress
	if !cfg.Compress {
		compress = true
	}
	l := &lumberjack.Logger{
		Filename:   filename, // 指定新的文件名
		MaxSize:    maxSize,  // 每个日志文件的大小限制，单位MB
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	return zapcore.AddSync(l)
}

func getCore(encoder zapcore.Encoder, write zapcore.WriteSyncer, zapLevel zapcore.Level) (core zapcore.Core) {
	level := zap.NewAtomicLevelAt(zapLevel)
	core = zapcore.NewCore(encoder, write, level)
	return
}

// Close 刷盘
func (helper *LoggerHelper) Close() { _ = helper.Normal.Sync() }
