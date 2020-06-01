package helper

import (
	"net/http"
	"os"
	"strings"
)

func Env(key string, missing ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		if len(missing) == 0 {
			return ""
		}
		return missing[0]
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
