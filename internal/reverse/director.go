package reverse

import (
	"net/http"
	"net/url"
	"strings"

	"goblin/internal/plugin"
	"goblin/pkg/utils"

	log "unknwon.dev/clog/v2"
)

func (reverse *Reverse) Director(host string) func(request *http.Request) {
	// goblin 到服务器发出的请求
	// remote 用户端请求的 url
	target := reverse.AllowSite[host]
	remote, err := url.Parse(target)
	if err != nil {
		log.Info("target parse fail: %s", err.Error())
	}
	return func(request *http.Request) {
		request.URL.Scheme = remote.Scheme
		targetQuery := remote.RawQuery
		request.URL.Scheme = remote.Scheme
		request.URL.Host = remote.Host
		request.Host = remote.Host //
		referer := request.Header.Get("Referer")
		if referer != "" {
			refURL, err := url.Parse(referer)
			if err != nil {
				request.Header.Set("Referer", target)
				log.Trace("target parse fail: %s", err.Error())
			} else {
				refURL.Host = remote.Host
				request.Header.Set("Referer", refURL.String())
				log.Trace("referer parse : %s", refURL.String())
			}
		}

		if targetQuery == "" || request.URL.RawQuery == "" {
			request.URL.RawQuery = targetQuery + request.URL.RawQuery
		} else {
			request.URL.RawQuery = targetQuery + "&" + request.URL.RawQuery
		}
		if _, ok := request.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
		}

		request.Header.Del("Accept-Encoding")
		request.Header.Del("Content-Encoding")

		// 删除默认代理头
		request.Header["X-Forwarded-For"] = nil
		request.Header.Del("X-Real-Ip")

		//  插件系统 rule 处理请求数据
		if rules, ok := plugin.Plugins[host]; ok {
			for _, rule := range rules.Rule {
				if strings.Contains(request.URL.Path, rule.URL) {
					for _, rp := range rule.Replace {
						// 判断请求方法是否在里面
						if utils.EleInArray(request.Method, rp.Request.Method) {
							//处理响应数据
							err = rp.Request.Request(request)
							if err != nil {
								log.Info(err.Error())
							}
						}
					}
				}
			}
		}
	}
}
