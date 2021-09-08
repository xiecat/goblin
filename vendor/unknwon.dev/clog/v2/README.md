# Clog 

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/go-clog/clog/Go?logo=github&style=for-the-badge)](https://github.com/go-clog/clog/actions?query=workflow%3AGo)
[![codecov](https://img.shields.io/codecov/c/github/go-clog/clog/master?logo=codecov&style=for-the-badge)](https://codecov.io/gh/go-clog/clog)
[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/unknwon.dev/clog/v2?tab=doc)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/go-clog/clog)

![](https://avatars1.githubusercontent.com/u/25576866?v=3&s=200)

Package clog is a channel-based logging package for Go.

This package supports multiple loggers across different levels of logging. It uses Go's native channel feature to provide goroutine-safe mechanism on large concurrency.

## Installation

The minimum requirement of Go is **1.11**.

	go get unknwon.dev/clog/v2
    
Please apply `-u` flag to update in the future.

## Getting Started

It is extremely easy to create one with all default settings. Generally, you would want to create new logger inside `init` or `main` function.

Let's create a logger that prints logs to the console:

```go
import (
	log "unknwon.dev/clog/v2"
)

func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}

func main() {
	log.Trace("Hello %s!", "World") // YYYY/MM/DD 12:34:56 [TRACE] Hello World!
	log.Info("Hello %s!", "World")  // YYYY/MM/DD 12:34:56 [ INFO] Hello World!
	log.Warn("Hello %s!", "World")  // YYYY/MM/DD 12:34:56 [ WARN] Hello World!

	// Graceful stopping all loggers before exiting the program.
	log.Stop()
}
```

The code inside `init` function is equivalent to the following:

```go
func init() {
	err := log.NewConsole(0, 
        log.ConsoleConfig{
		    Level: log.LevelTrace,
	    },
    )
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}
```

Or expand further:

```go
func init() {
	err := log.NewConsoleWithName(log.DefaultConsoleName, 0, 
        log.ConsoleConfig{
		    Level: log.LevelTrace,
	    },
    )
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}
```

- The `0` is an integer type so it is used as underlying buffer size. In this case, `0` creates synchronized logger (call hangs until write is finished).
- Any non-integer type is used as the config object, in this case `ConsoleConfig` is the respective config object for the console logger.
- The `LevelTrace` used here is the lowest logging level, meaning prints every log to the console. All levels from lowest to highest are: `LevelTrace`, `LevelInfo`, `LevelWarn`, `LevelError`, `LevelFatal`, each of them has at least one respective function, e.g. `log.Trace`, `log.Info`, `log.Warn`, `log.Error` and `log.Fatal`.

In production, you may want to make log less verbose and be asynchronous:

```go
func init() {
	// The buffer size mainly depends on number of logs could be produced at the same time, 
	// 100 is a good default.
	err := log.NewConsole(100,
        log.ConsoleConfig{
		    Level:      log.LevelInfo,
	    },
    )
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}
```

- When you set level to be `LevelInfo`, calls to the `log.Trace` will be simply noop.
- The console logger comes with color output, but for non-colorable destination, the color output will be disabled automatically.

Other builtin loggers are file (`log.NewFile`), Slack (`log.NewSlack`) and Discord (`log.NewDiscord`), see later sections in the documentation for usage details.

### Multiple Loggers

You can have multiple loggers in different modes across levels.

```go
func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
	err := log.NewFile(
        log.FileConfig{
		    Level:    log.LevelInfo,
		    Filename: "clog.log",
	    },
    )
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}
```

In this example, all logs will be printed to console, and only logs with level Info or higher (i.e. Warn, Error and Fatal) will be written into file.

### Write to a specific logger

When multiple loggers are registered, it is also possible to write logs to a special logger by giving its name.

```go
func main() {
	log.TraceTo(log.DefaultConsoleName, "Hello %s!", "World")
	log.InfoTo(log.DefaultConsoleName, "Hello %s!", "World")
	log.WarnTo(log.DefaultConsoleName, "Hello %s!", "World")
	log.ErrorTo(log.DefaultConsoleName, "So bad... %v", err)
	log.FatalTo(log.DefaultConsoleName, "Boom! %v", err)

	// ...
}
```

### Caller Location

When using `log.Error` and `log.Fatal` functions, the caller location is written along with logs. 

```go
func main() {
	log.Error("So bad... %v", err) // YYYY/MM/DD 12:34:56 [ERROR] [...er/main.go:64 main()] ...
	log.Fatal("Boom! %v", err)     // YYYY/MM/DD 12:34:56 [FATAL] [...er/main.go:64 main()] ...

	// ...
}
```

- Calling `log.Fatal` will exit the program.
- If you want to have different skip depth than the default, use `log.ErrorDepth` or `log.FatalDepth`.

### Clean Exit

You should always call `log.Stop()` to wait until all logs are processed before program exits.

## Builtin Loggers

### File Logger

File logger is the single most powerful builtin logger, it has the ability to rotate based on file size, line, and date:

```go
func init() {
	err := log.NewFile(100, 
        log.FileConfig{
            Level:              log.LevelInfo,
            Filename:           "clog.log",  
            FileRotationConfig: log.FileRotationConfig {
                Rotate: true,
                Daily:  true,
            },
        },
    )
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}
```

In case you have some other packages that write to a file, and you want to take advatange of this file rotation feature. You can do so by using the `log.NewFileWriter` function. It acts like a standard `io.Writer`.

```go
func init() {
	w, err := log.NewFileWriter("filename",
        log.FileRotationConfig{
            Rotate: true,
            Daily:  true,
        },
    )
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}
```

### Slack Logger

Slack logger is also supported in a simple way:

```go
func init() {
	err := log.NewSlack(100,
        log.SlackConfig{
            Level: log.LevelInfo,
            URL:   "https://url-to-slack-webhook",
        },
    )
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}
```

This logger also works for [Discord Slack](https://discordapp.com/developers/docs/resources/webhook#execute-slackcompatible-webhook) endpoint.

### Discord Logger

Discord logger is supported in rich format via [Embed Object](https://discordapp.com/developers/docs/resources/channel#embed-object):

```go
func init() {
	err := log.NewDiscord(100,
        log.DiscordConfig{
            Level: log.LevelInfo,
            URL:   "https://url-to-discord-webhook",
        },
    )
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}
```

This logger automatically retries up to 3 times if hits rate limit with respect to `retry_after`.

## Build Your Own Logger

You can implement your own logger and all the concurrency stuff are handled automatically!

Here is an example which sends all logs to a channel, we call it `chanLogger` here:

```go
import log "unknwon.dev/clog/v2"

type chanConfig struct {
	c chan string
}

var _ log.Logger = (*chanLogger)(nil)

type chanLogger struct {
	name  string
	level log.Level
	c     chan string
}

func (l *chanLogger) Name() string     { return l.name }
func (l *chanLogger) Level() log.Level { return l.level }

func (l *chanLogger) Write(m log.Messager) error {
	l.c <- m.String()
	return nil
}

func main() {
	log.New("channel", func(name string, vs ...interface{}) (log.Logger, error) {
		var cfg *chanConfig
		for i := range vs {
			switch v := vs[i].(type) {
			case chanConfig:
				cfg = &v
			}
		}

		if cfg == nil {
			return nil, fmt.Errorf("config object with the type '%T' not found", chanConfig{})
		} else if cfg.c == nil {
			return nil, errors.New("channel is nil")
		}

		return &chanLogger{
			name: name,
			c:    cfg.c,
		}, nil
	})
}
```

Have fun!

## Credits

- Avatar is a modified version based on [egonelbre/gophers' scientist](https://github.com/egonelbre/gophers/blob/master/vector/science/scientist.svg).

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
