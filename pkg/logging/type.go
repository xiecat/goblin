package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"goblin/pkg/utils"
	"strings"
)

var outType []string = []string{"es7", "es6", "file", "syslog"}
var modeType = []string{"json", "text"}

type Config struct {
	Type     string
	LogLevel logrus.Level
	EsLog    *EsLog
	FileLog  *FileLog
	Syslog   *Syslog
}

type EsLog struct {
	LogLevel logrus.Level
	DSN      string
	Index    string // name
	Host     string
}
type FileLog struct {
	Mode string
	DSN  string
}
type Syslog struct {
	Mode string
	DSN  string
}

func (conf *Config) New() (log *logrus.Logger) {
	switch conf.Type {
	case "es7":
		return conf.EsLog.Es7Setup(conf.LogLevel)
	case "es6":
		return conf.EsLog.Es6Setup(conf.LogLevel)
	case "file":
		return conf.FileLog.FileSetup(conf.LogLevel)
	case "syslog":
		return conf.Syslog.SyslogSetup(conf.LogLevel)

	}
	return
}

func (cf *Config) ValidateType() error {
	if !utils.StrEqualOrInList(cf.Type, outType) {
		return fmt.Errorf("log  value %s type must %s ", cf.Type, strings.Join(outType, ","))
	}
	switch cf.Type {
	case "file":
		if !utils.StrEqualOrInList(cf.FileLog.Mode, modeType) {
			return fmt.Errorf("file log  value %s type must %s ", cf.Type, strings.Join(modeType, ","))
		}
	case "syslog":
		if !utils.StrEqualOrInList(cf.Syslog.Mode, modeType) {
			return fmt.Errorf("syslog  value %s type must %s ", cf.Type, strings.Join(modeType, ","))
		}

	}
	return nil
}
