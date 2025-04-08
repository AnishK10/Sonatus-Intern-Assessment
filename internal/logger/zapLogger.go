package logger

import (
	"go.uber.org/zap"
)

var Zap *zap.Logger

func InitZap() {
	logger, _ := zap.NewProduction()
	Zap = logger
}
