package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	if err := cfg.Level.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}

	cfg.Encoding = "json"
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build(
		zap.AddCaller(),
	)
	if err != nil {
		return nil, err
	}

	return logger, nil
}
