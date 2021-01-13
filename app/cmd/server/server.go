package server

import (
	"fmt"
	"github.com/DarthPestilane/aliddns/app/helper"
	"github.com/urfave/cli"
	"strconv"
)

func Command() cli.Command {
	defaultPort, err := strconv.Atoi(helper.Env("PORT", "8888"))
	if err != nil {
		panic(fmt.Errorf("parse env PORT failed: %v", err))
	}
	cmd := cli.Command{
		Name: "run",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port",
				Value: defaultPort,
				Usage: "指定`监听端口`",
			},
		},
		Action: handler,
	}
	return cmd
}
