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
	str := "string"
	myCache.Set("string", str, time.Duration(1)*time.Minute)
	cacheStr, ok := myCache.Get("string").(string)
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
	str := "string"
	myCache.Set("string", &str, time.Duration(1)*time.Minute)
	cacheStr, ok := myCache.Get("string").(*string)
	if !ok {
		t.Fatalf("cache string not exist")
	}
	if *cacheStr != str {
		t.Fatalf("cache string not equal")
	}
}

func TestCacheExpire(t *testing.T) {
	myCache, err := New()
	if err != nil {
		t.Fatal(err)
	}
	str := "string"
	myCache.Set("string", str, time.Duration(1)*time.Second)
	cacheStr, ok := myCache.Get("string").(string)
	if !ok {
		t.Fatalf("cache string not exist")
	}
	if cacheStr != str {
		t.Fatalf("cache string not equal")
	}
	time.Sleep(time.Duration(1) * time.Second)
	_, ok = myCache.Get("string").(string)
	if ok {
		t.Fatalf("cache expire error")
	}
}
