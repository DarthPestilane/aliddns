package dns

import (
	"fmt"
	"github.com/DarthPestilane/aliddns/app"
	"github.com/DarthPestilane/aliddns/app/helper"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"time"
)

type Handler struct {
	client *alidns.Client
	ip     string
	domain string
	rr     string
}

func New(domain, ip, rr string) *Handler {
	if domain == "" || ip == "" || rr == "" {
		panic(fmt.Errorf("domain, ip or rr cannot be empty"))
	}
	client, err := alidns.NewClientWithAccessKey(helper.Env("REGION", "cn-hangzhou"), helper.Env("ACCESS_KEY"), helper.Env("ACCESS_KEY_SECRET"))
	if err != nil {
		panic(fmt.Errorf("new alidns client failed: %v", err))
	}
	return &Handler{
		client: client,
		ip:     ip,
		domain: domain,
		rr:     rr,
	}
}

func (dns *Handler) findRecords() (*alidns.DescribeDomainRecordsResponse, error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.DomainName = dns.domain
	resp, err := dns.client.DescribeDomainRecords(request)
	if err != nil {
		// try to fix timeout issue
		if clientErr, ok := err.(*errors.ClientError); ok && clientErr.ErrorCode() == errors.TimeoutErrorCode {
			// retry
			app.Log().WithError(err).Error("timeout. retry...")
			time.Sleep(time.Second)
			return dns.findRecords()
		}
		app.Log().WithError(err).Error("finding records failed")
		return nil, fmt.Errorf("finding records failed: %v", err)
	}
	return resp, nil
}

func (dns *Handler) addRecord() (*alidns.AddDomainRecordResponse, error) {
	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = dns.domain
	request.Type = "A"
	request.RR = dns.rr
	request.Value = dns.ip
	resp, err := dns.client.AddDomainRecord(request)
	if err != nil {
		app.Log().WithError(err).Error("adding record failed")
		return nil, fmt.Errorf("adding record failed: %v", err)
	}
	app.Log().Infof(`set ip of '%s.%s' to %s`, dns.rr, dns.domain, dns.ip)
	return resp, nil
}

func (dns *Handler) updateRecord(recordId string) (*alidns.UpdateDomainRecordResponse, error) {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.RecordId = recordId
	request.Type = "A"
	request.RR = dns.rr
	request.Value = dns.ip
	resp, err := dns.client.UpdateDomainRecord(request)
	if err != nil {
		app.Log().WithError(err).Error("updating record failed")
		return nil, fmt.Errorf("updating record failed: %v", err)
	}
	app.Log().Infof(`set ip of '%s.%s' to %s`, dns.rr, dns.domain, dns.ip)
	return resp, nil
}

func (dns *Handler) Bind() error {
	app.Log().Infof("current ip is %s", dns.ip)
	recordResp, err := dns.findRecords()
	if err != nil {
		return err
	}
	records := recordResp.DomainRecords.Record
	shouldAdd := true
	var recordId, recordValue string
	for _, r := range records {
		if r.RR == dns.rr {
			// 如果找到RR和输入里的rr相同的记录，则更新这条记录的解析。反之则添加一条新解析
			shouldAdd = false
			recordId = r.RecordId
			recordValue = r.Value
			break
		}
	}
	// add
	if shouldAdd {
		app.Log().Info("add domain record")
		if _, err := dns.addRecord(); err != nil {
			return err
		}
		return nil
	}
	// update record
	app.Log().Infof("domain ip is %s", recordValue)
	if recordValue == dns.ip {
		// no need updating
		app.Log().Info("ip not changed, no need updating")
		return nil
	}
	app.Log().Info("ip changed, update domain record")
	if _, err := dns.updateRecord(recordId); err != nil {
		return err
	}
	return nil
}
