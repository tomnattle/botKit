package cache

import (
	"testing"
	"time"
)

func TestCacheValue(t *testing.T) {
	myCache, err := New()
	if err != nil {
		t.Fatal(err)
	}
	str := "value"
	myCache.Set("key", str, time.Duration(1)*time.Minute)
	cacheStr, ok := myCache.Get("key").(string)
	if !ok {
		t.Fatalf("cache string not exist")
	}
	if cacheStr != str {
		t.Fatalf("cache string not equal")
	}
}

func TestCachePointer(t *testing.T) {
	myCache, err := New()
	if err != nil {
		t.Fatal(err)
	}
	str := "value"
	myCache.Set("key", &str, time.Duration(1)*time.Minute)
	cacheStr, ok := myCache.Get("key").(*string)
	if !ok {
		t.Fatalf("cache string not exist")
	}
	if *cacheStr != str {
		t.Fatalf("cache string not equal")
	}
}

func TestCacheDel(t *testing.T) {
	myCache, err := New()
	if err != nil {
		t.Fatal(err)
	}
	str := "value"
	myCache.Set("key", str, time.Duration(1)*time.Second)
	cacheStr, ok := myCache.Get("key").(string)
	if !ok {
		t.Fatalf("cache string not exist")
	}
	if cacheStr != str {
		t.Fatalf("cache string not equal")
	}
	myCache.Del("key")
	_, ok = myCache.Get("key").(string)
	if ok {
		t.Fatalf("cache expire error")
	}
}

func TestCacheExpire(t *testing.T) {
	myCache, err := New()
	if err != nil {
		t.Fatal(err)
	}
	str := "value"
	myCache.Set("key", str, time.Duration(1)*time.Second)
	cacheStr, ok := myCache.Get("key").(string)
	if !ok {
		t.Fatalf("cache string not exist")
	}
	if cacheStr != str {
		t.Fatalf("cache string not equal")
	}
	time.Sleep(time.Duration(1) * time.Second)
	_, ok = myCache.Get("key").(string)
	if ok {
		t.Fatalf("cache expire error")
	}
}
