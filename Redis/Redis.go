package Redis

import (
	"fmt"
	"github.com/ifchange/botKit/config"
	redisCluster "github.com/mediocregopher/radix.v2/cluster"
	"github.com/mediocregopher/radix.v2/redis"
	"net"
	"time"
)

func init() {
	cfg = config.GetConfig().Redis
	if cfg == nil {
		panic("Redis config is nil")
	}
	redis, err := getRedis(cfg)
	if err != nil {
		panic(fmt.Errorf("init config error, redis start error, %v %v",
			cfg.Addr, err))
	}
	redis.Close()
}

var (
	cfg *config.RedisConfig
)

type RedisCommon struct {
	isCluster bool
	cluster   *redisCluster.Cluster
	client    *redis.Client
}

func (ins *RedisCommon) Cmd(cmd string, args ...interface{}) *redis.Resp {
	if ins.isCluster {
		return ins.cluster.Cmd(cmd, args...)
	}
	return ins.client.Cmd(cmd, args...)
}

func (ins *RedisCommon) Close() {
	if ins.isCluster {
		ins.cluster.Close()
		return
	}
	ins.client.Close()
}

func GetRedis() (*RedisCommon, error) {
	ins, err := getRedis(cfg)
	if err != nil {
		return nil, fmt.Errorf("try get redis ins error %v %v",
			cfg, err)
	}
	return ins, nil
}

func getRedis(cfg *config.RedisConfig) (*RedisCommon, error) {
	ins := &RedisCommon{
		isCluster: cfg.IsCluster,
	}
	if ins.isCluster {
		cluster, err := redisCluster.NewWithOpts(
			redisCluster.Opts{
				Addr:             cfg.Addr,
				PoolSize:         cfg.PoolSize,
				MaxRedirectCount: 5,
			})
		if err != nil {
			return ins, err
		}
		ins.cluster = cluster
		return ins, nil
	}
	conn, err := net.DialTimeout("tcp", cfg.Addr, time.Second * 30)
	if err != nil {
		return ins, err
	}
	client, err := redis.NewClient(conn)
	if err != nil {
		return ins, err
	}
	ins.client = client
	return ins, nil
}

func FormatKey(key string) string {
	return cfg.Prefix + key
}
