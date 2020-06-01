package server

import (
	"fmt"
	"github.com/DarthPestilane/aliddns/app"
	"github.com/DarthPestilane/aliddns/app/dns"
	"github.com/DarthPestilane/aliddns/app/helper"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli"
	"net/http"
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
		Action: func(ctx *cli.Context) {
			port := ctx.Int("port")
			app.Log().Info(fmt.Sprintf("listening at port %d", port))
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				// query strings
				query := r.URL.Query()

				// domain name
				var domainName string
				if domains, has := query["domain_name"]; !has || domains[0] == "" {
					w.WriteHeader(422)
					b, _ := jsoniter.Marshal(map[string]interface{}{
						"success": false,
						"errors":  "domain_name is required",
					})
					_, _ = w.Write(b)
					return
				} else {
					domainName = domains[0]
				}

				// rr
				var rr = "@"
				if rrs, has := query["rr"]; has && rrs[0] != "" {
					rr = rrs[0]
				}

				// ip
				currentIP := helper.IP(r)

				// bind dns
				dnsHandler := dns.New(domainName, currentIP, rr)
				app.Log().Info("=====")
				if err := dnsHandler.Bind(); err != nil {
					b, _ := jsoniter.Marshal(map[string]interface{}{
						"success": false,
						"errors":  err.Error(),
					})
					w.WriteHeader(400)
					_, _ = w.Write(b)
					return
				}

				w.WriteHeader(200)
				b, err := jsoniter.Marshal(map[string]interface{}{
					"success": true,
					"message": fmt.Sprintf("set ip of '%s.%s' to %s", rr, domainName, currentIP),
				})
				if err != nil {
					app.Log().Error("decode response failed", err)
					return
				}
				_, _ = w.Write(b)
			})
			if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
				panic(fmt.Errorf("start http server failed: %s", err))
			}
		},
	}
	return cmd
}
