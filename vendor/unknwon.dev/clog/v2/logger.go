package clog

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/fatih/color"
)

// Logger is an interface for a logger with a specific name and level.
type Logger interface {
	// Name returns the name can used to identify the logger.
	Name() string
	// Level returns the minimum logging level of the logger.
	Level() Level
	// Write processes a Messager entry.
	Write(Messager) error
}

var _ Logger = (*noopLogger)(nil)

type noopLogger struct {
	name  string
	level Level
}

func (l *noopLogger) Name() string           { return l.name }
func (l *noopLogger) Level() Level           { return l.level }
func (l *noopLogger) Write(_ Messager) error { return nil }

func noopIniter(name string, _ ...interface{}) Initer {
	return func(string, ...interface{}) (Logger, error) {
		return &noopLogger{name: name}, nil
	}
}

type cancelableLogger struct {
	cancel  context.CancelFunc
	msgChan chan Messager
	done    chan struct{}
	Logger
}

var errLogger = log.New(color.Output, "", log.Ldate|log.Ltime)
var errSprintf = color.New(color.FgRed).Sprintf

func (l *cancelableLogger) error(err error) {
	if err == nil {
		return
	}

	errLogger.Print(errSprintf("[clog] [%s]: %v", l.Name(), err))
}

const (
	stateStopping int64 = iota
	stateRunning
)

type manager struct {
	state         int64
	ctx           context.Context
	cancel        context.CancelFunc
	loggers       []*cancelableLogger
	loggersByName map[string]*cancelableLogger
}

func (m *manager) len() int {
	return len(m.loggers)
}

// write attempts to send message to all loggers.
func (m *manager) write(level Level, skip int, format string, v ...interface{}) {
	if mgr.len() == 0 {
		errLogger.Print(errSprintf("[clog] no logger is available"))
		return
	}

	var msg *message
	for i := range mgr.loggers {
		if mgr.loggers[i].Level() > level {
			continue
		}

		if msg == nil {
			msg = newMessage(level, skip, format, v...)
		}

		mgr.loggers[i].msgChan <- msg
	}
}

// writeTo attempts to send message to the logger with given name.
func (m *manager) writeTo(name string, level Level, skip int, format string, v ...interface{}) {
	l, ok := mgr.loggersByName[name]
	if !ok {
		errLogger.Print(errSprintf("[clog] logger with name %q is not available", name))
		return
	}

	if l.Level() > level {
		return
	}

	l.msgChan <- newMessage(level, skip, format, v...)
}

func (m *manager) stop() {
	// Make sure cancellation is only propagated once to prevent deadlock of WaitForStop.
	if !atomic.CompareAndSwapInt64(&m.state, stateRunning, stateStopping) {
		return
	}

	m.cancel()
	for _, l := range m.loggers {
		<-l.done
	}
}

var mgr *manager

func init() {
	ctx, cancel := context.WithCancel(context.Background())
	mgr = &manager{
		state:         stateRunning,
		ctx:           ctx,
		cancel:        cancel,
		loggersByName: make(map[string]*cancelableLogger),
	}
}

// Initer takes a name and arbitrary number of parameters needed for initalization
// and returns an initalized logger.
type Initer func(string, ...interface{}) (Logger, error)

// New initializes and appends a new logger to the managed list.
// Calling this function multiple times will overwrite previous initialized
// logger with the same name.
//
// Any integer type (i.e. int, int32, int64) will be used as buffer size.
// Otherwise, the value will be passed to the initer.
//
// NOTE: This function is not concurrent safe.
func New(name string, initer Initer, opts ...interface{}) error {
	bufferSize := 0

	vs := opts[:0]
	for i := range opts {
		switch opt := opts[i].(type) {
		case int:
			bufferSize = opt
		case int32:
			bufferSize = int(opt)
		case int64:
			bufferSize = int(opt)
		default:
			vs = append(vs, opt)
		}
	}

	l, err := initer(name, vs...)
	if err != nil {
		return fmt.Errorf("initialize logger: %v", err)
	}

	if bufferSize < 0 {
		bufferSize = 0
	}

	ctx, cancel := context.WithCancel(mgr.ctx)
	cl := &cancelableLogger{
		cancel:  cancel,
		msgChan: make(chan Messager, bufferSize),
		done:    make(chan struct{}),
		Logger:  l,
	}

	// Check and replace previous logger
	found := false
	for i, l := range mgr.loggers {
		if l.Name() == name {
			found = true

			// Release previous logger
			l.cancel()
			<-l.done

			mgr.loggers[i] = cl
			break
		}
	}
	if !found {
		mgr.loggers = append(mgr.loggers, cl)
	}
	mgr.loggersByName[name] = cl

	go func() {
	loop:
		for {
			select {
			case m := <-cl.msgChan:
				cl.error(cl.Write(m))
			case <-ctx.Done():
				break loop
			}
		}

		// Drain the msgChan at best effort
		for {
			if len(cl.msgChan) == 0 {
				break
			}

			cl.error(cl.Write(<-cl.msgChan))
		}

		// Notify the cleanup is done
		cl.done <- struct{}{}
	}()
	return nil
}

// Remove removes a logger with given name from the managed list.
//
// NOTE: This function is not concurrent safe.
func Remove(name string) {
	loggers := mgr.loggers[:0]
	for _, l := range mgr.loggers {
		if l.Name() == name {
			go func(l *cancelableLogger) {
				l.cancel()
				<-l.done
			}(l)
			continue
		}
		loggers = append(loggers, l)
	}
	mgr.loggers = loggers
	delete(mgr.loggersByName, name)
}
