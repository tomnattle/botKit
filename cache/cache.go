/*
achieve stupid time-based memory cache
lazy release big map-cap per day

// init
myCache := cache.New()

// set
var value1 TYPE
myCache.Set("key1", value1, time.Duration(1)*time.Second)

// get
value1, ok := myCache.Get("key1").(TYPE)
if !ok {
	// cache not exist
}

// delete key (always safe)
myCache.Del("key1")
*/
package cache

import (
	"sync"
	"time"
)

type Cache struct {
	storage  map[string]*item
	released time.Time
	lock     *sync.RWMutex
}

type item struct {
	value  interface{}
	expire time.Time
}

func New() (*Cache, error) {
	return &Cache{
		storage:  make(map[string]*item),
		released: time.Now(),
		lock:     new(sync.RWMutex),
	}, nil
}

func (ins *Cache) Get(key string) interface{} {
	ins.lock.RLock()
	defer ins.lock.RUnlock()

	value, ok := ins.storage[key]
	if !ok {
		return nil
	}
	if value.expire.Before(time.Now()) {
		return nil
	}
	return value.value
}

func (ins *Cache) Set(key string, value interface{}, expire time.Duration) {
	now := time.Now()
	ins.lock.Lock()
	defer ins.lock.Unlock()
	defer ins.releseCap(now)

	ins.storage[key] = &item{
		value:  value,
		expire: now.Add(expire),
	}
}

func (ins *Cache) Del(key string) {
	ins.lock.Lock()
	defer ins.lock.Unlock()
	defer ins.releseCap(time.Now())

	delete(ins.storage, key)
}

// MUST handle in lock
func (ins *Cache) releseCap(now time.Time) {
	if now.Sub(ins.released) < time.Duration(24)*time.Hour {
		return
	}
	storage := make(map[string]*item)
	for key, value := range ins.storage {
		if value.expire.Before(now) {
			continue
		}
		storage[key] = value
	}
	ins.storage = storage
	ins.released = now
}
