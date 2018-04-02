package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net/http"
)

func httpCmd() cli.Command {
	cmd := cli.Command{
		Name: "http",
		Action: func(ctx *cli.Context) {
			log.Println("listening at port 8888")
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				header := w.Header()
				header.Set("Content-Type", "application/json")

				log.Println(r.URL.Query())
				query := r.URL.Query()

				// domain name
				var has bool
				_, has = query["domain_name"]
				if !has {
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
				}
				domainName = query["domain_name"][0]
				if domainName == "" {
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
				}

				// rr
				rr = query["rr"][0]
				if rr == "" {
					rr = "@"
				}

				// handle
				currentIP := r.RemoteAddr
				log.Printf("current ip is \t %s", currentIP)
				recordResp := findRecords()
				records := recordResp.DomainRecords.Record
				shouldAdd := true
				var recordId, recordValue string
				for _, r := range records {
					if r.RR == rr {
						// 如果找到RR和env里的rr相同的记录，则更新这条记录的解析。反之则添加一条新解析
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
			http.ListenAndServe(":8888", nil)
		},
	}
	return cmd
}
