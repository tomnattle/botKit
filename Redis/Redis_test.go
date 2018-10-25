package Redis

import (
	"math/rand"
	"testing"
)

func TestRedis(t *testing.T) {
	redis, err := GetRedis()
	if err != nil {
		t.Fatal(err)
	}
	defer redis.Close()

	random := rand.Intn(100)
	err = redis.Cmd("SET", FormatKey("TESTKEY"), random).Err
	if err != nil {
		t.Fatal(err)
	}
	result, err := redis.Cmd("GET", FormatKey("TESTKEY")).Int()
	if err != nil {
		t.Fatal(err)
	}
	if result != random {
		t.Fatalf("input:%v != output:%v", random, result)
	}
}
