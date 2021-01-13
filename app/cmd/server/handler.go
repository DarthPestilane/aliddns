package server

import (
	"fmt"
	"github.com/DarthPestilane/aliddns/app"
	"github.com/DarthPestilane/aliddns/app/dns"
	"github.com/DarthPestilane/aliddns/app/helper"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli"
	"net/http"
)

func handler(ctx *cli.Context) {
	port := ctx.Int("port")
	app.Log().Infof("listening at port %d", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// query strings
		query := r.URL.Query()

		// domain name
		domainName := query.Get("domain_name")
		if domainName == "" {
			responseFail(w, 422, "domain_name is required")
			return
		}

		// rr
		rr := query.Get("rr")
		if rr == "" {
			rr = "@"
		}

		// ip
		currentIP := helper.IP(r)

		// bind dns
		dnsHandler := dns.New(domainName, currentIP, rr)
		app.Log().Info("=====")
		if err := dnsHandler.Bind(); err != nil {
			responseFail(w, 400, err.Error())
			return
		}
		responseOK(w, 200, fmt.Sprintf("set ip of '%s.%s' to %s", rr, domainName, currentIP))
	})
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(fmt.Errorf("start http server failed: %s", err))
	}
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func responseFail(w http.ResponseWriter, code int, msg string) {
	setHeader(w, code)
	b, _ := jsoniter.Marshal(Response{
		Success: false,
		Error:   msg,
	})
	_, _ = w.Write(b)
}

func responseOK(w http.ResponseWriter, code int, msg string) {
	setHeader(w, code)
	b, _ := jsoniter.Marshal(Response{
		Success: true,
		Message: msg,
	})
	_, _ = w.Write(b)
}

func setHeader(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}
