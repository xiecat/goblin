package reverse

import (
	"github.com/sirupsen/logrus"
	"goblin/internal/plugin"
	"goblin/pkg/logging"
	"goblin/pkg/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

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
		if logLevel == logrus.InfoLevel {
			start := time.Now()
			reqraw := dumpJson(r)
			logging.AccLogger.WithFields(logrus.Fields{
				"method":     r.Method,
				"url":        r.URL.RequestURI(),
				"request_ip": GetClientIP(r),
				"user-agent": r.UserAgent(),
				"cost":       time.Since(start),
				"type":       "clientReq",
			}).Info(reqraw)
		} else if logLevel == logrus.WarnLevel && r.Method == "POST" {
			start := time.Now()
			reqraw := dumpJson(r)
			logging.AccLogger.WithFields(logrus.Fields{
				"method":     r.Method,
				"url":        r.URL.RequestURI(),
				"request_ip": GetClientIP(r),
				"user-agent": r.UserAgent(),
				"cost":       time.Since(start),
				"type":       "clientReq",
			}).Warn(reqraw)
		}

		log.Info("[c->p] host: %s,RemoteAddr: %s,URI: %s", host, GetClientIP(r), r.RequestURI)
		//response
		//插件系统 rule 处理响应数据
		if rules, ok := plugin.Plugins[host]; ok {
			//dump

			for _, rule := range rules.Rule {
				for _, dp := range rule.Dump {
					if dp.NeedCache(r) {
						uuidstr := utils.GenerateUUID()
						r.Header["X-Request-ID"] = []string{uuidstr}
						cache.DumpCache.Set(uuidstr, dumpJson(r), 60*time.Second)
					}
				}

				urlmatch := false
				// url 匹配规则
				switch strings.ToLower(rule.Match) {
				case "word":
					urlmatch = r.URL.Path == rule.URL
				case "prefix":
					urlmatch = strings.HasPrefix(r.URL.Path, rule.URL)
				case "suffix":
					urlmatch = strings.HasPrefix(r.URL.Path, rule.URL)
				}

				if urlmatch {
					log.Info("[plugin:%s] hit url: %s", rules.Name, r.RequestURI)
					// replace
					if rule.Replace != nil {
						for _, rp := range rule.Replace {
							// 判断请求方法是否在里面
							if utils.EleInArray(r.Method, rp.Request.Method) {
								log.Info("[plugin:%s.Replace.%s] Method match:%s", rules.Name, rule.URL, rp.Request.Method)
								//处理响应数据

								if rp.Response != nil {
									if rp.Response.Location != "" {
										log.Info("[plugin: %s.Location]: %s", rules.Name, rp.Response.Location)
										w.Header().Set("Location", rp.Response.Location)
										w.WriteHeader(302)
										return
									}
								}

							}
						}
					}
				}
			}
		}

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
	//w.Header().Set("GoblinServer", Version)
	w.Write([]byte("403: Host forbidden"))
}
