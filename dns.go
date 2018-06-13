package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"time"
)

type Dns struct {
	client *alidns.Client
	IP     string
	Domain string
	RR     string
}

func NewDns(domain, ip, rr string) *Dns {
	if domain == "" || ip == "" || rr == "" {
		panic(fmt.Errorf("domain ip or rr cannot be empty"))
	}
	client, err := alidns.NewClientWithAccessKey(env("REGION", "cn-hangzhou"), env("ACCESS_KEY"), env("ACCESS_KEY_SECRET"))
	if err != nil {
		panic(fmt.Errorf("new alidns client failed: %v", err))
	}
	return &Dns{
		client: client,
		IP:     ip,
		Domain: domain,
		RR:     rr,
	}
}

func (dns *Dns) FindRecords() (*alidns.DescribeDomainRecordsResponse, error) {
	reqest := alidns.CreateDescribeDomainRecordsRequest()
	reqest.DomainName = dns.Domain
	resp, err := dns.client.DescribeDomainRecords(reqest)
	if err != nil {
		if clientErr, ok := err.(*errors.ClientError); ok && clientErr.ErrorCode() == errors.TimeoutErrorCode {
			// retry
			Log.Error("timeout. retry...", clientErr)
			time.Sleep(time.Second)
			return dns.FindRecords()
		}
		Log.Error("finding records failed", err)
		return nil, fmt.Errorf("finding records failed: %v", err)
	}
	return resp, nil
}

func (dns *Dns) AddRecord() (*alidns.AddDomainRecordResponse, error) {
	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = dns.Domain
	request.Type = "A"
	request.RR = dns.RR
	request.Value = dns.IP
	resp, err := dns.client.AddDomainRecord(request)
	if err != nil {
		Log.Error("adding record failed", err)
		return nil, fmt.Errorf("adding record failed: %v", err)
	}
	Log.Info(fmt.Sprintf(`set ip of '%s.%s' to %s`, dns.RR, dns.Domain, dns.IP))
	return resp, nil
}

func (dns *Dns) UpdateRecord(recordId string) (*alidns.UpdateDomainRecordResponse, error) {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.RecordId = recordId
	request.Type = "A"
	request.RR = dns.RR
	request.Value = dns.IP
	resp, err := dns.client.UpdateDomainRecord(request)
	if err != nil {
		Log.Error("updating record failed", err)
		return nil, fmt.Errorf("updating record failed: %v", err)
	}
	Log.Info(fmt.Sprintf(`set ip of '%s.%s' to %s`, dns.RR, dns.Domain, dns.IP))
	return resp, nil
}

func (dns *Dns) Bind() error {
	Log.Info(fmt.Sprintf("current ip is %s", dns.IP))
	recordResp, err := dns.FindRecords()
	if err != nil {
		return err
	}
	records := recordResp.DomainRecords.Record
	shouldAdd := true
	var recordId, recordValue string
	for _, r := range records {
		if r.RR == dns.RR {
			// 如果找到RR和输入里的rr相同的记录，则更新这条记录的解析。反之则添加一条新解析
			shouldAdd = false
			recordId = r.RecordId
			recordValue = r.Value
			break
		}
	}
	// add
	if shouldAdd {
		Log.Info("add domain record")
		if _, err := dns.AddRecord(); err != nil {
			return err
		}
		return nil
	}
	// update record
	Log.Info(fmt.Sprintf("domain ip is %s", recordValue))
	if recordValue == dns.IP {
		// no need updating
		Log.Info("ip not changed, no need updating")
		return nil
	}
	Log.Info("ip changed, update domain record")
	if _, err := dns.UpdateRecord(recordId); err != nil {
		return err
	}
	return nil
}
