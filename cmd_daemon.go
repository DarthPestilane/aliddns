package main

import (
	"errors"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"time"
)

func daemonCmd() cli.Command {
	cmd := cli.Command{
		Name: "start",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "domain-name",
				Value: "",
				Usage: "请告诉我`域名`",
			},
			cli.StringFlag{
				Name:  "rr",
				Value: "@",
				Usage: "请告诉我`RR`",
			},
			cli.IntFlag{
				Name:  "interval-min",
				Value: 0,
				Usage: "请告诉我刷新时间间隔`分钟`数，如果是0则不会自动刷新",
			},
		},
		Action: func(ctx *cli.Context) {
			// prepare vars
			domainName = ctx.String("domain-name")
			rr = ctx.String("rr")
			intervalMinutes = ctx.Int("interval-min")
			if domainName == "" {
				panic(errors.New("domain-name must be specified"))
			}
			if intervalMinutes < 0 {
				panic(errors.New("interval minutes must be greater than 0"))
			}
			var err error
			for {
				currentIP, err = getCurrentIP()
				if err == http.ErrHandlerTimeout {
					log.Println("request current ip timeout, try again now")
					return
				} else if err != nil {
					panic(err)
				}
				bind()
				if intervalMinutes == 0 {
					return
				}
				time.Sleep(time.Duration(intervalMinutes) * time.Minute)
			}
		},
	}
	return cmd
}
