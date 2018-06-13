package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"strings"
)

func cmdRun() cli.Command {
	cmd := cli.Command{
		Name: "run",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port",
				Value: 8888,
				Usage: "指定`监听端口`",
			},
		},
		Action: func(ctx *cli.Context) {
			port := ctx.Int("port")
			Log.Info(fmt.Sprintf("listening at port %d", port))
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				// query strings
				query := r.URL.Query()

				// domain name
				var domainName string
				if domains, has := query["domain_name"]; !has || domains[0] == "" {
					w.WriteHeader(422)
					b, _ := json.Marshal(map[string]interface{}{
						"success": false,
						"errors":  "domain_name is required",
					})
					w.Write(b)
					return
				} else {
					domainName = domains[0]
				}

				// rr
				var rr string
				if rrs, has := query["rr"]; !has || rrs[0] == "" {
					rr = "@"
				} else {
					rr = rrs[0]
				}

				addr := strings.TrimSpace(r.RemoteAddr)
				idx := strings.Index(addr, ":")
				var currentIP string
				if idx != -1 {
					currentIP = addr[:idx]
				} else {
					currentIP = addr
				}

				// bind
				dns := NewDns(domainName, currentIP, rr)
				Log.Info("=====")
				if err := dns.Bind(); err != nil {
					b, _ := json.Marshal(map[string]interface{}{
						"success": false,
						"errors":  err.Error(),
					})
					w.WriteHeader(400)
					w.Write(b)
					return
				}

				w.WriteHeader(200)
				b, err := json.Marshal(map[string]interface{}{
					"success": true,
					"message": fmt.Sprintf("set ip of '%s.%s' to %s", rr, domainName, currentIP),
				})
				if err != nil {
					Log.Error("decode response failed", err)
					return
				}
				w.Write(b)
				return
			})
			http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		},
	}
	return cmd
}
