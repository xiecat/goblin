package options

import (
	"goblin/internal/plugin/replace"
	"goblin/pkg/cache"
	"goblin/pkg/ipinfo"
	"goblin/pkg/logging"
	"goblin/pkg/notice"
	"plugin"
	"time"
)

const (
	MiniLen = 10
)

type Options struct {
	Version     bool   `yaml:"-"`
	VersionInfo string `yaml:"-"`
	BinDir      string `yaml:"-"`
	// write config to file
	WConfile   bool         `yaml:"-"`
	Loglevel   int          `yaml:"Loglevel"`
	ConfigFile string       `yaml:"-"`
	GenPOC     string       `yaml:"-"`
	TestNotice bool         `yaml:"-"`
	Server     *Server      `yaml:"Server"`
	Proxy      *Proxy       `yaml:"Proxy"`
	Notice     noticeConfig `yaml:"Notice"`
	IPLocation *ipinfo.Config
	LogFile    string        `yaml:"log_file"`
	Cache      *cache.Config `yaml:"cache"`
	// PrintConfig print config file
	PrintConfig bool                 `yaml:"-"`
	CacheType   []string             `yaml:"CacheType"`
	CacheSize   int64                `yaml:"CacheSize"`
	Plugin      []*plugin.Plugin     `yaml:"-"`
	SupportMIME *replace.SupportMIME `yaml:"SupportMIME"`
	OutLog      *logging.Config      `yaml:"OutLog"`
}

type noticeConfig struct {
	DingTalk *notice.DingTalk
}

type Server struct {
	IdleTimeout       time.Duration `yaml:"IdleTimeout"`
	ReadTimeout       time.Duration `yaml:"ReadTimeout"`
	WriteTimeout      time.Duration `yaml:"WriteTimeout"`
	ReadHeaderTimeout time.Duration `yaml:"ReadHeaderTimeout"`
	ProxyHeader       string        `yaml:"ProxyHeader"`
	StaticDir         string        `yaml:"StaticDir"`
	StaticURI         string        `yaml:"StaticURI"`
}
