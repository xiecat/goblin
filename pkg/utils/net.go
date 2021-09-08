package utils

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "unknwon.dev/clog/v2"
)

// IsCidr determines if the given ipctl is a cidr range
func IsCidr(ip string) bool {
	_, _, err := net.ParseCIDR(ip)
	return err == nil
}

// IsIP determines if the given string is a valid ipctl
func IsIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsPort determines if the given int is a valid port
func IsPort(port int) bool {
	return port > 0 && port <= 65535
}

func VisitURL(url1 string) bool {
	res, err := http.Get(url1)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	_, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return true
}

func IsURL(url1 string) bool {
	_, err := url.Parse(url1)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return true
}

// IsHost determines if the given int is a valid host (ipctl/domainmde)
func IsHost(host string) bool {
	_, err := net.ResolveIPAddr("ip", host)
	if err == nil || IsIP(host) {
		return true
	}
	return false
}

func SplitHost(host string) (addr string, port int, err error) {
	host = strings.TrimSpace(host)
	sphost := strings.Split(host, ":")
	if len(sphost) != 2 {
		return "", 0, errors.New("host fomat err: " + host)
	}
	port, err = strconv.Atoi(sphost[1])
	if err != nil {
		return "", 0, errors.New("port is error")
	}
	return sphost[0], port, nil
}

func ValidProxy(pURL, acURL string) bool {
	urli := url.URL{}
	urlproxy, _ := urli.Parse(pURL)
	c := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
		},
	}
	resp, err := c.Get(acURL)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Trace("%s\n", body)
	return true
}
