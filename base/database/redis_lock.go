package database

import (
	redisGo "github.com/gomodule/redigo/redis"
)

type (
	// Redis 锁定义
	RedisLock interface {
		// 尝试能否获得锁的权利
		TryLock(string) bool

		// 加锁
		Lock()

		// 解锁
		UnLock()

		// 关闭
		Close() error
	}

	redisLock struct {
		redisConn redisGo.Conn
	}
)

const (
	RedisTryLockCommand = "get"
)

var ()

// 新建 Redis 锁
func NewRedisLock(conn redisGo.Conn) RedisLock {
	return &redisLock{redisConn: conn}
}

// 尝试能否获得锁的权利
func (rl *redisLock) TryLock(key string) bool {
	if rl.redisConn == nil {
		return false
	}

	_, err := redisGo.String(rl.redisConn.Do(RedisTryLockCommand, key))
	if err != nil && err == redisGo.ErrNil {
		return true
	} else {
		return false
	}
}

// 加锁
func (rl *redisLock) Lock() {}

// 解锁
func (rl *redisLock) UnLock() {}

// 关闭
func (rl *redisLock) Close() error {
	return nil
}
