//go:build windows || nacl || plan9
// +build windows nacl plan9

package logging

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	es6 "github.com/olivere/elastic"
	es7 "github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	els6 "gopkg.in/sohlich/elogrus.v3"
	els7 "gopkg.in/sohlich/elogrus.v7"
)

var (
	AccLogger   *logrus.Logger
	ErrorLogger *logrus.Logger
)

func (es *EsLog) Es7Setup(level logrus.Level) (log *logrus.Logger) {
	log = logrus.New()
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Open Src File err", err)
	}
	writer := bufio.NewWriter(src)
	log.SetOutput(writer)
	log.SetLevel(level)
	client, err := es7.NewClient(es7.SetURL(es.DSN))
	if err != nil {
		fmt.Printf("es conn fail please check: %s\n", err.Error())
		os.Exit(-1)
	}
	hook, err := els7.NewAsyncElasticHook(client, es.Host, es.LogLevel, es.Index)
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)
	return log
}

func (es *EsLog) Es6Setup(level logrus.Level) (log *logrus.Logger) {
	log = logrus.New()
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Open Src File err", err)
	}
	writer := bufio.NewWriter(src)
	log.SetOutput(writer)
	log.SetLevel(level)
	client, err := es6.NewClient(es6.SetURL(es.DSN))
	if err != nil {
		fmt.Printf("es conn fail please check: %s\n", err.Error())
		os.Exit(-1)
	}
	hook, err := els6.NewAsyncElasticHook(client, es.Host, es.LogLevel, es.Index)
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)
	return log
}

func (flog *FileLog) FileSetup(level logrus.Level) (log *logrus.Logger) {

	log = &logrus.Logger{
		Formatter: &prefixed.TextFormatter{
			ForceFormatting: true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
	if flog.Mode == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	log.SetLevel(level)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("logging.Setup, fail to get current dir")

	}

	file := path.Join(dir, flog.DSN)
	fileOutput, err := rotatelogs.New(file)
	if err != nil {
		log.Fatalf("logging.Setup, fail to create '%s': %v", flog.DSN, err)
	}
	log.SetOutput(fileOutput)
	return log
}
