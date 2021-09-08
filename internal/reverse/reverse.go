package reverse

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"goblin/pkg/cache"

	log "unknwon.dev/clog/v2"
)

func (reverse *Reverse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	if target, ok := reverse.AllowSite[host]; ok {
		remote, err := url.Parse(target)
		if err != nil {
			log.Info("target parse fail: %s", err.Error())
			return
		}
		r.URL.Scheme = remote.Scheme
		r.URL.Host = remote.Host
		r.Host = remote.Host
		//dump req log
		if logLevel == 1 {
			log.Trace("record:\n%s", dumpReq(r))
		}
		// dump
		if logLevel == 2 && r.Method == "POST" {
			log.Info("record:\n%s", dumpReq(r))
		}
		log.Info("[c->p] host: %s,RemoteAddr: %s,URI: %s", host, GetClientIP(r), r.RequestURI)
		// 处理缓存
		if cache.GlobalCache.Type != "none" {
			urlobj, err := decodeCache(r.URL.String())
			if err != nil {
				log.Trace("cache type:%s, no hit file: %s", cache.GlobalCache.Type, r.URL.String())
			} else {
				// 取缓存文件
				log.Trace("hit file: %s", r.URL.String())

				headers := urlobj.ResponseHeaders
				//处理 header
				for hkey, hvalue := range headers {
					w.Header().Set(hkey, strings.Join(hvalue, ";"))
				}
				// 处理状态码
				w.WriteHeader(urlobj.StatusCode)
				// 处理内容一定要最后处理。先处理 header 和状态码无法生效。被重写了
				w.Write(urlobj.Content)
				return
			}
		}
		// 直接从缓存取出
		if fn, ok := reverse.HostProxy[host]; ok {
			fn.ServeHTTP(w, r)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = reverse.Director(host)
		proxy.Transport = reverse.Transport()
		proxy.ModifyResponse = reverse.ModifyResponse(host) // nolint
		proxy.BufferPool = Bufferpool
		reverse.HostProxy[host] = proxy // 放入缓存
		proxy.ServeHTTP(w, r)
		return
	}
	w.Write([]byte("403: Host forbidden"))
}
