package dump

var Method = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "TRACE"}

type Dump struct {
	Request  *Request  `yaml:"Request"`
	Response *Response `yaml:"Response"`
	Notice   bool
}

// Request  请求头
type Request struct {
	Method []string `yaml:"Method"`
}

// Response 响应头
type Response struct {
	Status int               `yaml:"Status"`
	Header map[string]string `yaml:"Header"`
	Body   string            `yaml:"Body"`
}
