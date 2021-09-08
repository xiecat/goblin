package clog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

var _ Messager = (*message)(nil)

// Messager is a message entry to be processed by logger.
type Messager interface {
	// Level returns the level of the message.
	Level() Level
	fmt.Stringer
}

type message struct {
	level Level
	body  string
}

func newMessage(level Level, skip int, format string, v ...interface{}) *message {
	var body string
	// Only error and fatal information needs locate position for debugging.
	// But if skip is 0 means caller doesn't care so we can skip.
	if level >= LevelError && skip > 0 {
		pc, file, line, ok := runtime.Caller(skip)
		if ok {
			// Get caller function name
			fn := runtime.FuncForPC(pc)
			var fnName string
			if fn == nil {
				fnName = "?()"
			} else {
				fnName = strings.TrimLeft(filepath.Ext(fn.Name()), ".") + "()"
			}

			if len(file) > 32 {
				file = "..." + file[len(file)-32:]
			}
			body = fmt.Sprintf("[%s:%d %s] %s", file, line, fnName, fmt.Sprintf(format, v...))
		}
	}
	if len(body) == 0 {
		body = fmt.Sprintf(format, v...)
	}
	return &message{
		level: level,
		body:  fmt.Sprintf("[%5s] %s", level, body),
	}
}

func (m *message) Level() Level   { return m.level }
func (m *message) String() string { return m.body }
