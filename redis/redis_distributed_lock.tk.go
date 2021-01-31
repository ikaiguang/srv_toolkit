// 锁原理 https://redis.io/topics/distlock
// 锁的池子 pools 参数，在redis集群下，尽可能多的提供pools
package tkredis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	tkredigo "github.com/ikaiguang/srv_toolkit/redis/redigo"
	"sync"
	"time"
)

// NewLock .
func NewLock(name string) (lock *DistributedLock, isLockFailed bool, err error) {
	lock = new(DistributedLock)
	isLockFailed, err = lock.Lock(context.Background(), name)
	return
}

// NewLockContext .
func NewLockContext(ctx context.Context, name string) (lock *DistributedLock, isLockFailed bool, err error) {
	lock = new(DistributedLock)
	isLockFailed, err = lock.Lock(ctx, name)
	return
}

// 锁信息
const (
	_lockExpire     = 8 * time.Second // 锁过期时间
	_lockResetDelay = 3 * time.Second // 重置锁时间，防止锁自动过期
	_lockTries      = 1               // 尝试次数
)

// lock
var (
	_lockOnce sync.Once
	_lockSync *redsync.Redsync
	_lockOpts []redsync.Option
)

// DistributedLock redis 分布式锁
type DistributedLock struct {
	mutex     *redsync.Mutex // 锁
	resetChan chan bool      // 信道
}

// lazySync
func (s *DistributedLock) lazySync() (*redsync.Redsync, []redsync.Option) {
	_lockOnce.Do(func() {
		_lockOpts = s.Options()
	})
	_lockSync = redsync.New(tkredigo.NewPool(redisConn))
	return _lockSync, _lockOpts
}

// Lock 锁
func (s *DistributedLock) Lock(ctx context.Context, name string) (isLockFailed bool, err error) {
	// mutex
	lockSync, lockOpts := s.lazySync()
	s.mutex = lockSync.NewMutex(name, lockOpts...)
	err = s.mutex.LockContext(ctx)
	if err != nil {
		isLockFailed = err == redsync.ErrFailed
		return
	}

	// 续期锁，防止锁自动过期
	s.resetChan = make(chan bool)
	go s.resetExpire()

	return
}

// Unlock 解锁
func (s *DistributedLock) Unlock() (bool, error) {
	if s.resetChan != nil {
		s.resetChan <- true
		close(s.resetChan)
	}
	return s.mutex.Unlock()
}

// Options .
func (s *DistributedLock) Options() (opts []redsync.Option) {
	opts = []redsync.Option{
		redsync.WithExpiry(_lockExpire),
		redsync.WithTries(_lockTries),
	}
	return
}

// resetExpire 重置锁时间，防止自动过期而解锁
func (s *DistributedLock) resetExpire() {
	// 计时器
	timer := time.NewTimer(_lockResetDelay)

	select {
	case <-timer.C: // 续期
		// 结束计时
		timer.Stop()
		// 续期
		if ok, err := s.mutex.Extend(); err != nil || !ok {
			// 调试
			//fmt.Println("redis mutex 续期失败")
			return
		}
		// 调试
		//fmt.Println("redis mutex 续期成功")
		// 再次续期
		s.resetExpire()
	case <-s.resetChan: // 停止
		timer.Stop()
		// 调试
		//fmt.Println("redis mutex 停止续期")
		return
	}
}
