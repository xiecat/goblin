package plugin

import (
	"goblin/internal/plugin/dump"
	"goblin/internal/plugin/inject"
	"goblin/internal/plugin/replace"
)

var (
	// todo 解耦
	Plugins = make(map[string]*BasePlugin)
	// MatchType URL 匹配类型
	MatchType          = []string{"Word", "Prefix", "Suffix"}
	StaticFiles        = make(map[string][]byte) //插件系统注入的js静态文件
	Version     string = "unknown"
)

// BasePlugin 插件结构体
type BasePlugin struct {
	Name        string  `yaml:"Name"`
	Version     string  `yaml:"Version"`
	Description string  `yaml:"Description"`
	WriteDate   string  `yaml:"WriteDate"`
	Author      string  `yaml:"Author"`
	Rule        []*Rule `yaml:"Rule"`
	UseBody     bool    `yaml:"-"` //响应body解包一次
}

// Rule 规则结构体
type Rule struct {
	URL      string
	Match    string             `yaml:"Match"`
	Replace  []*replace.Replace `yaml:"Replace"`
	Dump     []*dump.Dump       `yaml:"Dump"`
	InjectJs *inject.InjectJs   `yaml:"InjectJs"`
}

// 模板变量
type TmpVariable struct {
	Static string

	FakeBaseURL string // http://www.baidu.com/
	FakeDomain  string //ip/domain
	FakePort    string // port
	FakeHost    string // ip:port

	RealBaseURL string // http://www.baidu.com/
	RealDomain  string
	RealHost    string
	RealPort    string
	ProxyPass   string
}
