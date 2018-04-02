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
			for {
				ip, err := getCurrentIP()
				if err == http.ErrHandlerTimeout {
					log.Println("request current ip timeout, try again now")
					continue
				} else if err != nil {
					panic(err)
				}
				currentIP = ip
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
				if intervalMinutes == 0 {
					return
				}
				time.Sleep(time.Duration(intervalMinutes) * time.Minute)
			}
		},
	}
	return cmd
}
