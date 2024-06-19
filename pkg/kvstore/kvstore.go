package kvstore

type KVStore interface {
	Set(key string, value any) error
	Get(key string) (any, error)
	Delete(key string) error
	Exist(key string) (bool, error)
	Ping() error
	Close() error
}
