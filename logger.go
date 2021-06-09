package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"reflect"
)

var logger *zap.Logger

type LoggerConfig struct {
	DefaultPath string
	DebugPath   string
	InfoPath    string
	WarnPath    string
	ErrorPath   string
	DPanicPath  string
	PanicPath   string
	FatalPath   string
}

var fields = []string{
	"DebugPath",
	"InfoPath",
	"WarnPath",
	"ErrorPath",
	"DPanicPath",
	"PanicPath",
	"FatalPath",
}

var str2level = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

var levels = []zapcore.Level{
	zapcore.DebugLevel,
	zapcore.InfoLevel,
	zapcore.WarnLevel,
	zapcore.ErrorLevel,
	zapcore.DPanicLevel,
	zapcore.PanicLevel,
	zapcore.FatalLevel,
}

func newHookedCore(filePath string, enab zap.LevelEnablerFunc) zapcore.Core {
	hook := lumberjack.Logger{
		Filename:   filePath, // 日志文件路径
		MaxSize:    128,      // 日志文件最大为128MB
		MaxBackups: 3,        // 最多保留3个备份
		MaxAge:     7,        // 最多保留7天
		Compress:   true,     // 是否压缩 disabled by default
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "lineNum",
		MessageKey:     "msg",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	} // 格式

	//l := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl == level
	//})

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)),
		enab,
	)

	return core
}

func Init(filePath, logLevel string) {

	level, ok := str2level[logLevel]
	if !ok {
		level = zapcore.InfoLevel
	}

	caller := zap.AddCaller()
	development := zap.Development()

	enab := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	logger = zap.New(newHookedCore(filePath, enab), caller, development)

	zap.ReplaceGlobals(logger)
}

func Init2(config *LoggerConfig) {

	var cores []zapcore.Core

	r := reflect.ValueOf(config)
	for idx, field := range fields {
		fieldValue := reflect.Indirect(r).FieldByName(field).String()

		l := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == levels[idx]
		})

		if fieldValue != "" {
			cores = append(cores, newHookedCore(fieldValue, l))
		} else {
			if config.DefaultPath != "" {
				cores = append(cores, newHookedCore(config.DefaultPath, l))
			}
		}
	}

	tee := zapcore.NewTee(cores...)

	caller := zap.AddCaller()
	development := zap.Development()
	logger = zap.New(tee, caller, development)

	zap.ReplaceGlobals(logger)

}

// NewSugar 参考了https://gist.github.com/rnyrnyrny/282fe705d6e8dc012e482582d7c8ec0b
func NewSugar(name string) *zap.SugaredLogger {
	return logger.Named(name).Sugar()
}
