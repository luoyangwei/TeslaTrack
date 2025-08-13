package zap

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
)

func TestZapLogger(t *testing.T) {
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
	logger := NewZapLogger(
		encoder,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)
	zlog := log.NewHelper(logger)
	zlog.Infow("name", "kratos", "from", "opensource")
	zlog.Infow("name", "kratos", "from")

	// zap stdout/stderr Sync bugs in OSX, see https://github.com/uber-go/zap/issues/370
	_ = logger.Sync()
}
