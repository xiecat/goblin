package clog

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	simpleDateFormat = "2006-01-02"
	logPrefixLength  = len("2017/02/06 21:20:08 ")
)

// FileRotationConfig represents rotation related configurations for file mode logger.
// All the settings can take effect at the same time, remain zero values to disable them.
type FileRotationConfig struct {
	// Do rotation for output files.
	Rotate bool
	// Rotate on daily basis.
	Daily bool
	// Maximum size in bytes of file for a rotation.
	MaxSize int64
	// Maximum number of lines for a rotation.
	MaxLines int64
	// Maximum lifetime of a output file in days.
	MaxDays int64
}

// FileConfig is the config object for the file logger.
type FileConfig struct {
	// Minimum level of messages to be processed.
	Level Level
	// File name to output messages.
	Filename string
	// Rotation related configurations.
	FileRotationConfig
}

var _ Logger = (*fileLogger)(nil)

type fileLogger struct {
	// Indicates whether it is being used as standalone logger.
	// It is only true when the logger is created by NewFileWriter.
	standalone bool

	*noopLogger

	filename       string
	rotationConfig FileRotationConfig

	// Rotation metadata
	file         *os.File
	openDay      int
	currentSize  int64
	currentLines int64

	*log.Logger
}

var newLineBytes = []byte("\n")

func (l *fileLogger) initFile() (err error) {
	l.file, err = os.OpenFile(l.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("open file %q: %v", l.filename, err)
	}

	l.Logger = log.New(l.file, "", log.Ldate|log.Ltime)
	return nil
}

// isExist returns true if the file or directory exists.
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// rotateFilename returns next available rotate filename in given date.
func rotateFilename(filename, date string) string {
	filename = fmt.Sprintf("%s.%s", filename, date)
	if !isExist(filename) {
		return filename
	}

	format := filename + ".%03d"
	for i := 1; i < 1000; i++ {
		filename := fmt.Sprintf(format, i)
		if !isExist(filename) {
			return filename
		}
	}

	panic("too many log files for yesterday, already reached 999")
}

func (l *fileLogger) deleteOutdatedFiles() error {
	return filepath.Walk(filepath.Dir(l.filename), func(path string, fi os.FileInfo, _ error) error {
		if !fi.IsDir() &&
			fi.ModTime().Before(time.Now().Add(-24*time.Hour*time.Duration(l.rotationConfig.MaxDays))) &&
			strings.HasPrefix(filepath.Base(path), filepath.Base(l.filename)) {
			return os.Remove(path)
		}
		return nil
	})
}

func (l *fileLogger) initRotation() error {
	// Gather basic file info for rotation.
	fi, err := l.file.Stat()
	if err != nil {
		return fmt.Errorf("stat: %v", err)
	}

	l.currentSize = fi.Size()

	// If there is any content in the file, count the number of lines.
	if l.rotationConfig.MaxLines > 0 && l.currentSize > 0 {
		data, err := ioutil.ReadFile(l.filename)
		if err != nil {
			return fmt.Errorf("read file %q: %v", l.filename, err)
		}

		l.currentLines = int64(bytes.Count(data, newLineBytes)) + 1
	}

	if l.rotationConfig.Daily {
		now := time.Now()
		l.openDay = now.Day()

		lastWriteTime := fi.ModTime()
		if lastWriteTime.Year() != now.Year() ||
			lastWriteTime.Month() != now.Month() ||
			lastWriteTime.Day() != now.Day() {

			if err = l.file.Close(); err != nil {
				return fmt.Errorf("close current file: %v", err)
			}
			if err = os.Rename(l.filename, rotateFilename(l.filename, lastWriteTime.Format(simpleDateFormat))); err != nil {
				return fmt.Errorf("rename rotate file: %v", err)
			}

			if err = l.initFile(); err != nil {
				return fmt.Errorf("init file: %v", err)
			}
		}
	}

	if l.rotationConfig.MaxDays > 0 {
		if err = l.deleteOutdatedFiles(); err != nil {
			return fmt.Errorf("delete outdated files: %v", err)
		}
	}
	return nil
}

func (l *fileLogger) write(m Messager) (int, error) {
	l.Logger.Print(m.String())

	bytesWrote := len(m.String())
	if !l.standalone {
		bytesWrote += logPrefixLength
	}
	if l.rotationConfig.Rotate {
		l.currentSize += int64(bytesWrote)
		l.currentLines += int64(strings.Count(m.String(), "\n")) + 1

		var (
			needsRotate = false
			rotateDate  time.Time
		)

		now := time.Now()
		if l.rotationConfig.Daily && now.Day() != l.openDay {
			needsRotate = true
			rotateDate = now.Add(-24 * time.Hour)

		} else if (l.rotationConfig.MaxSize > 0 && l.currentSize >= l.rotationConfig.MaxSize) ||
			(l.rotationConfig.MaxLines > 0 && l.currentLines >= l.rotationConfig.MaxLines) {
			needsRotate = true
			rotateDate = now
		}

		if needsRotate {
			_ = l.file.Close()
			if err := os.Rename(l.filename, rotateFilename(l.filename, rotateDate.Format(simpleDateFormat))); err != nil {
				return bytesWrote, fmt.Errorf("rename rotated file %q: %v", l.filename, err)
			}

			if err := l.initFile(); err != nil {
				return bytesWrote, fmt.Errorf("init file %q: %v", l.filename, err)
			}

			l.openDay = now.Day()
			l.currentSize = 0
			l.currentLines = 0

			if err := l.deleteOutdatedFiles(); err != nil {
				return bytesWrote, fmt.Errorf("delete outdated file: %v", err)
			}
		}
	}
	return bytesWrote, nil
}

func (l *fileLogger) Write(m Messager) error {
	_, err := l.write(m)
	return err
}

func (l *fileLogger) init() error {
	_ = os.MkdirAll(filepath.Dir(l.filename), os.ModePerm)
	if err := l.initFile(); err != nil {
		return fmt.Errorf("init file %q: %v", l.filename, err)
	}

	if l.rotationConfig.Rotate {
		if err := l.initRotation(); err != nil {
			return fmt.Errorf("init rotation: %v", err)
		}
	}
	return nil
}

// DefaultFileName is the default name for the file logger.
const DefaultFileName = "file"

// NewFile initializes and appends a new file logger with default name
// to the managed list.
func NewFile(vs ...interface{}) error {
	return NewFileWithName(DefaultFileName, vs...)
}

// NewFileWithName initializes and appends a new file logger with given
// name to the managed list.
func NewFileWithName(name string, vs ...interface{}) error {
	return New(name, FileIniter(), vs...)
}

// FileIniter returns the initer for the file logger.
func FileIniter() Initer {
	return func(name string, vs ...interface{}) (Logger, error) {
		var cfg *FileConfig
		for i := range vs {
			switch v := vs[i].(type) {
			case FileConfig:
				cfg = &v
			}
		}

		if cfg == nil {
			cfg = &FileConfig{
				Filename: "clog.log",
			}
		}

		l := &fileLogger{
			noopLogger: &noopLogger{
				name:  name,
				level: cfg.Level,
			},
			filename:       cfg.Filename,
			rotationConfig: cfg.FileRotationConfig,
		}

		if err := l.init(); err != nil {
			return nil, err
		}

		return l, nil
	}
}

var _ io.Writer = (*fileWriter)(nil)

type fileWriter struct {
	*fileLogger
}

// NewFileWriter returns an io.Writer for synchronized file logger.
func NewFileWriter(filename string, cfg FileRotationConfig) (io.Writer, error) {
	f := &fileLogger{
		standalone:     true,
		filename:       filename,
		rotationConfig: cfg,
	}
	if err := f.init(); err != nil {
		return nil, fmt.Errorf("init: %v", err)
	}

	return &fileWriter{f}, nil
}

// Write implements method of io.Writer interface.
func (w *fileWriter) Write(p []byte) (int, error) {
	return w.write(&message{
		body: string(p),
	})
}
