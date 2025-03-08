package ds

import (
	"sync"
	"time"
)

type ExpiringHashTable struct {
	data          map[string]any
	expiringTable map[string]int64

	sync.RWMutex
}

func NewExpiringHashTable() *ExpiringHashTable {
	return &ExpiringHashTable{
		data:          make(map[string]any),
		expiringTable: make(map[string]int64),
	}
}

type SetOptions struct {
	// ExpireAt is unix timestamp milliseconds
	ExpireAt       int64
	SetIfNotExists bool
	SetIfExists    bool
	KeepTTL        bool
}

func (h *ExpiringHashTable) Set(key string, val any, opts *SetOptions) error {
	defer h.Unlock()

	h.Lock()

	var opt SetOptions
	if opts != nil {
		opt = *opts
	}

	shouldSet := true
	_, isExist := h.data[key]

	if opt.SetIfExists {
		shouldSet = isExist && opt.SetIfExists
	}
	if opt.SetIfNotExists {
		shouldSet = !isExist && opt.SetIfNotExists
	}
	if !shouldSet {
		return nil
	}

	h.data[key] = val
	if opt.ExpireAt > 0 && !opt.KeepTTL {
		h.expiringTable[key] = opt.ExpireAt
	}

	return nil
}

func (h *ExpiringHashTable) Get(key string) (any, bool) {
	if h.IsExpired(key) {
		h.Remove(key)
		return nil, false
	}

	defer h.RUnlock()

	h.RLock()

	v, ok := h.data[key]
	return v, ok
}

func (h *ExpiringHashTable) Remove(key string) {
	defer h.Unlock()

	h.Lock()
	h.remove(key)
}

func (h *ExpiringHashTable) RemoveMany(keys []string) {
	defer h.Unlock()

	h.Lock()

	for _, key := range keys {
		h.remove(key)
	}
}

func (h *ExpiringHashTable) remove(key string) {
	delete(h.data, key)
	delete(h.expiringTable, key)
}

func (h *ExpiringHashTable) KeysChunk(size int) []string {
	defer h.RUnlock()

	h.RLock()

	keys := make([]string, 0, min(len(h.data), size))

	i := 0
	for k := range h.data {
		keys = append(keys, k)
		i++
		if i >= size {
			break
		}
	}

	return keys
}

func (h *ExpiringHashTable) RemoveIfExpiredMany(keys []string) {
	expiredKeys := h.getExpiredKeys(keys)

	h.RemoveMany(expiredKeys)
}

func (h *ExpiringHashTable) IsExpired(key string) bool {
	defer h.RUnlock()

	h.RLock()

	if expiredAt, ok := h.expiringTable[key]; !ok {
		return false
	} else {
		return time.Now().UnixMilli() > expiredAt
	}
}

func (h *ExpiringHashTable) getExpiredKeys(keys []string) []string {
	defer h.RUnlock()

	h.RLock()

	var expiredKeys []string

	for _, key := range keys {
		if expiredAt, ok := h.expiringTable[key]; ok {
			if time.Now().UnixMilli() > expiredAt {
				expiredKeys = append(expiredKeys, key)
			}
		}
	}

	return expiredKeys
}
