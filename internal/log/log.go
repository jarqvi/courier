package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func NewZapLogger() error {
	TZ := os.Getenv("TZ")
	if TZ == "" {
		TZ = "UTC"
	}

	location, err := time.LoadLocation(TZ)
	if err != nil {
		return fmt.Errorf("error in loading timezone: %w", err)
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
		return err
	}

	Logger = logger.Sugar()

	logger.Info("logger initialized")

	return nil
}
