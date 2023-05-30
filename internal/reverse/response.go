package reverse

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"goblin/internal/plugin"
	"goblin/pkg/cache"
	"goblin/pkg/ipinfo"
	"goblin/pkg/logging"
	"goblin/pkg/notice"
	"goblin/pkg/utils"

	"github.com/sirupsen/logrus"
	log "unknwon.dev/clog/v2"
)

func (reverse *Reverse) ModifyResponse(shost string) func(response *http.Response) error {
	//shost 代理的主机或者域名
	return func(response *http.Response) error {
		// Stop CSPs and anti-XSS headers from ruining our fun
		response.Header.Del("Content-Security-Policy")
		response.Header.Del("X-XSS-Protection")
		//https://stackoverflow.com/questions/27358966/how-to-set-x-frame-options-on-iframe
		response.Header.Del("X-Frame-Options")

		response.Header.Del("Content-Security-Policy-Report-Only")

		// 删除缓存策略
		response.Header.Del("Expires")
		response.Header.Del("Last-Modified")
		response.Header.Del("Date")

		if response.Header.Get("Access-Control-Allow-Origin") != "" {
			//https://stackoverflow.com/questions/1653308/access-control-allow-origin-multiple-origin-domains
			if response.Request.Header.Get("Origin") != "" {
				response.Header.Set("Access-Control-Allow-Origin", response.Request.Header.Get("Origin"))
			} else {
				response.Header.Set("Access-Control-Allow-Origin", "*")
			}

		}
		//response.Header.Add("GoblinServer", Version)
		err := reverse.modifyLocationHeader(shost, response)
		if err != nil {
			log.Info(err.Error())
			return err
		}
		err = reverse.modifyCookieHeader(shost, response)
		if err != nil {
			log.Info(err.Error())
			return err
		}
		// todo 大文件替换问题方法需要合并
		// 文件大小处理
		cleng, err := strconv.Atoi(response.Header.Get("Content-Length"))
		if err != nil {
			cleng = -1
			log.Info("response  %s   maxContentLen:%d", response.Request.RequestURI, reverse.MaxContentLength)
		}
		//插件系统 rule 处理响应数据
		if rules, ok := plugin.Plugins[shost]; ok {
			for _, rule := range rules.Rule {
				urlmatch := false
				// url 匹配规则
				switch strings.ToLower(rule.Match) {
				case "word":
					urlmatch = response.Request.URL.Path == rule.URL
				case "prefix":
					urlmatch = strings.HasPrefix(response.Request.URL.Path, rule.URL)
				case "suffix":
					urlmatch = strings.HasSuffix(response.Request.URL.Path, rule.URL)
				}

				if urlmatch {
					log.Info("[plugin:%s] hit url: %s", rules.Name, response.Request.RequestURI)
					// replace
					if rule.Replace != nil {
						for _, rp := range rule.Replace {
							// 判断请求方法是否在里面
							if utils.EleInArray(response.Request.Method, rp.Request.Method) {
								log.Info("[plugin:%s.Replace.%s] Method match:%s", rules.Name, rule.URL, rp.Request.Method)
								//处理响应数据
								err = rp.Response.Response(reverse.MaxContentLength, response)
								if err != nil {
									log.Info("[plugin: %s.Replace] err:%s", rules.Name, err.Error())
								}
							}
						}
					}

					// 插件系统 dump 和提示
					if rule.Dump != nil {
						for _, dp := range rule.Dump {
							// 判断请求方法是否在里面
							if utils.EleInArray(response.Request.Method, dp.Request.Method) {
								//处理响应数据
								dete, msg := dp.Determine(reverse.MaxContentLength, response)
								start := time.Now()
								if dete {
									uid := strings.Join(response.Request.Header["X-Request-ID"], "")
									dplog, isCache := cache.DumpCache.Get(uid)
									cache.DumpCache.Delete(uid)
									if !isCache {
										log.Warn("[Plugin:%s.%s]not cache : %s\n", rules.Name, rule.URL, dplog)
										dplog = dumpJson(response.Request)
									}

									logging.AccLogger.WithFields(logrus.Fields{
										"method":     response.Request.Method,
										"url":        response.Request.URL.RequestURI(),
										"request_ip": GetClientIP(response.Request),
										"cost":       time.Since(start),
										"user-agent": response.Request.UserAgent(),
										"type":       "realReq",
									}).Warn(dplog)
									log.Warn("[Plugin:%s.%s]record: %s\n", rules.Name, rule.URL, dplog)
								}
								if msg {
									//todo 抑制消息, 异步告警
									target := reverse.AllowSite[shost]
									realHost := GetClientIP(response.Request)
									addr, _, err := utils.SplitHost(realHost)
									ipLocation := ipinfo.DB.Area(addr)
									log.Info("[dintalk]:目标: %s, 来源地址: %s, 地理位置: %s, 请求链接:%s", target, realHost, ipLocation, response.Request.RequestURI)
									go func(target, remote, location, ruleName, path string) {
										err = reverse.DingTalk.Send(&notice.Msg{
											Target:   target,
											Remote:   remote,
											Location: location,
											Rule:     ruleName,
											Path:     path,
										})
										if err != nil {
											log.Info(err.Error())
										}
									}(target, realHost, ipLocation, rules.Name, response.Request.RequestURI)
								}
							}
						}
					}

					//插件系统植入js
					if rule.InjectJs != nil {
						err = rule.InjectJs.ReplaceJs(response)
						if err != nil {
							log.Info("[plugin: %s] url: %s Error: %s", rules.Name, rule.URL, err.Error())
							return err
						}
					}

				}
			}
		}
		//检查 Content-Length 如果大于 MaxContentLength 直接返回不处理
		if cleng == -1 && reverse.MaxContentLength > 0 {
			log.Info("[response] Content-Length no set, set max: %d,will ignore and return", reverse.MaxContentLength)
			return nil
		}
		if cleng > reverse.MaxContentLength {
			log.Info("[response] Content-Length is :%s  set max: %d,will ignore and return", response.Header.Get("Content-Length"), reverse.MaxContentLength)
			return nil
		}
		//缓存处理放到最后
		if cache.GlobalCache.Type == "none" {
			return nil
		}
		return encodeCache(response)
	}
}

func (reverse *Reverse) modifyLocationHeader(shost string, response *http.Response) error {
	location, err := response.Location()
	if err != nil {
		// 没有及时返回
		if err == http.ErrNoLocation {
			return nil
		}
		return err
	}
	log.Trace("Location: %s", location.String())
	target := reverse.AllowSite[shost]
	targetHost, _ := url.Parse(target)
	locationHost, err := url.Parse(location.String())
	if err != nil {
		return err
	}
	log.Trace("targetHost.Host:%s, locationHost.Host: %s", targetHost.Host, locationHost.Host)

	if targetHost.Host == locationHost.Host {
		location.Scheme = ""
		location.Host = ""
	} else {
		log.Trace("url: %s,Location: %s", response.Request.URL, location.String())
	}
	if location.String() == "" {
		log.Trace("url: %s, Location is empty", response.Request.URL)
		return nil
	}
	response.Header.Set("Location", location.String())
	return nil
}

func (reverse *Reverse) modifyCookieHeader(shost string, response *http.Response) error {
	rcookies := response.Cookies()
	// 没有及时返回
	if len(rcookies) == 0 {
		return nil
	}
	var mcook []string
	addr, _, err := utils.SplitHost(shost)
	if err != nil {
		addr = shost
	}
	for _, value := range rcookies {
		// 关闭 Secure
		value.Secure = false
		// 关闭 httponly
		value.HttpOnly = true
		//代理 host
		value.Domain = addr
		mcook = append(mcook, value.String())
	}

	response.Header.Del("Set-Cookie")
	for _, mc := range mcook {
		response.Header.Add("Set-Cookie", mc)
	}
	return nil
}
