package log

import (
	"time"

	"go.uber.org/zap/zapcore"
)

func customTimeEncoder(location *time.Location) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(location).Format("2006-01-02 15:04:05"))
	}
}
