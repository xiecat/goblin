package reverse

import (
	"context"
	"crypto/tls"
	"fmt"
	"goblin/internal/plugin/replace"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"goblin/internal/options"
	"goblin/internal/plugin"
	"goblin/pkg/utils"

	log "unknwon.dev/clog/v2"
)

const (
	DefaultTimeOut = 5
)

func (s *Servers) InitServer(revConf *Reverse) *http.Server {

	mux := http.NewServeMux()
	//

	for uri, content := range plugin.StaticFiles {
		log.Info("set route: %s", uri)
		mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/javascript")
			w.Write(content)
		})
	}

	// 路径处理 https://zhuanlan.zhihu.com/p/102619401
	mux.Handle(s.StaticURI, http.StripPrefix(s.StaticURI, http.FileServer(http.Dir(s.StaticDir)))) // /static
	// 代理路由
	mux.HandleFunc("/", revConf.ServeHTTP)

	var muxMiddleware http.Handler = mux

	server := &http.Server{
		ReadTimeout:       s.ReadHeaderTimeout * time.Second,
		WriteTimeout:      s.WriteTimeout * time.Second,
		IdleTimeout:       s.IdleTimeout * time.Second,
		ReadHeaderTimeout: s.ReadHeaderTimeout * time.Second,
		Handler:           muxMiddleware,
		TLSConfig:         tlsConfig,
		TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	return server
}

func (s *Servers) startListeners() {
	for _, srvHTTP := range s.HTTP {
		go func(srv *http.Server) { log.Fatal("%v", srv.ListenAndServe()) }(srvHTTP)
	}

	for _, srvHTTPS := range s.HTTPS {
		go func(srv *http.Server) { log.Fatal("%v", srv.ListenAndServeTLS("", "")) }(srvHTTPS)
	}
}

func (s *Servers) shutdownServers(ctx context.Context) {
	for k, v := range s.HTTP {
		err := v.Shutdown(ctx)
		if err != nil {
			log.Error("Cannot shutdown server %s: %s\n", k, err)
		}
	}

	for k, v := range s.HTTPS {
		err := v.Shutdown(ctx)
		if err != nil {
			log.Fatal("Cannot shutdown server %s: %s", k, err)
		}
	}
}

func initReverse(options *options.Options) (revMap map[string]struct {
	SSL     bool
	Reverse *Reverse
	Listen  string
}) {
	var err error
	// 代理初始化
	if options.Proxy.ProxyServerAddr != "" {
		ProxyServerAddr, err = url.Parse(options.Proxy.ProxyServerAddr)
		if err != nil {
			log.Error("proxy resolve:%s", err.Error())
		}
	} else {
		ProxyServerAddr = nil
	}
	//todo 待重构处理
	// 设置代理请求头
	ProxyHeader = options.Server.ProxyHeader
	// 模板变量 StaticURI 初始化
	//plugin.PluginVariable.Static = options.Server.StaticURI
	// 设置日志
	logLevel = options.OutLog.LogLevel
	// 初始化证书
	tlsConfig.Certificates = []tls.Certificate{}
	// 初始化版本
	Version = options.VersionInfo
	plugin.Version = options.VersionInfo
	// 初始化 supportMIME
	replace.AllowMIMEType = options.SupportMIME
	revMap = make(map[string]struct {
		SSL     bool
		Reverse *Reverse
		Listen  string
	})
	// 默认证书配置
	certName, keyName := options.Proxy.CertDir+"/"+"default.crt", options.Proxy.CertDir+"/"+"default.key"
	// 检查 cert 默认证书是否存在
	if utils.FileExist(certName) && utils.FileExist(keyName) {
		tlsc, err := tls.LoadX509KeyPair(certName, keyName)
		if err != nil {
			log.Fatal(err.Error())
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, tlsc)
	} else {
		log.Fatal("default cert error")
	}
	fmt.Printf("Plugin Dir: %s\n", options.Proxy.PluginDir)
	for host, v := range options.Proxy.Sites {
		hAddr, port, err := utils.SplitHost(host)
		if err != nil {
			log.Fatal("SplitHost: %s err", v.ProxyPass)
		}
		portstr := strconv.Itoa(port)
		PluginTmpVar := &plugin.TmpVariable{}
		// 模板变量 StaticURI 初始化
		PluginTmpVar.Static = options.Server.StaticURI
		PluginTmpVar.FakeDomain = hAddr
		PluginTmpVar.FakePort = portstr
		PluginTmpVar.FakeHost = host
		// real
		realURL, err := utils.ParseWithScheme(v.ProxyPass)
		if err != nil {
			log.Fatal("url.Parse: %s err", v.ProxyPass)
		}
		// 模板变量初始化
		PluginTmpVar.ProxyPass = v.ProxyPass
		PluginTmpVar.RealDomain = realURL.Host
		PluginTmpVar.RealPort = realURL.Port
		PluginTmpVar.RealHost = fmt.Sprintf("%s:%s", realURL.Host, realURL.Port)
		rpRequestURI := realURL.RequestURI
		realURL.RequestURI = "/"
		PluginTmpVar.RealBaseURL = realURL.String()
		realURL.RequestURI = rpRequestURI

		if v.SSL {
			tlsc, err := tls.LoadX509KeyPair(options.Proxy.CertDir+"/"+v.CACert, options.Proxy.CertDir+"/"+v.CAKey)
			if err != nil {
				log.Fatal(err.Error())
			}
			tlsConfig.Certificates = append(tlsConfig.Certificates, tlsc)
			// 模板变量 FakeHost 初始化
			PluginTmpVar.FakeBaseURL = fmt.Sprintf("https://%s", host)

		} else {
			// 模板变量 FakeHost 初始化
			PluginTmpVar.FakeBaseURL = fmt.Sprintf("http://%s", host)
		}
		//代理地址
		if options.Proxy.ProxyServerAddr != "" {
			if v.SSL {
				// SSL 支持
				//log.Fatal("Temporarily does not support https Please wait")
				fmt.Printf("goblin: https://%s ==> [ proxy: %s ] ==> %s, Plugin: [ %v ]\n", host, options.Proxy.ProxyServerAddr, v.ProxyPass, v.Rules)
			} else {
				fmt.Printf("goblin: http://%s ==> [ proxy: %s ] ==> %s, Plugin: [ %v ]\n", host, options.Proxy.ProxyServerAddr, v.ProxyPass, v.Rules)
			}
		} else {
			if v.SSL {
				//todo SSL 支持
				//log.Fatal("Temporarily does not support https Please wait")
				fmt.Printf("goblin: https://%s ==> %s, Plugin: [ %v ]\n", host, v.ProxyPass, v.Rules)
			} else {
				fmt.Printf("goblin: http://%s  ==> %s, Plugin: [ %v ]\n", host, v.ProxyPass, v.Rules)
			}
		}

		// todo 应该为域名或者 host 不该为 IP:Port
		// revmap 填充数据
		if rev, ok := revMap[portstr]; ok {
			if port == 80 || port == 443 {
				rev.Reverse.AllowSite[hAddr] = v.ProxyPass
			}
			rev.Reverse.AllowSite[host] = v.ProxyPass
			rev.Reverse.AllowSite[hAddr] = v.ProxyPass

		} else {
			revMap[portstr] = struct {
				SSL     bool
				Reverse *Reverse
				Listen  string
			}{
				SSL:    v.SSL,
				Listen: v.ListenIP,
				Reverse: &Reverse{
					AllowSite: map[string]string{
						host:  v.ProxyPass, // todo 尽量移除
						hAddr: v.ProxyPass,
					},
					HostProxy:             make(map[string]*httputil.ReverseProxy),
					MaxIdleConns:          options.Proxy.MaxIdleConns,
					MaxIdleConnsPerHost:   options.Proxy.MaxIdleConnsPerHost,
					MaxConnsPerHost:       options.Proxy.MaxConnsPerHost,
					IdleConnTimeout:       options.Proxy.IdleConnTimeout,
					TLSHandshakeTimeout:   options.Proxy.TLSHandshakeTimeout,
					ExpectContinueTimeout: options.Proxy.ExpectContinueTimeout,
					MaxContentLength:      options.Proxy.MaxContentLength,
					DingTalk:              options.Notice.DingTalk,
				}}

		}
		if v.Rules != "" {
			rule, err := plugin.LoadPlugin(options.Proxy.PluginDir + "/" + v.Rules + ".yaml")
			if err != nil {
				log.Fatal("plugin err please check: %s", v.Rules)

			}
			// 初始化插件配置
			rule.SetInitConfig(PluginTmpVar)
			plugin.Plugins[host] = rule
			plugin.Plugins[hAddr] = rule

		}
	}
	return revMap
}

func InitServerConfig(options *options.Options) *Servers {
	// cache 初始化
	cacheRspFile.Type = options.CacheType
	cacheRspFile.Size = options.CacheSize

	// server 初始化
	revMap := initReverse(options)
	servers := &Servers{
		HTTP:              make(map[string]*http.Server),
		HTTPS:             make(map[string]*http.Server),
		ReadTimeout:       options.Server.ReadTimeout,
		WriteTimeout:      options.Server.WriteTimeout,
		IdleTimeout:       options.Server.IdleTimeout,
		ReadHeaderTimeout: options.Server.ReadHeaderTimeout,
		StaticURI:         options.Server.StaticURI,
		StaticDir:         options.Server.StaticDir,
	}
	for portstr, revconf := range revMap {
		if revconf.SSL {
			servers.HTTPS[portstr] = servers.InitServer(revconf.Reverse)
			servers.HTTPS[portstr].Addr = revconf.Listen + ":" + portstr
			fmt.Printf("ListenServer: https://%s\n", servers.HTTPS[portstr].Addr)
		} else {
			servers.HTTP[portstr] = servers.InitServer(revconf.Reverse)
			servers.HTTP[portstr].Addr = revconf.Listen + ":" + portstr
			fmt.Printf("ListenServer: http://%s\n", servers.HTTP[portstr].Addr)
		}
	}

	return servers
}

func (s *Servers) Start() {
	s.startListeners()
}

func (s *Servers) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeOut*time.Second)
	defer cancel()

	log.Warn("Shutting down servers...")
	s.shutdownServers(ctx)
}
