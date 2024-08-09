package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func Init(logLevel string) error {
	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		return err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	globalLogger = zap.Must(cfg.Build())
	return nil
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}
