package reverse

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

func (reverse *Reverse) Transport() http.RoundTripper {
	//ProxyAddr, err := url.Parse(ProxyServerAddr)
	// 检查代理
	proxyURLFunc := func(purl *url.URL) func(*http.Request) (*url.URL, error) {
		if ProxyServerAddr == nil {
			return nil
		}
		return func(*http.Request) (*url.URL, error) {
			return purl, nil
		}
	}(ProxyServerAddr)

	return &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			//RootCAs:            rootCAs,
		},
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second, //连接超时
			KeepAlive: 30 * time.Second, //长连接超时时间
		}).DialContext,
		Proxy: proxyURLFunc,
		//DialContext: dialer.Dial,
		MaxIdleConns:          reverse.MaxIdleConns, //最大空闲连接
		MaxConnsPerHost:       reverse.MaxConnsPerHost,
		MaxIdleConnsPerHost:   reverse.MaxIdleConnsPerHost,
		IdleConnTimeout:       reverse.IdleConnTimeout,       //空闲超时时间
		TLSHandshakeTimeout:   reverse.TLSHandshakeTimeout,   //tls握手超时时间
		ExpectContinueTimeout: reverse.ExpectContinueTimeout, //100-continue 超时时间
	}
}
