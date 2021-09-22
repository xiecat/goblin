package elogrus

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

var (
	// ErrCannotCreateIndex Fired if the index is not created
	ErrCannotCreateIndex = fmt.Errorf("cannot create index")
)

// IndexNameFunc get index name
type IndexNameFunc func() string

type fireFunc func(entry *logrus.Entry, hook *ElasticHook) error

// ElasticHook is a logrus
// hook for ElasticSearch
type ElasticHook struct {
	client    *elastic.Client
	host      string
	index     IndexNameFunc
	levels    []logrus.Level
	ctx       context.Context
	ctxCancel context.CancelFunc
	fireFunc  fireFunc
}

type message struct {
	Host      string
	Timestamp string `json:"@timestamp"`
	Message   string
	Data      logrus.Fields
	Level     string
}

// NewElasticHook creates new hook.
// client - ElasticSearch client with specific es version (v5/v6/v7/...)
// host - host of system
// level - log level
// index - name of the index in ElasticSearch
func NewElasticHook(client *elastic.Client, host string, level logrus.Level, index string) (*ElasticHook, error) {
	return NewElasticHookWithFunc(client, host, level, func() string { return index })
}

// NewAsyncElasticHook creates new  hook with asynchronous log.
// client - ElasticSearch client with specific es version (v5/v6/v7/...)
// host - host of system
// level - log level
// index - name of the index in ElasticSearch
func NewAsyncElasticHook(client *elastic.Client, host string, level logrus.Level, index string) (*ElasticHook, error) {
	return NewAsyncElasticHookWithFunc(client, host, level, func() string { return index })
}

// NewBulkProcessorElasticHook creates new hook that uses a bulk processor for indexing.
// client - ElasticSearch client with specific es version (v5/v6/v7/...)
// host - host of system
// level - log level
// index - name of the index in ElasticSearch
func NewBulkProcessorElasticHook(client *elastic.Client, host string, level logrus.Level, index string) (*ElasticHook, error) {
	return NewBulkProcessorElasticHookWithFunc(client, host, level, func() string { return index })
}

// NewElasticHookWithFunc creates new hook with
// function that provides the index name. This is useful if the index name is
// somehow dynamic especially based on time.
// client - ElasticSearch client with specific es version (v5/v6/v7/...)
// host - host of system
// level - log level
// indexFunc - function providing the name of index
func NewElasticHookWithFunc(client *elastic.Client, host string, level logrus.Level, indexFunc IndexNameFunc) (*ElasticHook, error) {
	return newHookFuncAndFireFunc(client, host, level, indexFunc, syncFireFunc)
}

// NewAsyncElasticHookWithFunc creates new asynchronous hook with
// function that provides the index name. This is useful if the index name is
// somehow dynamic especially based on time.
// client - ElasticSearch client with specific es version (v5/v6/v7/...)
// host - host of system
// level - log level
// indexFunc - function providing the name of index
func NewAsyncElasticHookWithFunc(client *elastic.Client, host string, level logrus.Level, indexFunc IndexNameFunc) (*ElasticHook, error) {
	return newHookFuncAndFireFunc(client, host, level, indexFunc, asyncFireFunc)
}

// NewBulkProcessorElasticHookWithFunc creates new hook with
// function that provides the index name. This is useful if the index name is
// somehow dynamic especially based on time that uses a bulk processor for
// indexing.
// client - ElasticSearch client with specific es version (v5/v6/v7/...)
// host - host of system
// level - log level
// indexFunc - function providing the name of index
func NewBulkProcessorElasticHookWithFunc(client *elastic.Client, host string, level logrus.Level, indexFunc IndexNameFunc) (*ElasticHook, error) {
	fireFunc, err := makeBulkFireFunc(client)
	if err != nil {
		return nil, err
	}
	return newHookFuncAndFireFunc(client, host, level, indexFunc, fireFunc)
}

func newHookFuncAndFireFunc(client *elastic.Client, host string, level logrus.Level, indexFunc IndexNameFunc, fireFunc fireFunc) (*ElasticHook, error) {
	var levels []logrus.Level
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	} {
		if l <= level {
			levels = append(levels, l)
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(indexFunc()).Do(ctx)
	if err != nil {
		// Handle error
		cancel()
		return nil, err
	}
	if !exists {
		createIndex, err := client.CreateIndex(indexFunc()).Do(ctx)
		if err != nil {
			cancel()
			return nil, err
		}
		if !createIndex.Acknowledged {
			cancel()
			return nil, ErrCannotCreateIndex
		}
	}

	return &ElasticHook{
		client:    client,
		host:      host,
		index:     indexFunc,
		levels:    levels,
		ctx:       ctx,
		ctxCancel: cancel,
		fireFunc:  fireFunc,
	}, nil
}

// Fire is required to implement
// Logrus hook
func (hook *ElasticHook) Fire(entry *logrus.Entry) error {
	return hook.fireFunc(entry, hook)
}

func asyncFireFunc(entry *logrus.Entry, hook *ElasticHook) error {
	go syncFireFunc(entry, hook)
	return nil
}

func createMessage(entry *logrus.Entry, hook *ElasticHook) *message {
	level := entry.Level.String()

	if e, ok := entry.Data[logrus.ErrorKey]; ok && e != nil {
		if err, ok := e.(error); ok {
			entry.Data[logrus.ErrorKey] = err.Error()
		}
	}

	return &message{
		hook.host,
		entry.Time.UTC().Format(time.RFC3339Nano),
		entry.Message,
		entry.Data,
		strings.ToUpper(level),
	}
}

func syncFireFunc(entry *logrus.Entry, hook *ElasticHook) error {
	_, err := hook.client.
		Index().
		Index(hook.index()).
		Type("log").
		BodyJson(*createMessage(entry, hook)).
		Do(hook.ctx)

	return err
}

// Create closure with bulk processor tied to fireFunc.
func makeBulkFireFunc(client *elastic.Client) (fireFunc, error) {
	processor, err := client.BulkProcessor().
		Name("elogrus.v3.bulk.processor").
		Workers(2).
		FlushInterval(time.Second).
		Do(context.Background())

	return func(entry *logrus.Entry, hook *ElasticHook) error {
		r := elastic.NewBulkIndexRequest().
			Index(hook.index()).
			Type("log").
			Doc(*createMessage(entry, hook))
		processor.Add(r)
		return nil
	}, err
}

// Levels Required for logrus hook implementation
func (hook *ElasticHook) Levels() []logrus.Level {
	return hook.levels
}

// Cancel all calls to elastic
func (hook *ElasticHook) Cancel() {
	hook.ctxCancel()
}
