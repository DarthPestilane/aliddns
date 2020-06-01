package cmd

import (
	"github.com/DarthPestilane/aliddns/app/cmd/server"
	"github.com/DarthPestilane/aliddns/app/cmd/version"
	"github.com/urfave/cli"
	"os"
)



func Run(buildTime, gitCommit, gitTag string) {
	info := &version.BuildInfo{
		BuildTime: buildTime,
		GitCommit: gitCommit,
		GitTag:    gitTag,
	}
	if info.GitTag == "" {
		info.GitTag = "in-dev"
	}

	app := cli.NewApp()
	app.Version = info.GitTag
	app.Commands = []cli.Command{
		server.Command(),
		version.Command(info),
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
