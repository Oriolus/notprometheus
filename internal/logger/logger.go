package logger

import (
	"fmt"
	"go.uber.org/zap"
)

var (
	Logger = &zap.Logger{}
)

func init() {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func Infof(template string, args ...any) {
	Logger.Info(fmt.Sprintf(template, args...))
}

func Errorf(template string, args ...any) {
	Logger.Error(fmt.Sprintf(template, args...))
}

func Fatalf(template string, args ...any) {
	Logger.Fatal(fmt.Sprintf(template, args...))
}
