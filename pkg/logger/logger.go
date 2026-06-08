package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log zap.Logger

func NewLogger(env string) *zap.Logger{
	var config zap.Config
	
	if env == "production"{
		config = zap.NewProductionConfig()
	}else{
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	logger , _ := config.Build()
	return logger
}
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Any(key string, value any) zap.Field {
	return zap.Any(key, value)
}

func String(key, value string) zap.Field {
	return zap.String(key, value)
}