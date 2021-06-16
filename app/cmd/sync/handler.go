package sync

import (
	"fmt"
	"github.com/DarthPestilane/aliddns/app/dns"
	"github.com/DarthPestilane/aliddns/app/helper"
	"github.com/urfave/cli"
	"strings"
)

func handle(ctx *cli.Context) error {
	// check domain
	if ctx.NArg() != 1 {
		return fmt.Errorf("there can be only one argument, and it should be the domain")
	}
	domain := strings.TrimSpace(ctx.Args()[0])
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	// check aliyun configs
	accessKey := strings.TrimSpace(ctx.String("access-key"))
	if accessKey == "" {
		return fmt.Errorf("access-key cannot be empty")
	}
	accessSecret := strings.TrimSpace(ctx.String("access-secret"))
	if accessSecret == "" {
		return fmt.Errorf("access-secret cannot be empty")
	}

	// check ip
	ip := strings.TrimSpace(ctx.String("ip"))
	if ip == "" {
		var err error
		ip, err = helper.GeoIP()
		if err != nil {
			return err
		}
	}
	rr := strings.TrimSpace(ctx.String("rr"))
	if rr == "" {
		rr = "@"
	}

	// bind
	dnsHandler := dns.New(domain, ip, rr)
	return dnsHandler.Bind()
}
