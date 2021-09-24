package replace

import (
	"bytes"
	"fmt"
	"goblin/pkg/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "unknwon.dev/clog/v2"
)

// https://annevankesteren.nl/2005/02/javascript-mime-type

func (rpRule *Response) Response(maxContentLength int, response *http.Response) error {
	start := time.Now()
	defer log.Info("[time] url: %s, replace hand time: %v", response.Request.RequestURI, time.Since(start))
	if rpRule == nil {
		return nil
	}

	cleng, err := strconv.Atoi(response.Header.Get("Content-Length"))
	if err != nil {
		// 如果没有 "Content-Length" 头设置头为 -1
		cleng = -1
		log.Info("response Header Error: %s Content-Length: no set maxContentLen: %d", response.Request.RequestURI, maxContentLength)
	}
	if rpRule.Status > 0 {
		response.StatusCode = rpRule.Status
	}

	// cookies 检查
	if rpRule.Cookie != nil {
		rcookies := response.Cookies()
		// 没有及时返回
		if len(rcookies) != 0 {
			var mcook []string

			for _, value := range rcookies {
				// 关闭 Secure
				value.Secure = rpRule.Cookie.Secure

				// 关闭 httponly
				value.HttpOnly = rpRule.Cookie.HttpOnly
				if rpRule.Cookie.SameSite != 0 {
					value.SameSite = rpRule.Cookie.SameSite
				}
				if rpRule.Cookie.Domain != "" {
					value.Domain = rpRule.Cookie.Domain
				}
				mcook = append(mcook, value.String())
			}

			response.Header.Del("Set-Cookie")
			for _, mc := range mcook {
				response.Header.Add("Set-Cookie", mc)
			}
		} else {
			log.Info("[plugin] Cookie URL: %s cookies not find", response.Request.RequestURI)
		}
	}
	// 为空会自动删除
	// header 处理可能为nil
	if rpRule.Header != nil {
		for hkey, hvalue := range rpRule.Header {
			if hvalue == "" {
				response.Header.Del(hkey)
			} else {
				response.Header.Set(hkey, hvalue)
			}

		}
	}
	if rpRule.Body == nil {
		return nil
	}
	// file 处理 有 file 了就直接处理结束不要替换追加
	if strings.TrimSpace(rpRule.Body.File) != "" {
		// 建议加到缓存里面
		b := BodyFiles[rpRule.Body.File]
		response.Body = ioutil.NopCloser(bytes.NewReader(b))
		response.ContentLength = int64(len(b))
		response.Header.Set("Content-Length", strconv.Itoa(len(b)))
	}
	// 检查 Content-Length 如果大于 MaxContentLength 直接返回不处理当 maxCountlength 小于 0 放行
	if cleng == -1 && maxContentLength > 0 {
		return fmt.Errorf("[exit] ReplaceStr response Content-Length  no set,  set max: %d, rule Replace str and Append will ignore ", maxContentLength)
	}
	if cleng > maxContentLength && maxContentLength > 0 {
		return fmt.Errorf("[exit] ReplaceStr response Content-Length is :%s  set max: %d, rule Replace str and Append will ignore ", response.Header.Get("Content-Length"), maxContentLength)

	}
	conType := response.Header.Get("Content-Type")
	if conType == "" {
		log.Trace("%s,Content-Type is empty", response.Request.URL)
		return nil
	}
	if AllowMIMEType.Enable {
		//只允许文本类的替换
		if !utils.StrPrefixOrinList(conType, AllowMIMEType.List) {
			log.Trace("%s,Content-Type is not plan: %s will ignore", response.Request.URL, conType)
			return nil
		}
	}

	// append
	if rpRule.Body.Append != "" {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		log.Info("url :%s, payload: %s\n", response.Request.RequestURI, rpRule.Body.Append)
		body = append(body, rpRule.Body.Append...)

		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		response.Header.Set("Content-Length", fmt.Sprint(len(body)))
	}

	// replace 可能为nil
	if rpRule.Body.ReplaceStr != nil {
		// str 处理
		if len(rpRule.Body.ReplaceStr) > 0 {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}
			for _, str := range rpRule.Body.ReplaceStr {
				log.Info("[plugin] str module URL: %s, oldStr: %s, NewStr:%s", response.Request.RequestURI, str.Old, str.New)
				body = bytes.Replace(body, []byte(str.Old), []byte(str.New), str.Count)
			}
			response.Body = ioutil.NopCloser(bytes.NewReader(body))
			response.ContentLength = int64(len(body))
			response.Header.Set("Content-Length", strconv.Itoa(len(body)))
		}
	}
	return nil
}
