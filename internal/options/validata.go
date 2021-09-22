package options

import (
	"crypto/tls"
	"fmt"
	"goblin/internal/plugin"
	"goblin/pkg/utils"
	"strings"
	log "unknwon.dev/clog/v2"
)

func (options *Options) validateOptions() {
	// Both verbose and silent flags were used
	options.validateLogLevel()
	options.validateCacheDsn()
	options.validateIpLocation()
	options.validataProxySite()
	options.Server.validateServer()
	options.Proxy.validatePlugin()
}

func (s *Server) validateServer() {
	// Both verbose and silent flags were used
	if len(s.StaticURI) < MiniLen {
		log.Fatal("StaticURI：%s must be no less than %d.  You can use it instead: %s ", s.StaticURI, MiniLen, utils.RandChar(MiniLen))
	}
	if !strings.HasPrefix(s.StaticURI, "/") || !strings.HasSuffix(s.StaticURI, "/") {
		log.Fatal("StaticURI：%s Must '/' begin, '/' end. For example /xadef34rxdd/ ", s.StaticURI)
	}
}
func (options *Options) validataProxySite() {
	// 判断端口不重复
	pSites := options.Proxy
	httpPort := map[int]string{}
	httpsPort := map[int]string{}
	checkPrefixMap := make(map[string]bool)
	for host, v := range pSites.Sites {
		// todo 检查 proxy 超时参数

		addr, port, err := utils.SplitHost(host)
		if err != nil {
			log.Fatal("Proxy host: %s  format error must have port eg(127.0.0.1:80,www.xxx.com:80)", addr)
		}
		//if !utils.IsHost(addr) {
		//	log.Fatal("Proxy host is not ip or host please check host: %s,Addr: %s", host, addr)
		//}
		if !utils.IsPort(port) {
			log.Fatal("Port range of (0,65535] host: %s", host)
		}

		if !utils.IsHost(v.ListenIP) {
			log.Fatal("Listen is not ip or host please check host: %s Listen: %s", host, v.ListenIP)
		}
		// ssl 检查
		if v.SSL {
			certName, keyName := options.Proxy.CertDir+"/"+v.CACert, options.Proxy.CertDir+"/"+v.CAKey
			if !utils.FileExist(certName) {
				log.Fatal("Host: %s ,CACert File not find: %s/%s", host, certName)
			}
			if !utils.FileExist(keyName) {
				log.Fatal("Host: %s ,Cakey File not find: %s", host, keyName)
			}
			_, err := tls.LoadX509KeyPair(certName, keyName)
			if err != nil {
				log.Fatal("cert format err: %s", err.Error())
			}

			httpsPort[port] = host
		} else {
			httpPort[port] = host
		}
		// 检查 pass
		if !utils.IsURL(v.ProxyPass) { // isURL 是否是 URL // 也可以使用 url 请求
			log.Fatal("ProxyPass Error! Host: %s ,ProxyPass: %s", host, v.ProxyPass)
		}
		if len(v.StaticPrefix) < MiniLen {
			log.Fatal("Host: %s, StaticPrefix must be no less than %d.  You can use it instead: %s ", host, MiniLen, utils.RandChar(MiniLen))
		}
		// 校验Static 是否唯一
		if v.StaticPrefix == "" {
			log.Fatal("Host: %s, StaticPrefix must not be empty. You can use it instead: %s ", host, utils.RandChar(MiniLen))
		}
		if _, ok := checkPrefixMap[v.StaticPrefix]; ok {
			log.Fatal("Host: %s, StaticPrefix:%s must not be repeated. You can use it instead:%s ", host, v.StaticPrefix, utils.RandChar(MiniLen))
		}

	}
	// 检查代理
	if pSites.ProxyServerAddr != "" {
		log.Trace("check proxy:%s", pSites.ProxyServerAddr)
		if !utils.ValidProxy(pSites.ProxyServerAddr, pSites.ProxyCheckURL) {
			log.Fatal("ProxyServerAddr:%s ", pSites.ProxyServerAddr)
		} else {
			fmt.Printf("use proxy: %s\n", pSites.ProxyServerAddr)
		}
	}
	//检查端口不匹配
	for k, v := range httpsPort {
		if p, ok := httpPort[k]; ok {
			log.Fatal("port not match http port:%s https port: %s", v, p)
		}
	}
	//检查端口不匹配
	for k, v := range httpPort {
		if p, ok := httpsPort[k]; ok {
			log.Fatal("port not match http port:%s https port: %s", v, p)
		}
	}

}

// validateCacheDsn 验证缓存的配置信息
func (options *Options) validateCacheDsn() {
	if err := options.Cache.ValidateCacheDsn(); err != nil {
		log.Fatal(err.Error())
	}
}

// validateHost 验证站点绑定的地址信息
func (options *Options) validateIpLocation() {
	options.IPLocation.ValidateType()
}

// validateHost 验证站点绑定的地址信息
func (options *Options) validateLogLevel() {
	if options.Loglevel < 1 || options.Loglevel > 5 {
		log.Fatal("logLevel range [1,5].and you set: %d", options.Loglevel)
	}
	err := options.OutLog.ValidateType()
	if err != nil {
		log.Fatal(err.Error())
	}
}

// validatePlugin 验证使用的插件是否正确
func (p *Proxy) validatePlugin() {
	for host, v := range p.Sites {
		if v.Rules == "" {
			continue
		}
		pRuleDir := p.PluginDir + "/" + v.Rules + ".yaml"
		if !utils.FileExist(pRuleDir) {
			log.Fatal("Host: %s, Plugin: %s not exist,please check\n", host, pRuleDir)
		}
		pg, err := plugin.LoadPlugin(pRuleDir)
		if err != nil {
			log.Fatal("LoadPlugin: %s", err.Error())
		}
		err = pg.CheckPlugin()
		if err != nil {
			log.Fatal("Plugin Check: %s", err.Error())
		}
	}
}
