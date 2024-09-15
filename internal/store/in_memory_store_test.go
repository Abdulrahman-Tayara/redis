package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	s := NewInMemoryStore()

	err := s.Set("key", "val")
	assert.NoError(t, err)

	val, err := s.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "val", val)
}

func TestInMemoryStore_Keys(t *testing.T) {
	s := NewInMemoryStore()

	_ = s.Set("key1", "val")
	_ = s.Set("key2", "val")
	_ = s.Set("key3", "val")

	_ = s.Set("username", "abdulrahman")

	keys, err := s.Keys("*")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(keys))

	keys, err = s.Keys("usernam?")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(keys))
}
