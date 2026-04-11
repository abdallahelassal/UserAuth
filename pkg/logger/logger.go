package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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