package zap

import (
	"fmt"
	"os"
	"strconv"

	kratoslog "github.com/go-kratos/kratos/v2/log"
	zaplog "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ kratoslog.Logger = (*ZapLogger)(nil)

// ZapLogger is a logger impl.
type ZapLogger struct {
	log  *zaplog.Logger
	Sync func() error
}

type LumberjackConfig struct {
	Enable     bool
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

func newLumberjackConfig() LumberjackConfig {
	envMaxSize := os.Getenv("LUMBERJACK_MAX_SIZE")
	if envMaxSize == "" {
		envMaxSize = "0"
	}
	maxSize, _ := strconv.Atoi(envMaxSize)

	envMaxAge := os.Getenv("LUMBERJACK_MAX_SIZE")
	if envMaxAge == "" {
		envMaxAge = "0"
	}
	maxAge, _ := strconv.Atoi(envMaxAge)

	envMaxBackups := os.Getenv("LUMBERJACK_MAX_SIZE")
	if envMaxBackups == "" {
		envMaxBackups = "0"
	}
	maxBackups, _ := strconv.Atoi(envMaxBackups)
	return LumberjackConfig{
		Enable:     os.Getenv("LUMBERJACK_ENABLE") == "true",
		Filename:   os.Getenv("LUMBERJACK_FILE_NAME"),
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		Compress:   true,
	}
}

// Log Implementation of logger interface.
func (l *ZapLogger) Log(level kratoslog.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}
	// Zap.Field is used when keyvals pairs appear
	var data []zaplog.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zaplog.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
	}
	switch level {
	case kratoslog.LevelDebug:
		l.log.Debug("", data...)
	case kratoslog.LevelInfo:
		l.log.Info("", data...)
	case kratoslog.LevelWarn:
		l.log.Warn("", data...)
	case kratoslog.LevelError:
		l.log.Error("", data...)
	}
	return nil
}

// NewZapLogger return a zap logger.
func NewZapLogger(encoder zapcore.EncoderConfig, level zaplog.AtomicLevel, opts ...zaplog.Option) *ZapLogger {
	ws := make([]zapcore.WriteSyncer, 0)
	ws = append(ws, zapcore.AddSync(os.Stdout))

	lumberjackConfig := newLumberjackConfig()
	if lumberjackConfig.Enable {
		ws = append(ws, zapcore.AddSync(&lumberjack.Logger{
			Filename:   lumberjackConfig.Filename,
			MaxSize:    lumberjackConfig.MaxSize,
			MaxAge:     lumberjackConfig.MaxAge,
			MaxBackups: lumberjackConfig.MaxBackups,
			Compress:   lumberjackConfig.Compress,
		}))
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.NewMultiWriteSyncer(ws...), level)
	zapLogger := zaplog.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}

// MustZapLogger return a zap logger
func MustZapLogger() *ZapLogger {
	encoder := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999999Z07"),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return NewZapLogger(
		encoder,
		zaplog.NewAtomicLevelAt(zapcore.DebugLevel),
		zaplog.AddStacktrace(
			zaplog.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zaplog.AddCaller(),
		zaplog.AddCallerSkip(2),
		zaplog.Development(),
	)
}
