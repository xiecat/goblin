package reverse

import (
	"fmt"
	"net/http"
	"net/http/httputil"

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
	return string(req)
}
