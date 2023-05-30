package dump

import (
	"bytes"
	"goblin/pkg/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "unknwon.dev/clog/v2"
)

func (dump *Dump) Determine(maxContentLength int, response *http.Response) (dete, notice bool) {
	start := time.Now()
	defer log.Info("[time] url: %s, dump hand time: %v", response.Request.RequestURI, time.Since(start))
	if dump == nil {
		return false, false
	}
	cleng, err := strconv.Atoi(response.Header.Get("Content-Length"))
	if err != nil {
		cleng = -1
		log.Info("response Header Error: %s Content-Length:%s ", response.Request.RequestURI, response.Header.Get("Content-Length"))
	}
	r := response.Request
	// 如何没有任何请求方式支持直接返回
	if len(dump.Request.Method) == 0 {
		return false, dump.Notice
	}
	// 判断 method 是否符合 dump 规则
	if utils.EleInArray(r.Method, dump.Request.Method) {
		//为 nil 不匹配  Response
		if dump.Response == nil {
			return true, dump.Notice
		}
		// 如果状态码为 0 或者相等则成功
		if response.StatusCode == dump.Response.Status || dump.Response.Status == 0 {
			if dump.Response.Header != nil {
				for hkey, hvalue := range dump.Response.Header {
					exist := strings.Contains(response.Header.Get(hkey), hvalue)
					if !exist {
						return false, dump.Notice
					}
				}
			}
			// 如果 body 没有值默认匹配
			if dump.Response.Body == "" {
				return true, dump.Notice
			}
			// 检查 Content-Length 如果大于 MaxContentLength 直接返回不处理 小于0 放行
			if cleng == -1 && maxContentLength > 0 {
				log.Info("response Content-Length no set  set max: %d, rule Dump Body will ignore ", maxContentLength)
				return true, dump.Notice
			}
			if cleng > maxContentLength && maxContentLength > 0 {
				log.Info("response Content-Length is :%s  set max: %d, rule Dump Body will ignore ", response.Header.Get("Content-Length"), maxContentLength)
				return true, dump.Notice
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Info("%s", err.Error())
				return false, dump.Notice
			}
			if bytes.Contains(body, []byte(dump.Response.Body)) {
				return true, dump.Notice
			}
			response.Body = io.NopCloser(bytes.NewReader(body))
		}
	}

	return false, dump.Notice
}

func (dump *Dump) NeedCache(r *http.Request) bool {
	start := time.Now()
	defer log.Info("[time] url: %s, dump cache hand time: %v", r.RequestURI, time.Since(start))
	if dump == nil {
		return false
	}

	// 如何没有任何请求方式支持直接返回
	if len(dump.Request.Method) == 0 {
		return false
	}
	// 判断 method 是否符合 dump 规则
	if utils.EleInArray(r.Method, dump.Request.Method) {
		//为 nil 不匹配  Response
		return true

	}

	return false
}
