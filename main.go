package main

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	client     *alidns.Client
	domainName string
	rr         string
	currentIP  string
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	domainName = env("DOMAIN_NAME")
	rr = env("RR", "@")
	client = newClient()
}

func main() {
	for {
		currentIP = getCurrentIP()
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
		time.Sleep(5 * 60 * time.Second)
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

func getCurrentIP(ch ...chan string) string {
	response, err := http.Get("http://members.3322.org/dyndns/getip")
	if err != nil {
		panic(err.Error())
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	ip := strings.TrimSpace(string(b))
	if len(ch) != 0 {
		ch[0] <- ip
	}
	return ip
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
