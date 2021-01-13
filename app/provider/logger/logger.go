package logger

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

var appLog *logrus.Logger

func Register() {
	rawLogger := logrus.New()
	rawLogger.SetOutput(ioutil.Discard)
	rawLogger.SetLevel(logrus.TraceLevel)
	rawLogger.AddHook(NewStdoutHook(logrus.InfoLevel, &TextFormatter{ColorPrint: true}))
	appLog = rawLogger
}

func Provide() *logrus.Logger {
	return appLog
}
