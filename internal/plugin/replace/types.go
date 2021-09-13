package replace

import (
	"net/http"
)

var Method = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "TRACE"}

var AllowMIMEType = &SupportMIME{}
var (
	BodyFiles = make(map[string][]byte) //插件系统注入的js静态文件
)

type Replace struct {
	Request  *Request  `yaml:"Request"`
	Response *Response `yaml:"Response"`
}

// Request  请求头
type Request struct {
	Method []string          `yaml:"Method"`
	Header map[string]string `yaml:"Header"`
}

// Response 响应头
type Response struct {
	Status   int               `yaml:"Status"`
	Header   map[string]string `yaml:"Header"`
	Cookie   *Cookie           `yaml:"Cookie"`
	Body     *Body             `yaml:"Body"`
	Location string            `yaml:"Location"`
}

type Cookie struct { // 由于有默认值需要一起设置
	Domain   string // optional
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

// Body 响应体
type Body struct {
	File       string        `yaml:"File"`
	ReplaceStr []*ReplaceStr `yaml:"ReplaceStr"`
	Append     string        `yaml:"Append"`
}

// ReplaceStr 替换字符串
type ReplaceStr struct {
	Old   string `yaml:"Old"`
	New   string `yaml:"New"`
	Count int    `yaml:"Count"`
}

type SupportMIME struct {
	Enable bool     `yaml:"Enable"`
	List   []string `yaml:"List"`
}
