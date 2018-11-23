package cache

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// type CacheType uint8

// const (
// 	_ CacheType = iota
// 	MemoryCache
// 	RedisCache
// )

// type Cache struct {
// 	cacheSaver saver
// 	created    time.Time
// 	lock       *sync.Mutex
// }

// type saver interface {
// 	get(string) interface{}
// 	set(string, interface{}, time.Duration)
// 	release(string)
// 	releaseAll()
// }

// func New(t CacheType) (*Cache, error) {
// 	if t == MemoryCache {

// 	}
// }

// func (ins *Cache) Get(key string) interface{} {

// }

// func (ins *Cache) Set(key string, value interface{}, expire time.Duration) {

// }

// func (ins *Cache) Del(key string) {

// }

// type item struct {
// 	value  interface{}
// 	expire time.Time
// }
