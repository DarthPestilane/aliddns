package main

import (
	"github.com/DarthPestilane/aliddns/app/cmd"
	"github.com/DarthPestilane/aliddns/bootstrap"
	_ "github.com/joho/godotenv/autoload" // autoload .env
)

// ldflags -X
var (
	buildTime string // eg: 2020-04-24T15:56:31+0800
	gitCommit string // eg: d7b9655
	gitTag    string // eg: v1.1.1
)

func main() {
	bootstrap.Boot()
	cmd.Run(buildTime, gitCommit, gitTag)
}
