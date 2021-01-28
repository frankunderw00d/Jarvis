package database

import (
	redisGo "github.com/gomodule/redigo/redis"
	"time"
)

type (
	// Redis 锁定义
	RedisLock interface {
		// 加锁
		Lock(string, int) bool

		// 解锁
		UnLock(string) bool

		// 关闭
		Close() error
	}

	// Redis 锁定义实现
	redisLock struct {
		redisConn redisGo.Conn // redis 连接
		ex        int          // 加锁超时时间(以 second 为单位)，防止死锁
	}
)

const (
	// set 加锁命令
	RedisLockCommand = "set"
	// set 加锁值
	RedisLockValue = "1"
	// set 加锁超时命令
	RedisLockEXCommand = "ex"
	// set 加锁排他命令
	RedisLockIfNotExist = "nx"
	// del 解锁命令
	RedisUnlockCommand = "del"
	// 默认加锁超时时间
	DefaultEX = 2
	// 最大尝试加锁次数 0.1+0.2+0.4+0.8+1.6+3.2+6.4+12.8+25.6+51.2 = 102.3 second
	MaxLockTryCount = 10
)

var ()

// 新建 Redis 锁
func NewRedisLock(conn redisGo.Conn, ex int) RedisLock {
	if ex <= 0 {
		ex = DefaultEX
	}

	return &redisLock{
		redisConn: conn,
		ex:        ex,
	}
}

// 加锁
func (rl *redisLock) Lock(key string, try int) bool {
	if rl.redisConn == nil {
		return false
	}

	count := 1
	waitTime := time.Millisecond * time.Duration(100)

	for count <= try {
		_, err := redisGo.String(rl.redisConn.Do(RedisLockCommand, key, RedisLockValue, RedisLockEXCommand, rl.ex, RedisLockIfNotExist))
		if err == nil {
			return true
		} else {
			time.Sleep(waitTime)
			count++
			waitTime *= 2
		}
	}

	return false
}

// 解锁
func (rl *redisLock) UnLock(key string) bool {
	if rl.redisConn == nil {
		return false
	}

	reply, err := redisGo.Int(rl.redisConn.Do(RedisUnlockCommand, key))
	if err != nil {
		return false
	}

	return reply == 1
}

// 关闭
func (rl *redisLock) Close() error {
	return rl.redisConn.Close()
}
