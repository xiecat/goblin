# ElasticSearch Hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>
<img src="https://travis-ci.org/sohlich/elogrus.svg?branch=master" />

Elasticsearch version | Elastic version | Package URL
----------------------|------------------|------------
6.x                   | 6.0              | [`gopkg.in/sohlich/elogrus.v3`](http://gopkg.in/sohlich/elogrus.v3)
5.x                   | 5.0              | [`gopkg.in/sohlich/elogrus.v2`](http://gopkg.in/sohlich/elogrus.v2)
2.x                   | 3.0              | [`gopkg.in/sohlich/elogrus.v1`](http://gopkg.in/sohlich/elogrus.v1)


## Changelog
- elastic 6.x support (currently in master)
- v2.1 - Added support for async hook


## Import

```
go get gopkg.in/sohlich/elogrus.v3
```

## Usage

```go
package main

import (
	"github.com/sirupsen/logrus.v3"
	"gopkg.in/sohlich/elogrus"
	"gopkg.in/olivere/elastic"
)


func main() {
	log := logrus.New()
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Panic(err)
	}	
	hook, err := elogrus.NewElasticHook(client, "localhost", logrus.DebugLevel, "mylog")
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