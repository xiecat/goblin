package reverse

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"goblin/pkg/notice"

	"github.com/valyala/bytebufferpool"
)

// Bufferpool 反向代理字符串缓存池
var (
	Bufferpool      = &bufferPool{ByteBuffer: &bytebufferpool.ByteBuffer{}}
	cacheRspFile    = &CacheFile{}
	ProxyServerAddr *url.URL
	ProxyHeader     = "RemoteAddr"
	logLevel        = logrus.InfoLevel
	tlsConfig       = &tls.Config{}
	Version         = "unknown"
)

// Reverse 反代配置
type Reverse struct {
	// proxy =>pass
	AllowSite             map[string]string
	HostProxy             map[string]*httputil.ReverseProxy
	MaxContentLength      int           //处理响应体的最大长度
	MaxIdleConns          int           //最大空闲连接
	MaxConnsPerHost       int           //每个host的最大连接数量
	MaxIdleConnsPerHost   int           //每个host的连接池最大空闲连接数,默认2
	IdleConnTimeout       time.Duration //空闲超时时间
	TLSHandshakeTimeout   time.Duration //tls握手超时时间
	ExpectContinueTimeout time.Duration
	DingTalk              *notice.DingTalk
}

// Servers 服务配置
type Servers struct {
	HTTP              map[string]*http.Server
	HTTPS             map[string]*http.Server
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	StaticDir         string
	StaticURI         string
}

type URIObj struct {
	URL             url.URL
	Method          string
	StatusCode      int
	RequestHeaders  http.Header
	ResponseHeaders http.Header
	Content         []byte
}

type CacheFile struct {
	Type []string
	Size int64
}

func (ca *CacheFile) search(matchURL string) bool {
	for _, v := range ca.Type {
		if strings.HasSuffix(matchURL, v) {
			return true
		}
	}
	return false
}

type bufferPool struct {
	*bytebufferpool.ByteBuffer
}

func (b *bufferPool) Get() []byte {
	return b.Bytes()
}

func (b *bufferPool) Put(payload []byte) {
	b.Set(payload)
}
