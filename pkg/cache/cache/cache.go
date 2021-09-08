package cache

import (
	"time"

	cacheNew "github.com/patrickmn/go-cache"
)

func Init(expTime time.Duration) *cacheNew.Cache {
	return cacheNew.New(5*time.Minute, expTime*time.Minute)
}
