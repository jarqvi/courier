package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger() (*zap.SugaredLogger, error) {
	TZ := os.Getenv("TZ")
	if TZ == "" {
		TZ = "UTC"
	}

	location, err := time.LoadLocation(TZ)
	if err != nil {
		return nil, fmt.Errorf("error in loading timezone: %w", err)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "Time",
		LevelKey:       "Level",
		MessageKey:     "Message",
		CallerKey:      "File",
		EncodeTime:     customTimeEncoder(location),
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	config := zap.Config{
		Encoding:          "console",
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:       []string{"stdout"},
		EncoderConfig:     encoderConfig,
		DisableCaller:     false,
		DisableStacktrace: true,
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func customTimeEncoder(location *time.Location) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(location).Format("2006-01-02 15:04:05"))
	}
}
