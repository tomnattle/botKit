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
	storage     *sync.Map
	released    time.Time
	releaseLock *sync.Mutex
}

type item struct {
	value  interface{}
	expire time.Time
}

func New() (*Cache, error) {
	return &Cache{
		storage:     new(sync.Map),
		released:    time.Now(),
		releaseLock: new(sync.Mutex),
	}, nil
}

func (ins *Cache) Get(key interface{}) interface{} {
	value, ok := ins.storage.Load(key)
	if !ok {
		return nil
	}
	i, ok := value.(*item)
	if !ok {
		return nil
	}
	if i.expire.Before(time.Now()) {
		return nil
	}
	return i.value
}

func (ins *Cache) Set(key interface{}, value interface{}, expire time.Duration) {
	now := time.Now()
	ins.releaseLock.Lock()
	defer ins.releaseLock.Unlock()
	defer ins.releseCap(now)

	ins.storage.Store(key, &item{
		value:  value,
		expire: now.Add(expire),
	})
}

func (ins *Cache) Del(key interface{}) {
	ins.releaseLock.Lock()
	defer ins.releaseLock.Unlock()
	defer ins.releseCap(time.Now())

	ins.storage.Delete(key)
}

// MUST handle in releaseLock
func (ins *Cache) releseCap(now time.Time) {
	if now.Sub(ins.released) < time.Duration(24)*time.Hour {
		return
	}
	storage := new(sync.Map)
	scanAll :=
		func(key, value interface{}) (continueIteration bool) {
			continueIteration = true
			i, ok := value.(*item)
			if !ok {
				return
			}
			if i.expire.Before(now) {
				return
			}
			storage.Store(key, value)
			return
		}
	ins.storage.Range(scanAll)
	ins.storage = storage
	ins.released = now
}
