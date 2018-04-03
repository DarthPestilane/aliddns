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
				Name: "port",
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

				// handle
				addr := strings.TrimSpace(r.RemoteAddr)
				idx := strings.Index(addr, ":")
				if idx != -1 {
					currentIP = addr[:idx]
				} else {
					currentIP = addr
				}
				log.Printf("current ip is \t %s", currentIP)
				recordResp := findRecords()
				records := recordResp.DomainRecords.Record
				shouldAdd := true
				var recordId, recordValue string
				for _, r := range records {
					if r.RR == rr {
						// 如果找到RR和输入的rr相同的记录，则更新这条记录的解析。反之则添加一条新解析
						shouldAdd = false
						recordId = r.RecordId
						recordValue = r.Value
						break
					}
				}
				if shouldAdd {
					log.Printf("add domain record")
					addRecord()
				} else {
					// update record
					log.Printf("domain ip is \t %s", recordValue)
					if recordValue != currentIP {
						log.Println("ip changed, update domain record")
						updateRecord(recordId)
					} else {
						// no need updating
						log.Println("ip not changed, no need updating")
					}
				}

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
