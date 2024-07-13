package interfaces

type CacheManager interface {
	Get(key string) ([]byte, error)
	Set(key string, value any) error
}
