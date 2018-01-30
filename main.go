package main

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	client          *alidns.Client
	domainName      string
	rr              string = "@" // default is `@`
	currentIP       string
	intervalMinutes int
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	for _, s := range os.Args {
		arg := strings.Split(s, "=")
		key := strings.Trim(strings.ToLower(arg[0]), "-")
		var val string
		if len(arg) > 1 {
			val = strings.TrimSpace(arg[1])
		}
		switch key {
		case "domain-name":
			if val != "" {
				domainName = val
			}
		case "rr":
			if val != "" {
				rr = val
			}
		case "interval-min":
			if val != "" {
				min, err := strconv.Atoi(val)
				if err != nil {
					panic(err)
				}
				if min < 1 {
					panic(errors.New("interval minutes must be greater than 0"))
				}
				intervalMinutes = min
			}
		}
	}
	// check vars
	if domainName == "" {
		panic(errors.New("domain-name must be specified"))
	}
	client = newClient()
}

func main() {
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
}

func env(key string, missing ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		if len(missing) == 0 {
			return ""
		}
		return missing[0]
	}
	return v
}

func newClient(ch ...chan *alidns.Client) *alidns.Client {
	client, err := alidns.NewClientWithAccessKey(env("REGION", "cn-hangzhou"), env("ACCESS_KEY"), env("ACCESS_KEY_SECRET"))
	if err != nil {
		panic(err)
	}
	if len(ch) != 0 {
		ch[0] <- client
	}
	return client
}

func getCurrentIP() (string, error) {
	// response, err := http.Get("http://members.3322.org/dyndns/getip")
	response, err := http.Get("http://35.194.248.24:81") // ip getter
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(string(b))
	return ip, nil
}

func findRecords() *alidns.DescribeDomainRecordsResponse {
	reqest := alidns.CreateDescribeDomainRecordsRequest()
	reqest.DomainName = domainName
	resp, err := client.DescribeDomainRecords(reqest)
	if err != nil {
		panic(err)
	}
	return resp
}

func addRecord() *alidns.AddDomainRecordResponse {
	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = domainName
	request.Type = "A"
	request.RR = rr
	request.Value = currentIP
	resp, err := client.AddDomainRecord(request)
	if err != nil {
		panic(err)
	}
	log.Printf(`set ip of '%s.%s' to %s`, rr, domainName, currentIP)
	return resp
}

func updateRecord(recordId string) *alidns.UpdateDomainRecordResponse {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.RecordId = recordId
	request.Type = "A"
	request.RR = rr
	request.Value = currentIP
	resp, err := client.UpdateDomainRecord(request)
	if err != nil {
		panic(err)
	}
	log.Printf(`set ip of '%s.%s' to %s`, rr, domainName, currentIP)
	return resp
}
