package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func Init(conf *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password, // no password set
		DB:       conf.DB,       // use default DB
	})
}
