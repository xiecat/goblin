# :fishing_pole_and_fish: Goblin 钓鱼演练工具

![GitHub branch checks state](https://img.shields.io/github/checks-status/xiecat/goblin/master)
[![Latest release](https://img.shields.io/github/v/release/xiecat/goblin)](https://github.com/xiecat/goblin/releases/latest)
![GitHub Release Date](https://img.shields.io/github/release-date/xiecat/goblin)
![GitHub All Releases](https://img.shields.io/github/downloads/xiecat/goblin/total)
[![GitHub issues](https://img.shields.io/github/issues/xiecat/goblin)](https://github.com/xiecat/goblin/issues)
[![Docker Pulls](https://img.shields.io/docker/pulls/becivells/goblin)](https://hub.docker.com/r/becivells/goblin)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/becivells/goblin)        
Goblin 是一款适用于红蓝对抗的钓鱼演练工具。通过反向代理，可以在不影响用户操作的情况下无感知的获取用户的信息，或者诱导用户操作。也可以通过使用代理方式达到隐藏服务端的目的。内置插件，通过简单的配置，快速调整网页内容以达到更好的演练效果

[:ledger:English Document](https://github.com/xiecat/goblin/blob/master/README_EN.md)   |   [:pushpin:下载地址](https://github.com/xiecat/goblin/releases)    |   [:book:使用文档](https://xiecat.github.io/goblin-doc/)

## :collision: 特点: 

* 支持缓存静态文件，加速访问
* 支持 dump 所有请求，dump 匹配规则的请求
* 支持通过插件快速配置，调整不合适的跳转或者内容
* 支持植入特定的 js
* 支持修改响应内容或者 goblin 请求的内容
* 支持通过代理方式隐藏真实 IP


## :tv: Demo:

demo效果演示：
![image](https://github.com/xiecat/goblin/blob/master/Demo.gif)

快速体验 demo
1. Flash demo
```shell
docker run -it --rm  -p 8083:8083 -p 8084:8084 -p 8085:8085 -p 8086:8086  becivells/goblin-demo-flash
```
本机访问 [http://127.0.0.1:8083](http://127.0.0.1:8083) 示例仓库为: [goblin-flash-demo](https://github.com/xiecat/goblin-demo/tree/master/goblin-demo-flash)

2. 默认代理百度的 demo
```shell
docker run -it --rm -v $(pwd):/goblin/ -p 8084:8084 becivells/goblin
```

本机访问 [http://127.0.0.1:8084](http://127.0.0.1:8084)

## :computer: 快速部署


### Docker 快速部署

运行如下命令获取镜像
```shell
docker pull becivells/goblin
```
Dockerfile 如下：
```shell
FROM scratch
COPY goblin /usr/bin/goblin
ENTRYPOINT ["/usr/bin/goblin"]
WORKDIR /goblin
```
工作目录在 goblin ，首先创建目录，切换到目录下，执行
```shell
docker run -it --rm -v $(pwd):/goblin/ -p 8084:8084 becivells/goblin
```


### Git安装

1.访问 [https://github.com/xiecat/goblin/releases](https://github.com/xiecat/goblin/releases) 从中选择适合自己操作系统的二进制文件（注:本系统全面支持国产芯片,相关文件可进微信群获取，进群二维码见文末）

2.根据需求修改配置文件的参数，配置文件详细介绍请移步使用文档 [:point_right:配置文件介绍](https://xiecat.github.io/goblin-doc/config/)

命令行参数如下

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
#### :warning: 注意

如果是在服务器端部署则需要修改 ip 地址

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

## :triangular_ruler: 插件系统


插件系统介绍详细使用方法见文档[:point_right:插件系统](https://xiecat.github.io/goblin-doc/plugin/)


## :battery: 高阶用法

goblin 使用反向代理，前端使用 cf 等代理 goblin， 即可隐藏 goblin 主机

### JS 注入

js 注入有两种方式一种是跟着页面走(Replace 需要自己追加` \<script\> ` 标签)，一种是跟着全局 js 文件走各有好处。

这两种其实都是使用 Replace 功能

#### 使用 InjectJs 注入

```yaml
- url: /base.js # 待替换的js 尽量选择全局 js
  Match: word   # 匹配方式
  InjectJs:
    File: aaa.js # 要替换的 js。 可以为文件或者 url
```

#### 使用 replace 注入

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



## :star: 案例

### [深信服 vpn 案例](https://xiecat.github.io/goblin-doc/example/sanfor.html)

### [Flash 钓鱼案例](https://xiecat.github.io/goblin-doc/example/flash.html)

## :question: FAQ
1. 有些网站改完会直接 302 跳转到正常页面怎么也取消不了,最终发现解决方法加 xff xfi 头
```yaml
      - Request:
          Method:
            - GET
          Header:
            x-forwarded-for: 127.0.0.1
            x-real-ip: 127.0.0.1
```

## :bar_chart: Todo 

3. websocket 支持
4. 插件系统增强（内置变量）更多匹配规则
3. 前端记录输入框输入


## :pray: 致谢


感谢`小明(Master)`的使用、反馈和建议，[\_0xf4n9x\_](https://github.com/FanqXu) 的建议。[judas](https://github.com/JonCooperWorks/judas) 带来的灵感，还有参考其他项目，才得以快速实现


## :speech_balloon: 意见交流

您可以直接在 GitHub 仓库中提交 Issue：https://github.com/xiecat/goblin/issues

与此同时您可以扫描下方群聊二维码加入我们的微信讨论群（如果群满，请稍等后续会更换二维码）：

<p align="center">
<img alt="QR-code" src="https://github.com/xiecat/goblin-doc/blob/dev/docs/.vuepress/public/wechat_group.png?raw=trueg" width="43%" height="43%" style="max-width:100%;">

## :loudspeaker: 免责声明
本工具仅能在取得足够合法授权的企业安全建设以及攻防演练中使用，在使用本工具过程中，您应确保自己所有行为符合当地的法律法规。 如您在使用本工具的过程中存在任何非法行为，您将自行承担所有后果，本工具所有开发者和所有贡献者不承担任何法律及连带责任。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。 您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。


