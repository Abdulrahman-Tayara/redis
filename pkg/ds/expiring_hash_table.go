package ds

import (
	"sync"
	"time"
)

type OptionsFunc[T any] func(*T)

func getFilledOptions[T any](opts ...OptionsFunc[T]) *T {
	var opt T
	for _, f := range opts {
		f(&opt)
	}
	return &opt
}

type ExpiringHashTable struct {
	data          map[string]any
	expiringTable map[string]int64

	sync.RWMutex
}

func NewExpiringHashTable() *ExpiringHashTable {
	return &ExpiringHashTable{
		data: make(map[string]any),

		expiringTable: make(map[string]int64),
	}
}

type SetOptions struct {
	// ExpireAt is unix timestamp milliseconds
	ExpireAt int64
}

func SetWithExpireAt(expireAt int64) OptionsFunc[SetOptions] {
	return func(opts *SetOptions) {
		opts.ExpireAt = expireAt
	}
}

func (h *ExpiringHashTable) Set(key string, val any, opts ...OptionsFunc[SetOptions]) error {
	defer h.Unlock()

	h.Lock()

	opt := getFilledOptions(opts...)

	h.data[key] = val
	if opt.ExpireAt > 0 {
		h.expiringTable[key] = opt.ExpireAt
	}

	return nil
}

func (h *ExpiringHashTable) Get(key string) (any, error) {
	defer h.RUnlock()

	h.RLock()

	if h.isExpired(key) {
		h.Remove(key)
		return nil, nil
	}
	return h.data[key], nil
}

func (h *ExpiringHashTable) Remove(key string) {
	defer h.Unlock()

	h.Lock()
	delete(h.data, key)
	delete(h.expiringTable, key)
}

func (h *ExpiringHashTable) isExpired(key string) bool {
	defer h.RUnlock()

	h.RLock()

	if expiredAt, ok := h.expiringTable[key]; !ok {
		return false
	} else {
		return time.Now().UnixMilli() > expiredAt
	}
}
