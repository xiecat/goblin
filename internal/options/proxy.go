package options

import (
	"time"
)

type Proxy struct {
	MaxIdleConns          int           `yaml:"MaxIdleConns"`
	MaxIdleConnsPerHost   int           `yaml:"MaxIdleConnsPerHost"`
	MaxConnsPerHost       int           `yaml:"MaxConnsPerHost"`
	IdleConnTimeout       time.Duration `yaml:"IdleConnTimeout"`
	TLSHandshakeTimeout   time.Duration `yaml:"TLSHandshakeTimeout"`
	ExpectContinueTimeout time.Duration `yaml:"ExpectContinueTimeout"`
	MaxContentLength      int           //处理响应体的最大长度
	ProxyServerAddr       string        `yaml:"ProxyServerAddr"` // socks5://
	ProxyCheckURL         string        `yaml:"ProxyCheckURL"`
	PluginDir             string        `yaml:"PluginDir"`
	CertDir               string        `yaml:"CertDir"`
	Sites                 ProxySite     `yaml:"Site"`
}

type ProxySite map[string]struct {
	ListenIP     string `yaml:"Listen"`
	StaticPrefix string `yaml:"StaticPrefix"`
	SSL          bool   `yaml:"SSL"` // http https
	CAKey        string `yaml:"CAKey"`
	CACert       string `yaml:"CACert"`
	ProxyPass    string `yaml:"ProxyPass"`
	Rules        string `yaml:"Plugin"`
}
