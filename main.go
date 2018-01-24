package main

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	chClient := make(chan *alidns.Client)
	chCurrentIP := make(chan string)
	go newClient(chClient)
	go getCurrentIP(chCurrentIP)
	client = <-chClient
	currentIP = <-chCurrentIP

	// client = newClient()
	// currentIP = getCurrentIP()

	log.Printf("current ip is \t %s", currentIP)
}

func main() {
	recordResp := findRecords()
	records := recordResp.DomainRecords.Record
	if len(records) == 0 {
		// add record
		log.Println("add domain record")
		addRecord()
	} else {
		// update record
		recordId := records[0].RecordId
		domainIP := records[0].Value
		log.Printf("domain ip is \t %s", domainIP)
		if domainIP != currentIP {
			log.Println("ip changed, update domain record")
			updateRecord(recordId)
		} else {
			log.Println("ip not changed, cancel update")
		}
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
	response, err := http.Get("http://ipinfo.io/json")
	if err != nil {
		panic(err.Error())
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(b, &body); err != nil {
		panic(err)
	}
	ip := body["ip"].(string)
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
	log.Printf(`set ip to %s`, currentIP)
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
	log.Printf(`set ip to %s`, currentIP)
	return resp
}
