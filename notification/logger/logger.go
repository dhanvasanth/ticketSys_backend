package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Log = l.Sugar()
}
