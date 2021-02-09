// 锁原理 https://redis.io/topics/distlock
// 锁的池子 pools 参数，在redis集群下，尽可能多的提供pools
package tkredis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	tke "github.com/ikaiguang/srv_toolkit/error"
	tkredigo "github.com/ikaiguang/srv_toolkit/redis/redigo"
	"sync"
	"time"
)

// NewLock .
func NewLock(ctx context.Context, name string) (lock *DLock, isLockFailed bool, err error) {
	lock = &DLock{tries: 1}
	isLockFailed, err = lock.Lock(ctx, name)
	return
}

// NewLockWithTries .
func NewLockWithTries(ctx context.Context, name string, opt *DLockOption) (lock *DLock, isLockFailed bool, err error) {
	lock = &DLock{tries: opt.Tries, tryDelay: opt.TryDelay}
	isLockFailed, err = lock.Lock(ctx, name)
	return
}

// GetLock .
func GetLock(ctx context.Context, name string) (lock *DLock, err error) {
	lock, isLockFail, err := NewLock(ctx, name)
	if isLockFail {
		err = tke.Newf(tke.TooManyRequests, err)
		return
	}
	if err != nil {
		err = tke.Newf(tke.Redis, err)
		return
	}
	return
}

// GetLockWithTries .
func GetLockWithTries(ctx context.Context, name string, opt *DLockOption) (lock *DLock, err error) {
	lock, isLockFail, err := NewLockWithTries(ctx, name, opt)
	if isLockFail {
		err = tke.Newf(tke.TooManyRequests, err)
		return
	}
	if err != nil {
		err = tke.Newf(tke.Redis, err)
		return
	}
	return
}

// 锁信息
const (
	_lockExpire     = 8 * time.Second        // 锁过期时间
	_lockResetDelay = 3 * time.Second        // 重置锁时间，防止锁自动过期
	_lockTries      = 1                      // 尝试次数
	_lockDelay      = 100 * time.Millisecond // 尝试间隔
)

// lock
var (
	_lockOnce sync.Once
	_lockSync *redsync.Redsync
)

// DLock redis 分布式锁
type DLock struct {
	tries     int
	tryDelay  time.Duration
	mutex     *redsync.Mutex // 锁
	resetChan chan bool      // 信道
}

// DLockOption .
type DLockOption struct {
	Tries    int
	TryDelay time.Duration
}

// lazySync
func (s *DLock) lazySync() *redsync.Redsync {
	_lockOnce.Do(func() {
		_lockSync = redsync.New(tkredigo.NewPool(Redis()))
	})
	return _lockSync
}

// Lock 锁
func (s *DLock) Lock(ctx context.Context, name string) (isLockFailed bool, err error) {
	// mutex
	s.mutex = s.lazySync().NewMutex(name, s.Options()...)
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
func (s *DLock) Unlock() (bool, error) {
	if s.resetChan != nil {
		s.resetChan <- true
		close(s.resetChan)
	}
	return s.mutex.Unlock()
}

// Options .
func (s *DLock) Options() (opts []redsync.Option) {
	if s.tries <= 1 {
		s.tries = _lockTries
	}
	if s.tryDelay <= time.Millisecond {
		s.tryDelay = _lockDelay
	}

	opts = []redsync.Option{
		redsync.WithExpiry(_lockExpire),
		redsync.WithTries(s.tries),
		redsync.WithRetryDelay(s.tryDelay),
	}
	return
}

// resetExpire 重置锁时间，防止自动过期而解锁
func (s *DLock) resetExpire() {
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
