package app

import "github.com/DarthPestilane/aliddns/app/logger"

func Log() *logger.Logger {
	return logger.Provide()
}
