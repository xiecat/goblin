package reverse

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"goblin/pkg/utils"
)

func GetClientIP(req *http.Request) string {
	RemoteIP := req.RemoteAddr
	if ProxyHeader != "Remote-Addr" {
		_, port, _ := utils.SplitHost(RemoteIP)
		rp := req.Header.Get(ProxyHeader)
		if rp != "" {
			RemoteIP = fmt.Sprintf("%s:%d", rp, port)
		}
	}
	return RemoteIP
}

func dumpReq(r *http.Request) string {
	info := "---------------------  %s start  ------------------------\n-------- requests from: %s\n----- requests raw:  -----\n%s\n--------------------  %s end  ---------------------------"
	req, _ := httputil.DumpRequest(r, true)
	return fmt.Sprintf(info, r.URL.RequestURI(), GetClientIP(r), string(req), r.URL.RequestURI())
}

func dumpJson(r *http.Request) string {
	req, _ := httputil.DumpRequest(r, true)
	if r.Method == "POST" {
		return string(req) + "\n\n"
	}
	return string(req)
}

// IsWebSocketRequest returns a boolean indicating whether the request has the
func IsWebSocketRequest(r *http.Request) bool {
	contains := func(key, val string) bool {
		vv := strings.Split(r.Header.Get(key), ",")
		for _, v := range vv {
			if val == strings.ToLower(strings.TrimSpace(v)) {
				return true
			}
		}
		return false
	}
	if !contains("Connection", "upgrade") {
		return false
	}
	if !contains("Upgrade", "websocket") {
		return false
	}
	return true
}
