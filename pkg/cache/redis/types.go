package redis

import (
	"fmt"
	"goblin/pkg/utils"
)

// Config redis 认证配置
type Config struct {
	Host     string `yaml:"host"`     // database address
	Port     int    `yaml:"port"`     // database port
	Password string `yaml:"password"` // database password
	DB       int    `yaml:"db"`       // database username
}

func (db *Config) ValidateDsn() error {
	if !utils.IsHost(db.Host) {
		return fmt.Errorf("the value of %s host is not in ipctl format or domainmde err", "redis")
	}
	if !utils.IsPort(db.Port) {
		return fmt.Errorf(" %s Port range of (0,65535], you set %d \n", "redis", db.Port)
	}
	return nil
}
