package nocache

import "time"

func Init() *Cache {
	return &Cache{}
}

type Cache struct {
}

func (c *Cache) Set(k string, x interface{}, d time.Duration) {

}

func (c *Cache) Get(k string) (interface{}, bool) {
	return nil, false
}
