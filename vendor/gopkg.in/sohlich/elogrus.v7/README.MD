# ElasticSearch Hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>
<img src="https://travis-ci.org/sohlich/elogrus.svg?branch=master" />

## Releases
**Notice that the master branch always refers to the latest version of Elastic. If you want to use stable versions of elogus, you should use the packages released via [gopkg.in](https://gopkg.in).**

*Here's the version matrix:*

ElasticSearch version | Elastic version  | Package URL                              | Remarks |
----------------------|------------------|------------------------------------------|---------|
7.x                   | 7.0              | [`gopkg.in/sohlich/elogrus.v7`](http://gopkg.in/sohlich/elogrus.v7)| Use Go modules.
6.x                   | 6.0              | [`gopkg.in/sohlich/elogrus.v3`](http://gopkg.in/sohlich/elogrus.v3)| Actively maintained.
5.x                   | 5.0              | [`gopkg.in/sohlich/elogrus.v2`](http://gopkg.in/sohlich/elogrus.v2)| Actively maintained.
2.x                   | 3.0              | [`gopkg.in/sohlich/elogrus.v1`](http://gopkg.in/sohlich/elogrus.v1)| Deprecated. Please update.

*For elastic search 7.x*
```bash
# We name v7 to align with elastic v7
go get gopkg.in/sohlich/elogrus.v7
go get github.com/olivere/elastic/v7
```

*For elastic search 6.x*
```bash
go get gopkg.in/sohlich/elogrus.v3
go get github.com/olivere/elastic

```

*For elastic search 5.x*
```bash
go get gopkg.in/sohlich/elogrus.v2
go get gopkg.in/olivere/elastic.v5

```

## Changelog
- elastic 7.x support (currently in master)

## Usage

```go
package main

import (
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
)

func main() {
	log := logrus.New()
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Panic(err)
	}
	hook, err := elogrus.NewAsyncElasticHook(client, "localhost", logrus.DebugLevel, "mylog")
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)

	log.WithFields(logrus.Fields{
		"name": "joe",
		"age":  42,
	}).Error("Hello world!")
}
```

### Asynchronous hook

```go
	...
	elogrus.NewAsyncElasticHook(client, "localhost", logrus.DebugLevel, "mylog")
	...
```
