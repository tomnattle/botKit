package lockPool

import (
	"fmt"
	"github.com/ifchange/botKit/Redis"
	"github.com/ifchange/botKit/logger"
	"time"
)

type LockPool struct {
	key    func(interface{}) string
	size   int
	expire time.Duration

	conn *Redis.RedisCommon
}

func New(prefix string) (*LockPool, error) {
	return NewWithConfig(prefix, 1, time.Duration(30)*time.Second)
}

func NewWithConfig(prefix string, size int, expire time.Duration) (*LockPool, error) {
	if len(prefix) == 0 {
		return nil, fmt.Errorf("Init lock pool error prefix unauthorized")
	}
	if size <= 0 {
		return nil, fmt.Errorf("Init lock pool error size unauthorized")
	}
	if expire <= 0 {
		return nil, fmt.Errorf("Init lock pool error expire unauthorized")
	}

	conn, err := Redis.GetRedis()
	if err != nil {
		return nil, fmt.Errorf("Init lock pool error Redis error %v", err)
	}
	return &LockPool{
		key:    func(unique interface{}) string { return Redis.FormatKey(fmt.Sprintf("%s_%v", prefix, unique)) },
		size:   size,
		expire: expire,
		conn:   conn,
	}, nil
}

func (lp *LockPool) Lock(unique interface{}) {
	pollTimes := time.Duration(0)
	for {
		// polling
		exist, err := lp.conn.Cmd("GET", lp.key(unique)).Int()
		if err != nil {
			break
		}
		if exist < lp.size {
			break
		}
		duration := time.Duration(5) * time.Millisecond
		time.Sleep(duration)
		pollTimes += duration
		if pollTimes >= lp.expire {
			break
		}
	}
	lp.conn.Cmd("INCR", lp.key(unique))
	lp.conn.Cmd("EXPIRE", lp.key(unique), int(lp.expire.Seconds()))
}

func (lp *LockPool) Unlock(unique interface{}) {
	exist, err := lp.conn.Cmd("DECR", lp.key(unique)).Int()
	if err != nil {
		return
	}
	if exist <= 0 {
		lp.conn.Cmd("DEL", lp.key(unique))
	}
}

func (lp *LockPool) Status(unique interface{}) (exist int) {
	exist, err := lp.conn.Cmd("GET", lp.key(unique)).Int()
	if err != nil {
		logger.Printf("lockPool ins %v get status error %v",
			lp, err)
		return 0
	}
	return
}
