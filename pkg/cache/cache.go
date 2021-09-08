package cache

import (
	"context"
	"fmt"
	"time"

	log "unknwon.dev/clog/v2"

	"goblin/pkg/cache/cache"
	"goblin/pkg/cache/nocache"
	"goblin/pkg/cache/redis"
)

var ctx = context.Background()

var GlobalCache *Cache

func (db *Config) Init() {
	c := &Cache{
		ExpireTime: db.ExpTime,
		Type:       db.Type,
	}
	switch db.Type {
	case "self":
		c.Self = cache.Init(db.ExpTime)
	case "redis":
		c.Redis = redis.Init(db.Redis)
	case "none":
		c.NoCache = nocache.Init()
	default:
		log.Fatal("unsupported database type: %s", db.Type)
	}
	GlobalCache = c
}

func (cache *Cache) Set(key string, v interface{}) {
	switch cache.Type {
	case "self":
		cache.Self.Set(key, v, time.Duration(-1))
	case "redis":
		cache.Redis.Set(ctx, key, v, 0)
	case "none":
		cache.NoCache.Set(key, v, time.Duration(-1))
	}
}

// SetNX 带过期时间的
func (cache *Cache) SetNX(key string, v interface{}) {
	expTime := cache.ExpireTime * time.Minute
	switch cache.Type {
	case "self":
		cache.Self.Set(key, v, expTime)
	case "redis":
		cache.Redis.Set(ctx, key, v, expTime)
	case "none":
		cache.NoCache.Set(key, v, time.Duration(-1))
	}
}

func (cache *Cache) Get(key string) (interface{}, error) {
	switch cache.Type {
	case "self":
		if val, hasKey := cache.Self.Get(key); hasKey {
			return val, nil
		}
		return nil, fmt.Errorf("no cache")
	case "redis":
		return cache.Redis.Get(ctx, key).Result()
	case "none":
		return nil, fmt.Errorf("no cache")
	}
	return nil, fmt.Errorf("no cache type")
}
