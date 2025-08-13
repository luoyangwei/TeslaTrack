package zap

import (
	"fmt"
	"os"

	kratoslog "github.com/go-kratos/kratos/v2/log"
	zaplog "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
)

var _ kratoslog.Logger = (*ZapLogger)(nil)

// ZapLogger is a logger impl.
type ZapLogger struct {
	log  *zaplog.Logger
	Sync func() error
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
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
		), level)
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
