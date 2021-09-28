# :fishing_pole_and_fish: Goblin for Phishing Exercise Tools

![GitHub branch checks state](https://img.shields.io/github/checks-status/xiecat/goblin/master) [![Latest release](https://img.shields.io/github/v/release/xiecat/goblin)](https://github.com/becivells/iconhash/releases/latest) ![GitHub Release Date](https://img.shields.io/github/release-date/xiecat/goblin) ![GitHub All Releases](https://img.shields.io/github/downloads/xiecat/goblin/total) [![GitHub issues](https://img.shields.io/github/issues/xiecat/goblin)](https://github.com/xiecat/goblin/issues) [![Docker Pulls](https://img.shields.io/docker/pulls/becivells/goblin)](https://hub.docker.com/r/becivells/goblin) ![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/becivells/goblin)       
Goblin is a phishing rehearsal tool for red-blue confrontation. By using a reverse proxy, it is possible to obtain information about a user without affecting the user's operation perceptibly, or to induce the user's operation. The purpose of hiding the server side can also be achieved by using a proxy. Built-in plug-in, through a simple configuration, quickly adjust the content of the web page to achieve a better exercise effect.

[:ledger: 中文 README](https://github.com/xiecat/goblin/blob/master/README.md)   |   [:pushpin: Releases Download](https://github.com/xiecat/goblin/releases)    |   [:book: Documents](https://goblin.xiecat.fun/)

## :collision: ​Features

* Support for caching static files to speed up access.
* Supports dumping all requests, dumping requests that match the rules.
* Support quick configuration through plug-ins to adjust inappropriate jumps or content.
* Support for implanting specific javacript code.
* Support for modifying the content of responses or goblin requests.
* Support hiding real IP by proxy.

## :tv: Demo:

![demo](https://raw.githubusercontent.com/xiecat/goblin/master/Demo.gif)

Quick Experience

1. Proxy Flash.cn

```bash
docker run -it --rm  -p 8083:8083 -p 8084:8084 -p 8085:8085 -p 8086:8086  becivells/goblin-demo-flash
```

Access to <http://127.0.0.1:8083>, corresponding example repo: [goblin-flash-demo](https://github.com/xiecat/goblin-demo/tree/master/goblin-demo-flash).

2. Proxy Baidu.com

```bash
docker run -it --rm -v $(pwd):/goblin/ -p 8084:8084 becivells/goblin
```

Access to <http://127.0.0.1:8084>.

## :computer: ​Quick Deployment

### Quick deployment with Docker

Run the following command to pull the image.

```bash
docker pull becivells/goblin
```

Dockerfile:

```yaml
FROM scratch
COPY goblin /usr/bin/goblin
ENTRYPOINT ["/usr/bin/goblin"]
WORKDIR /goblin
```

The working directory is in `goblin`, first create the directory, go to the directory and execute the following command.

```bash
docker run -it --rm -v $(pwd):/goblin/ -p 8084:8084 becivells/goblin
```

### Installing from GitHub

1. Visit [releases](https://github.com/xiecat/goblin/releases) to select the appropriate binary for your operating system from there.

2. Modify the parameters of the configuration file according to your needs. For details of the configuration file, please refer to the usage documentation 👉 [Introduction to the configuration file](https://goblin.xiecat.fun/config/).

Command-line arguments:

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

#### ⚠️ Cautions

If the server-side deployment requires changing the ip address. if you have any questions, please refer to the [`site`](https://goblin.xiecat.fun/config/site.html) explanation.

```yaml
  Site:
    server_ip:8084:  ## Change to domain name or server IP
      Listen: 0.0.0.0
      StaticPrefix: x9ut17jbqa
      SSL: false
      CAKey: ""
      CACert: ""
      ProxyPass: https://www.baidu.com
      Plugin: demo
```

## :triangular_ruler: Plugin System

See documentation for introduction details and usage 👉 [Plug-in system](https://goblin.xiecat.fun/plugin/).

## :battery: Advanced Usage

goblin uses a reverse proxy. The frontend uses a proxy such as cloudflare goblin, that can hide the goblin host. Documentation details can be found in the [goblin proxy configuration](https://goblin.xiecat.fun/guide/proxy.html).

### JS Injection

There are two ways to inject javascript: one is to follow the page (Replace requires you to append `\<script\>` tags), and the other is to follow the global js file, each approach has its own benefits.

Both of these actually use the Replace function.

#### Injecting with InjectJs

For details, please refer to [goblin InjectJs module](https://goblin.xiecat.fun/plugin/injectjs.html).

```yaml
- url: /base.js # JS to be replaced, with preference for global JS.
  Match: word   # Matching method.
  InjectJs:
    File: aaa.js # The JS to be replaced, it can be a file or a url.
```

#### Use replace to inject

For details, please refer to [goblin Replace module](https://goblin.xiecat.fun/plugin/replace.html).

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

## :star: ​Some cases

### [Flash Phishing Case](https://goblin.xiecat.fun/example/flash.html)

For more cases, please enter the Discord group for discussion, or submit an [issue](https://github.com/xiecat/goblin/issues/new/). 

## :bar_chart: Todo


- Front-end record input box input.

## :pray: Acknowledgements

Thanks to Master(小明)'s use, feedback and suggestions, and [\_0xf4n9x\_](https://github.com/FanqXu)'s suggestions. [judas](https://github.com/JonCooperWorks/judas) brought inspiration, and references to other projects, to enable quick implementation.

## :speech_balloon: Exchange of opinions

You can submit an [issue](https://github.com/xiecat/goblin/issues/new/). 

In the meantime, you can join our [Discord discussion group](https://discord.gg/BXrSruuU).

## :loudspeaker: Disclaimers

This tool can only be used in enterprise security construction and offensive and defensive exercises with sufficient legal authorization. In the process of using this tool, you should ensure that all your actions comply with local laws and regulations. If you have any illegal behavior in the process of using this tool, you will bear all the consequences by yourself, and all developers and all contributors of this tool will not bear any legal and joint liability. Please do not install and use this tool unless you have fully read, fully understood and accepted all the terms of this agreement. You are deemed to have read and agreed to be bound by this Agreement by your act of use or by your acceptance of this Agreement in any other way, express or implied.

## :laughing: Stargazers

[![Stargazers over time](https://starchart.cc/xiecat/goblin.svg)](https://starchart.cc/xiecat/goblin)

