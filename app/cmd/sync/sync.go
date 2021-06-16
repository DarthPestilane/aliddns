package sync

import (
	"github.com/urfave/cli"
)

func Command() cli.Command {
	// aliddns sync abc.com --ip=x.x.x.x --access-key=xxx --access-secret=xxx --region=xx
	return cli.Command{
		Name:      "sync",
		Usage:     "resolve an IP to the domain",
		UsageText: "aliddns sync abc.com --ip='1.2.3.4' -rr='www'",
		Aliases:   []string{"s"},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "rr",
				Usage: "The subdomain name such as 'www'",
				Value: "@",
			},
			cli.StringFlag{
				Name:  "ip",
				Usage: "Specify the ip to bind. current IP will be used if empty",
			},
			cli.StringFlag{
				Name:   "access-key",
				Usage:  "aliyun access key",
				EnvVar: "ACCESS_KEY",
			},
			cli.StringFlag{
				Name:   "access-secret",
				Usage:  "aliyun access key secret",
				EnvVar: "ACCESS_KEY_SECRET",
			},
			cli.StringFlag{
				Name:   "region",
				Usage:  "aliyun region",
				EnvVar: "REGION",
				Value:  "cn-hangzhou",
			},
		},
		Action: handle,
	}
}
