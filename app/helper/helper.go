package helper

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Env(key string, alternate ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		if len(alternate) == 0 {
			return ""
		}
		return alternate[0]
	}
	return v
}

func IP(req *http.Request) string {
	var ip string
	ip = strings.TrimSpace(req.Header.Get("x-forwarded-for"))
	if ip != "" {
		return ip
	}
	addr := strings.TrimSpace(req.RemoteAddr)
	idx := strings.Index(addr, ":")
	if idx != -1 {
		ip = addr[:idx]
	} else {
		ip = addr
	}
	return ip
}

func GeoIP() (string, error) {
	resp, err := http.DefaultClient.Get("https://api.ip.sb/geoip/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // nolint

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := gjson.GetBytes(data, "ip").String()
	if ip == "" {
		return "", fmt.Errorf("cannot find current IP")
	}
	return ip, nil
}
