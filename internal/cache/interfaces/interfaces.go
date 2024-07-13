package interfaces

import "github.com/go-redsync/redsync/v4"

type CacheManager interface {
	Get(key string) ([]byte, error)
	Set(key string, value any) error
	AcquireLock(key string) (*redsync.Mutex, error)
	SafeUnlock(mutex *redsync.Mutex)
	DeleteByKeys(keys []string) error
}
