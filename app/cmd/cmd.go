package cmd

import (
	"github.com/DarthPestilane/aliddns/app/cmd/server"
	"github.com/urfave/cli"
	"os"
)

func Run() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		server.Command(),
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
