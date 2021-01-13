package app

import (
	"github.com/DarthPestilane/aliddns/app/provider/logger"
	"github.com/sirupsen/logrus"
)

func Log() *logrus.Logger {
	return logger.Provide()
}
