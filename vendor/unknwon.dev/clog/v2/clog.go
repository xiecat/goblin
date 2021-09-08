// Package clog is a channel-based logging package for Go.
package clog

import (
	"fmt"
	"os"
)

// Level is the logging level.
type Level int

// Available logging levels.
const (
	LevelTrace Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		fmt.Printf("Unexpected Level value: %v\n", int(l))
		panic("unreachable")
	}
}

// Trace writes formatted log in Trace level.
func Trace(format string, v ...interface{}) {
	mgr.write(LevelTrace, 0, format, v...)
}

// Info writes formatted log in Info level.
func Info(format string, v ...interface{}) {
	mgr.write(LevelInfo, 0, format, v...)
}

// Warn writes formatted log in Warn level.
func Warn(format string, v ...interface{}) {
	mgr.write(LevelWarn, 0, format, v...)
}

// Error writes formatted log in Error level.
func Error(format string, v ...interface{}) {
	ErrorDepth(4, format, v...)
}

// ErrorDepth writes formatted log with given skip depth in Error level.
func ErrorDepth(skip int, format string, v ...interface{}) {
	mgr.write(LevelError, skip, format, v...)
}

// Fatal writes formatted log in Fatal level then exits.
func Fatal(format string, v ...interface{}) {
	FatalDepth(4, format, v...)
}

// isTestEnv is true when running tests.
// In test environment, Fatal or FatalDepth won't stop the manager or exit the program.
var isTestEnv = false

func exit() {
	if isTestEnv {
		return
	}

	Stop()
	os.Exit(1)
}

// FatalDepth writes formatted log with given skip depth in Fatal level then exits.
func FatalDepth(skip int, format string, v ...interface{}) {
	mgr.write(LevelFatal, skip, format, v...)
	exit()
}

// TraceTo writes formatted log in Trace level to the logger with given name.
func TraceTo(name, format string, v ...interface{}) {
	mgr.writeTo(name, LevelTrace, 0, format, v...)
}

// InfoTo writes formatted log in Info level to the logger with given name.
func InfoTo(name, format string, v ...interface{}) {
	mgr.writeTo(name, LevelInfo, 0, format, v...)
}

// WarnTo writes formatted log in Warn level to the logger with given name.
func WarnTo(name, format string, v ...interface{}) {
	mgr.writeTo(name, LevelWarn, 0, format, v...)
}

// ErrorTo writes formatted log in Error level to the logger with given name.
func ErrorTo(name, format string, v ...interface{}) {
	ErrorDepthTo(name, 4, format, v...)
}

// ErrorDepthTo writes formatted log with given skip depth in Error level to
// the logger with given name.
func ErrorDepthTo(name string, skip int, format string, v ...interface{}) {
	mgr.writeTo(name, LevelError, skip, format, v...)
}

// FatalTo writes formatted log in Fatal level to the logger with given name
// then exits.
func FatalTo(name, format string, v ...interface{}) {
	FatalDepthTo(name, 4, format, v...)
}

// FatalDepthTo writes formatted log with given skip depth in Fatal level to
// the logger with given name then exits.
func FatalDepthTo(name string, skip int, format string, v ...interface{}) {
	mgr.writeTo(name, LevelFatal, skip, format, v...)
	exit()
}

// Stop propagates cancellation to all loggers and waits for completion.
// This function should always be called before exiting the program.
func Stop() {
	mgr.stop()
}
