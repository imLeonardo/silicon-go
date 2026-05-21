package logger

import (
	"fmt"
	"silicon-go/config"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	registry.NopHook

	Config    *config.Config
	ZapLogger *zap.Logger
}

type LogField = zap.Field

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// 时间格式：2006-01-02 15:04:05
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger", // 可选，一般可隐藏
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,                   // 带颜色：DEBUG/INFO/WARN/ERROR
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), // 简洁时间
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 显示为 server/server.go:42
	}
}

func newRotatelog(config *config.Config) (*rotatelogs.RotateLogs, error) {
	// 替换 {%type} 为 cluster.type
	pattern := strings.Replace(config.Logger.FilePattern, "{$cluster.type}", config.Cluster.Type, -1)
	// 创建轮转器
	rl, err := rotatelogs.New(
		pattern,
		// 每小时轮转一次
		rotatelogs.WithRotationTime(time.Hour),
		// 可选：保留最近MaxAge天的日志
		rotatelogs.WithMaxAge(time.Duration(config.Logger.MaxAge*24)*time.Hour),
		// 可选：设置文件权限
		// rotatelogs.WithFileMode(0644),
	)
	if err != nil {
		return nil, err
	}
	return rl, nil
}

func newHourlyLogger(config *config.Config) (*zap.Logger, error) {
	if config.Logger.IsDocker {
		return zap.NewProduction()
	}

	// 创建控制台core
	consoleEncoder := zapcore.NewConsoleEncoder(newEncoderConfig())
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.Level(config.Logger.Level))

	// 创建文件core
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	rl, err := newRotatelog(config)
	if err != nil {
		return nil, err
	}
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(rl), zapcore.Level(config.Logger.Level))

	combinedCore := zapcore.NewTee(consoleCore, fileCore)
	logger := zap.New(
		combinedCore,
		zap.AddCaller(),
		zap.AddCallerSkip(config.Logger.CallerSkip), // 在默认基础上再跳过 1 层，以跳过 logger 包
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return logger, nil
}

func NewLogger(config *config.Config) (*Logger, error) {
	zapLogger, err := newHourlyLogger(config)
	if err != nil {
		return nil, err
	}

	return &Logger{Config: config, ZapLogger: zapLogger}, nil
}

func (l *Logger) OnStart() error {
	l.Info("logger initialized successfully")
	return nil
}

func (l *Logger) OnStop() error {
	err := l.ZapLogger.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (l *Logger) Debug(v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.DebugLevel) {
		return
	}

	msg := fmt.Sprintf(v[0].(string), v[1:]...)
	l.ZapLogger.Debug(msg)
}
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.DebugLevel) {
		return
	}

	msg := fmt.Sprintf(format, v...)
	l.ZapLogger.Debug(msg)
}
func (l *Logger) Debugln(v ...interface{}) {
	l.Debug(v...)
}
func (l *Logger) DebugWithFields(msg string, fields ...LogField) {
	if l.Config.Logger.Level > int32(zapcore.DebugLevel) {
		return
	}

	l.ZapLogger.Debug(msg, fields...)
}

func (l *Logger) Info(v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.InfoLevel) {
		return
	}

	msg := fmt.Sprintf(v[0].(string), v[1:]...)
	l.ZapLogger.Info(msg)
}
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.InfoLevel) {
		return
	}

	msg := fmt.Sprintf(format, v...)
	l.ZapLogger.Info(msg)
}
func (l *Logger) Infoln(v ...interface{}) {
	l.Info(v)
}
func (l *Logger) InfoWithFields(msg string, fields ...LogField) {
	if l.Config.Logger.Level > int32(zapcore.InfoLevel) {
		return
	}
	l.ZapLogger.Info(msg, fields...)
}

func (l *Logger) Warn(v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.WarnLevel) {
		return
	}

	msg := fmt.Sprintf(v[0].(string), v[1:]...)
	l.ZapLogger.Warn(msg)
}
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.WarnLevel) {
		return
	}

	msg := fmt.Sprintf(format, v...)
	l.ZapLogger.Warn(msg)
}
func (l *Logger) Warnln(v ...interface{}) {
	l.Warn(v)
}
func (l *Logger) WarnWithFields(msg string, fields ...LogField) {
	if l.Config.Logger.Level > int32(zapcore.WarnLevel) {
		return
	}

	l.ZapLogger.Warn(msg, fields...)
}

func (l *Logger) Error(v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.ErrorLevel) {
		return
	}

	msg := fmt.Sprintf(v[0].(string), v[1:]...)
	l.ZapLogger.Error(msg)
}
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.ErrorLevel) {
		return
	}

	msg := fmt.Sprintf(format, v...)
	l.ZapLogger.Error(msg)
}
func (l *Logger) Errorln(v ...interface{}) {
	l.Error(v)
}
func (l *Logger) ErrorWithFields(msg string, fields ...LogField) {
	if l.Config.Logger.Level > int32(zapcore.ErrorLevel) {
		return
	}

	l.ZapLogger.Error(msg, fields...)
}

func (l *Logger) Panic(v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.PanicLevel) {
		return
	}

	msg := fmt.Sprintf(v[0].(string), v[1:]...)
	l.ZapLogger.Panic(msg)
}
func (l *Logger) Panicf(format string, v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.PanicLevel) {
		return
	}

	msg := fmt.Sprintf(format, v...)
	l.ZapLogger.Panic(msg)
}
func (l *Logger) Panicln(v ...interface{}) {
	l.Panic(v)
}
func (l *Logger) PanicWithFields(msg string, fields ...LogField) {
	if l.Config.Logger.Level > int32(zapcore.PanicLevel) {
		return
	}

	l.ZapLogger.Panic(msg, fields...)
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.FatalLevel) {
		return
	}

	msg := fmt.Sprintf(v[0].(string), v[1:]...)
	l.ZapLogger.Fatal(msg)
}
func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.Config.Logger.Level > int32(zapcore.FatalLevel) {
		return
	}

	msg := fmt.Sprintf(format, v...)
	l.ZapLogger.Fatal(msg)
}
func (l *Logger) Fatalln(v ...interface{}) {
	l.Fatal(v)
}
func (l *Logger) FatalWithFields(msg string, fields ...LogField) {
	if l.Config.Logger.Level > int32(zapcore.FatalLevel) {
		return
	}

	l.ZapLogger.Fatal(msg, fields...)
}

func (l *Logger) GetInternalLogger() any {
	return l.ZapLogger
}

func (l *Logger) LogField(key string, value interface{}) LogField {
	return zap.Any(key, value)
}

func (l *Logger) WithFields(fields map[string]interface{}) interfaces.Logger {
	zapFields := make([]LogField, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	l.ZapLogger = l.ZapLogger.With(zapFields...)
	return l
}

func (l *Logger) WithField(key string, value interface{}) interfaces.Logger {
	l.ZapLogger = l.ZapLogger.With(LogField{Key: key, Type: zapcore.ReflectType, Interface: value})
	return l
}

func (l *Logger) WithError(err error) interfaces.Logger {
	l.ZapLogger = l.ZapLogger.With(zap.Error(err))
	return l
}
