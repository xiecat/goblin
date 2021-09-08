package clog

import (
	"log"

	"github.com/fatih/color"
)

// consoleColors is the color set for different levels.
var consoleColors = []func(a ...interface{}) string{
	color.New(color.FgBlue).SprintFunc(),   // Trace
	color.New(color.FgGreen).SprintFunc(),  // Info
	color.New(color.FgYellow).SprintFunc(), // Warn
	color.New(color.FgRed).SprintFunc(),    // Error
	color.New(color.FgHiRed).SprintFunc(),  // Fatal
}

// ConsoleConfig is the config object for the console logger.
type ConsoleConfig struct {
	// Minimum logging level of messages to be processed.
	Level Level
}

var _ Logger = (*consoleLogger)(nil)

type consoleLogger struct {
	*noopLogger
	*log.Logger
}

func (l *consoleLogger) Write(m Messager) error {
	l.Print(consoleColors[m.Level()](m.String()))
	return nil
}

// DefaultConsoleName is the default name for the console logger.
const DefaultConsoleName = "console"

// NewConsole initializes and appends a new console logger with default name
// to the managed list.
func NewConsole(vs ...interface{}) error {
	return NewConsoleWithName(DefaultConsoleName, vs...)
}

// NewConsoleWithName initializes and appends a new console logger with given
// name to the managed list.
func NewConsoleWithName(name string, vs ...interface{}) error {
	return New(name, ConsoleIniter(), vs...)
}

// ConsoleIniter returns the initer for the console logger.
func ConsoleIniter() Initer {
	return func(name string, vs ...interface{}) (Logger, error) {
		var cfg *ConsoleConfig
		for i := range vs {
			switch v := vs[i].(type) {
			case ConsoleConfig:
				cfg = &v
			}
		}

		if cfg == nil {
			cfg = &ConsoleConfig{}
		}

		return &consoleLogger{
			noopLogger: &noopLogger{
				name:  name,
				level: cfg.Level,
			},
			Logger: log.New(color.Output, "", log.Ldate|log.Ltime),
		}, nil
	}
}
