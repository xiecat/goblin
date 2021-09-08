package reverse

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"goblin/internal/msgpack"
	"goblin/pkg/cache"
	"goblin/pkg/utils"

	log "unknwon.dev/clog/v2"
)

func encodeCache(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		log.Info("no cache url:%s status_code %d", response.Request.URL.String(), response.StatusCode)
		return nil
	}
	req := response.Request
	currUrl := req.URL.String()

	if cacheRspFile.search(currUrl) { // 判断 后缀类型
		// 读取 body
		body, err := ioutil.ReadAll(response.Body)
		// 将 body 放入
		response.Body = ioutil.NopCloser(bytes.NewReader(body))
		if int64(len(body)) > cacheRspFile.Size {
			return nil
		}
		if err != nil {
			log.Error("currurl: %s,err:%s", currUrl, err.Error())
		}
		if _, err := cache.GlobalCache.Get(currUrl); err != nil {
			log.Trace("cache type %s, Content-Type: %s, cached file: %s", cache.GlobalCache.Type, response.Header.Get("Content-Type"), currUrl)
			// 如果为 self 不做编码处理
			if cache.GlobalCache.Type == "self" {
				cache.GlobalCache.SetNX(currUrl, &URIObj{
					URL:             *req.URL,
					Method:          req.Method,
					RequestHeaders:  req.Header,
					ResponseHeaders: response.Header,
					StatusCode:      response.StatusCode,
					Content:         body,
				})
				return nil
			}

			// 使用 redis 等外部需要编码
			value, err := msgpack.Encode(&URIObj{
				URL:             *req.URL,
				Method:          req.Method,
				RequestHeaders:  req.Header,
				ResponseHeaders: response.Header,
				StatusCode:      response.StatusCode,
				Content:         body,
			})
			if err != nil {
				log.Error(err.Error())
				return err
			}
			encoded, err := utils.B64Encode(value)
			if err != nil {
				return err
			}
			cache.GlobalCache.SetNX(currUrl, encoded)
		}
	}
	return nil
}

func decodeCache(currURL string) (*URIObj, error) {
	if ff, ok := cache.GlobalCache.Get(currURL); ok == nil {
		if cache.GlobalCache.Type == "self" {
			uff, ok := ff.(*URIObj)
			if !ok {
				return nil, errors.New("the data decode error")
			}
			return uff, nil
		}
		decoded, err := utils.B64Decode(ff.(string))
		if err != nil {
			return nil, err
		}
		urobj := URIObj{}
		err = msgpack.Decode(decoded, &urobj)
		return &urobj, err
	}
	return nil, errors.New("no cache")
}
