package store

import (
	"errors"
	"path/filepath"
	"sync"
)

var ErrKeyNotFound = errors.New("KEY_NOT_FOUND")

type inMemoryStore struct {
	data map[string]any

	sync.Mutex
}

func NewInMemoryStore() Store {
	return &inMemoryStore{
		data: make(map[string]any),
	}
}

func (s *inMemoryStore) Set(key string, val any) error {
	defer s.Unlock()

	defer s.Lock()

	s.data[key] = val

	return nil
}

func (s *inMemoryStore) Get(key string) (any, error) {
	defer s.Unlock()

	s.Lock()

	val, ok := s.data[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	return val, nil
}

func (s *inMemoryStore) Keys(pattern string) ([]string, error) {
	defer s.Unlock()

	s.Lock()

	size := 0

	for key, _ := range s.data {
		if isRegexMatched(key, pattern) {
			size++
		}
	}

	keys := make([]string, size)
	i := 0

	for key, _ := range s.data {
		if isRegexMatched(key, pattern) {
			keys[i] = key
			i++
		}
	}

	return keys, nil
}

func isRegexMatched(val string, pattern string) bool {
	matched, _ := filepath.Match(pattern, val)

	return matched
}
