package options

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"goblin/internal/plugin"
	"goblin/internal/plugin/replace"
	"goblin/pkg/cache"
	"goblin/pkg/cache/redis"
	"goblin/pkg/ipinfo"
	"goblin/pkg/logging"
	"goblin/pkg/notice"
	"goblin/pkg/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
	log "unknwon.dev/clog/v2"
)

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	var options = &Options{
		ConfigFile:  "goblin.yaml",
		VersionInfo: Version,
		LogFile:     "goblin.log",
		OutLog: &logging.Config{
			Type:     "file",
			LogLevel: logrus.InfoLevel,
			FileLog: &logging.FileLog{
				Mode: "text",
				DSN:  "access.log",
			},
			EsLog: &logging.EsLog{
				LogLevel: logrus.InfoLevel,
				DSN:      "http://127.0.0.1:9200",
				Host:     "localhost",
				Index:    "goblin",
			},
			Syslog: &logging.Syslog{
				DSN:  "127.0.0.1:514",
				Mode: "text",
			},
		},
		Cache: &cache.Config{
			Type:    "self",
			ExpTime: 10 * time.Minute,
			Redis: &redis.Config{
				Host:     "127.0.0.1",
				Port:     6379,
				Password: utils.RandChar(MiniLen),
			},
		},
		Server: &Server{
			IdleTimeout:       180 * time.Second,
			ReadTimeout:       300 * time.Second,
			WriteTimeout:      300 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			ProxyHeader:       "RemoteAddr",
			StaticDir:         "static",
			StaticURI:         "/" + strings.ToLower(utils.RandChar(MiniLen)) + "/",
		},
		Loglevel: 2,
		IPLocation: &ipinfo.Config{
			Type: "qqwry",
		},
		Notice: noticeConfig{
			DingTalk: &notice.DingTalk{
				URL: "",
			},
		},
		CacheType: []string{"png", "jpg", "js", "jpeg", "css", "otf", "ttf"},
		CacheSize: 1024 * 1024 * 12, // 缓存不超过 12M
		BinDir:    utils.BinBaseDir(),
		Proxy: &Proxy{
			MaxIdleConns:          512,
			MaxConnsPerHost:       20,
			MaxIdleConnsPerHost:   20,
			IdleConnTimeout:       120 * time.Second,
			TLSHandshakeTimeout:   60 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxContentLength:      -1, // 默认不检查
			ProxyCheckURL:         "https://myip.ipip.net/",
			PluginDir:             "plugins",
			CertDir:               "cert",
			Sites: ProxySite{
				"127.0.0.1:8084": {
					ListenIP:     "0.0.0.0",
					SSL:          false,
					CACert:       "",
					CAKey:        "",
					ProxyPass:    "https://www.baidu.com",
					StaticPrefix: strings.ToLower(utils.RandChar(MiniLen)),
					Rules:        "demo",
				},
				"127.0.0.1:8083": {
					ListenIP:     "0.0.0.0",
					SSL:          false,
					CACert:       "",
					CAKey:        "",
					ProxyPass:    "https://www.douban.com/",
					StaticPrefix: strings.ToLower(utils.RandChar(MiniLen)),
				},
			},
		},
		SupportMIME: &replace.SupportMIME{
			Enable: false,
			List: []string{"text", "application/json", "application/javascript",
				"application/x-javascript", "message", "application/hta",
				"application/rtf", "application/ecmascript", "image/svg+xml",
				"application/xhtml", "application/xml"},
		},
	}
	// 显示banner
	showBanner()
	// 配置的优先级是 命令行 <- 配置文件 <- 默认配置
	// 如果不存在配置文件则创建
	options.makeConfig()
	// 获取配置文件的路径
	options.readFileconfigure()
	// 读取配置文件配置
	err := options.readConfigFile()
	if err != nil {
		os.Exit(1)
	}
	// 读取命令行中的配置
	options.readConfigCMD()
	// 如果需要刷新配置文件可以使用 -w 参数
	if options.WConfile {
		options.writeConfig()
		os.Exit(0)
	}
	logging.AccLogger = options.OutLog.New()

	// 检查 cert 插件目录是否存在
	if !utils.FileExist(options.Proxy.CertDir) {
		fmt.Printf("Plugin not exist will create: %s", options.Proxy.CertDir)
		err = os.Mkdir(options.Proxy.CertDir, 755) //nolint:
		if err != nil {
			log.Fatal("%s", err.Error())
		}
		//https://www.cnblogs.com/feiquan/p/11429065.html
		os.Chmod(options.Proxy.CertDir, 0755) //nolint:
	}
	//
	certName, KeyName := options.Proxy.CertDir+"/"+"default.crt", options.Proxy.CertDir+"/"+"default.key"
	// 检查 cert 默认证书是否存在
	if !utils.FileExist(certName) || !utils.FileExist(KeyName) {
		utils.MakeDefautCert(certName, KeyName)
	}

	// 检查 plugin 插件目录是否存在
	if !utils.FileExist(options.Proxy.PluginDir) {
		fmt.Printf("Plugin not exist will create: %s", options.Proxy.PluginDir)
		err = os.Mkdir(options.Proxy.PluginDir, 755) //nolint:
		if err != nil {
			log.Fatal("%s", err.Error())
		}
		//https://www.cnblogs.com/feiquan/p/11429065.html
		os.Chmod(options.Proxy.PluginDir, 0755) //nolint:
		// plugin pam.yaml
		demo := `Name: demo
Version: 0.0.1
Description: this is a description
WriteDate: "2021-09-06"
Author: goblin
Rule:
- url: /
  Match: Word
  Replace: ## 替换模块
    - Request:
        Method: ## 匹配到如下请求方式方可替换
          - GET
        Header:
          goblin: 1.0.1  # 替换的 header 头内容。为空则是删除。
      Response: # 替换的响应内容
        Body:
          Append: <script type='text/javascript'>setTimeout(function(){alert("hello goblin!");}, 2000);</script> # 追加字符串`
		pluginDemo := options.Proxy.PluginDir + "/" + "demo.yaml"
		err = ioutil.WriteFile(pluginDemo, []byte(demo), 0755) //nolint:
		if err != nil {
			log.Fatal("%s", err.Error())
		}
		// https://www.cnblogs.com/feiquan/p/11429065.html
		err = os.Chmod(pluginDemo, 0755) // nolint:
		if err != nil {
			log.Fatal("%s", err.Error())
		}
	}

	// 检查 static server 目录是否存在
	if !utils.FileExist(options.Server.StaticDir) {
		fmt.Printf("static server dir not exist will create: %s\n", options.Server.StaticDir)
		err = os.Mkdir(options.Server.StaticDir, 0755) //nolint: 权限不需要检查
		if err != nil {
			log.Fatal("%s", err.Error())
		}
		// https://www.cnblogs.com/feiquan/p/11429065.html
		err = os.Chmod(options.Server.StaticDir, 0755)
		if err != nil {
			log.Fatal("%s", err.Error())
		}
	}
	// 静态服务器器写入 index.html
	sIndex := options.Server.StaticDir + "/index.html"
	if !utils.FileExist(sIndex) {
		log.Info("will create file index.html")
		err = ioutil.WriteFile(sIndex, []byte(""), 0755) //nolint:
		if err != nil {
			log.Fatal("%s", err.Error())
		}
		// https://www.cnblogs.com/feiquan/p/11429065.html
		err = os.Chmod(sIndex, 0755) // nolint:
		if err != nil {
			log.Fatal("%s", err.Error())
		}
	}
	fmt.Printf("Static Server route: %s ==> file dir: %s\n", options.Server.StaticURI, options.Server.StaticDir)
	// 检查配置
	options.validateOptions()
	return options
}

// readConfigCMD 读配置文件
func (options *Options) readConfigCMD() {
	flag.StringVar(&options.LogFile, "log", options.LogFile, "Webserver log file")
	flag.StringVar(&options.ConfigFile, "config", options.ConfigFile, "Webserver port")
	flag.StringVar(&options.GenPOC, "gen-plugin", options.GenPOC, "Generate rule file")
	flag.BoolVar(&options.TestNotice, "test-notice", false, "Test message alarm")
	flag.BoolVar(&options.WConfile, "w", options.WConfile, "Write config to config file")
	flag.BoolVar(&options.Version, "v", false, "Show version of goblin")
	flag.IntVar(&options.Loglevel, "log-level", options.Loglevel, "Log mode [1-5] 1.dump All logs include GET log and POST log, 2. Record POST log, 3. Record dump log in rules, 4. Record error log, and 5. Record exception exit log") //nolint:
	flag.BoolVar(&options.PrintConfig, "print-config", false, "print config file")
	flag.Parse()

	if options.Version {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Branch: %s\n", Branch)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("ResDate: %s\n", Release)
		os.Exit(0)
	}
	if options.GenPOC != "" {
		poc := options.GenPOC
		if !strings.HasSuffix(poc, "yaml") {
			poc = options.GenPOC + ".yaml"
		}
		fmt.Printf("Generate plug-in: %s, please edit it and put it in the [ plugins ] directory\n", poc)
		err := plugin.GenDefaultPlugin(poc)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
		}
		os.Exit(0)
	}
	if options.TestNotice {
		log.Warn("this is a test notice")
		options.Notice.DingTalk.SendTest()
		os.Exit(0)
	}
	if options.PrintConfig {
		options.printConfig()
		os.Exit(0)
	}
}

// readConfigFile 加载配置文件
func (options *Options) readConfigFile() error {
	if len(options.Proxy.Sites) > 0 {
		options.Proxy.Sites = make(ProxySite)
	}
	path := options.ConfigFile
	configFile, err := os.Open(path)
	if err != nil {
		log.Fatal("ConfigFile: %s os.Open failed: %v", path, err)
	}
	defer configFile.Close()

	content, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal("ConfigFile: %s ioutil.ReadAll failed: %v", path, err)
		return err
	}

	err = yaml.Unmarshal(content, options)
	if err != nil {
		log.Fatal("readConfigFile(%s) yaml.Unmarshal failed: %v", path, err)
		return err
	}
	return nil
}

// readFileconfigure 找到配置文件
func (options *Options) readFileconfigure() {
	for i, c := range os.Args {
		// 如果指定了 -config, 它必须是第一个参数
		if strings.HasPrefix(c, "-config") && i != 1 {
			log.Fatal("-config must be the first arg")
		}
		// 等号两边请不要加空格
		if c == "=" {
			// -config = goblin.yaml not support
			log.Fatal("wrong format, no space between '=', eg: -config=goblin.yaml")
		}
	}
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-config") {
		if os.Args[1] == "-config" && len(os.Args) > 2 {
			if os.Args[2] == "=" && len(os.Args) > 3 {
				// -config = goblin.yaml not support
				log.Fatal("wrong format, no space between '=', eg: -config=goblin.yaml")
			} else {
				// -config goblin.yaml
				options.ConfigFile = os.Args[2]
			}
			if strings.HasPrefix(options.ConfigFile, "=") {
				//nolintlint -config =goblin.yaml
				options.ConfigFile = strings.Split(options.ConfigFile, "=")[1]
			}
		}
		if strings.Contains(os.Args[1], "=") {
			//nolintlint  -config=goblin.yaml
			options.ConfigFile = strings.Split(os.Args[1], "=")[1]
		}
	} else {
		for i, c := range os.Args {
			if strings.HasPrefix(c, "-config") && i != 1 {
				log.Fatal("-config must be the first arg")
			}
		}
	}
}

// makeConfig 生成配置文件
func (options *Options) makeConfig() {
	if !utils.FileExist(options.ConfigFile) {
		data, _ := yaml.Marshal(options)
		log.Info("can't find config file for goblin in %s. and will write config file\n", options.ConfigFile)
		err := ioutil.WriteFile(options.ConfigFile, data, 0644) //nolint: 写入
		if err != nil {
			log.Fatal("%s can't  write config file\n", options.ConfigFile)
		}
	}
}

// writeConfig 写入配置文件
func (options *Options) writeConfig() {
	data, _ := yaml.Marshal(options)
	fmt.Printf("find config file for goblin in %s. and will update config file\n", options.ConfigFile)
	err := ioutil.WriteFile(options.ConfigFile, data, 0644) //nolint: 写入
	if err != nil {
		log.Fatal("%s can't  write config file\n", options.ConfigFile)
	}
}

// printConfig 打印配置文件
func (options *Options) printConfig() {
	if utils.FileExist(options.ConfigFile) {
		data, _ := yaml.Marshal(options)
		fmt.Println(string(data))
	} else {
		log.Fatal("can't find config file for goblin in %s.\n", options.ConfigFile)
	}
}
func (options *Options) SetLogLevel() (logLevel log.Level) {
	switch options.Loglevel {
	case 1: //nolint:
		logLevel = log.LevelTrace
	case 2: //nolint:
		logLevel = log.LevelInfo
	case 3: //nolint:
		logLevel = log.LevelWarn
	case 4: //nolint:
		logLevel = log.LevelError
	case 5: //nolint:
		logLevel = log.LevelFatal
	default:
		logLevel = log.LevelInfo
	}
	return logLevel
}
