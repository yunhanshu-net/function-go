package runner

import (
	"time"
)

type Locker interface {

	// Lock 尝试加锁，返回是否加锁成功，非阻塞
	Lock(key string, ttl ...time.Duration) bool

	// Unlock 删除锁，返回是否删除成功
	Unlock(key string) bool
}

type Lock struct {
}

func newLock() *Lock {
	return &Lock{}
}

func (l *Lock) Lock(key string, ttl ...time.Duration) bool {
	//todo待实现
	return true
}
func (l *Lock) Unlock(key string) bool {
	//todo待实现
	return true
}
