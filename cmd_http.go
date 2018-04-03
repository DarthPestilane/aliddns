package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"strings"
)

func httpCmd() cli.Command {
	cmd := cli.Command{
		Name: "http",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port",
				Value: 8888,
				Usage: "指定`监听端口`",
			},
		},
		Action: func(ctx *cli.Context) {
			port := ctx.Int("port")
			log.Printf("listening at port %d\n", port)
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				header := w.Header()
				header.Set("Content-Type", "application/json")

				// query strings
				query := r.URL.Query()

				// domain name
				if domains, has := query["domain_name"]; !has || domains[0] == "" {
					w.WriteHeader(422)
					b, err := json.Marshal(map[string]interface{}{
						"success": false,
						"errors":  "domain_name is required",
					})
					if err != nil {
						panic(err)
					}
					w.Write(b)
					return
				} else {
					domainName = domains[0]
				}

				// rr
				if rrs, has := query["rr"]; !has || rrs[0] == "" {
					rr = "@"
				} else {
					rr = rrs[0]
				}

				// bind
				addr := strings.TrimSpace(r.RemoteAddr)
				idx := strings.Index(addr, ":")
				if idx != -1 {
					currentIP = addr[:idx]
				} else {
					currentIP = addr
				}
				bind()

				w.WriteHeader(200)
				b, err := json.Marshal(map[string]interface{}{
					"success": true,
					"message": fmt.Sprintf("set ip of '%s.%s' to %s", rr, domainName, currentIP),
				})
				if err != nil {
					panic(err)
				}
				w.Write(b)
				return
			})
			http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		},
	}
	return cmd
}
