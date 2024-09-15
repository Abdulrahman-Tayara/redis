package store

type Store interface {
	Set(key string, val any) error

	Get(key string) (any, error)

	Keys(pattern string) ([]string, error)
}
