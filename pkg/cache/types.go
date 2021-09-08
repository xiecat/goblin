package cache

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"goblin/pkg/cache/nocache"
	"goblin/pkg/utils"
	"strings"
	"time"
	log "unknwon.dev/clog/v2"

	cacheNew "github.com/patrickmn/go-cache"
	newRedis "goblin/pkg/cache/redis"
)

var // 支持的缓存类型
cacheType = []string{"self", "redis", "none"}

type Config struct {
	Type    string           `yaml:"type"`
	ExpTime time.Duration    `yaml:"expire_time"`
	Redis   *newRedis.Config `yaml:"redis"`
}

type Cache struct {
	Type       string
	ExpireTime time.Duration
	Self       *cacheNew.Cache
	Redis      *redis.Client
	NoCache    *nocache.Cache
}

// ValidateCacheDsn 验证缓存的配置信息
func (db *Config) ValidateCacheDsn() error {
	if !utils.StrEqualOrInList(db.Type, cacheType) {
		return fmt.Errorf("cache value is %s type must %s\n", db.Type, strings.Join(cacheType, ","))
	}
	if db.Type == "redis" {
		return db.Redis.ValidateDsn()
	}
	return nil
}

func (config *Config) ValidateCachePing() error {
	if config.Type == "redis" {
		_, err := GlobalCache.Redis.Ping(ctx).Result()

		if err != nil {
			return fmt.Errorf("redis conn error %s", err.Error())
		}
		log.Info("use redis cache conn ok")
	} else if config.Type == "self" {
		log.Info("use self cache conn ok")
	} else {
		log.Info("no use cache")
	}
	return nil
}
