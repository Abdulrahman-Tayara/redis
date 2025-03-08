package ds

import (
	"fmt"
	"testing"
	"time"
)

func TestExpiringHashTable(t *testing.T) {
	hs := NewExpiringHashTable()

	_ = hs.Set("key1", "value1", &SetOptions{
		ExpireAt: time.Now().Add(time.Hour).UnixMilli(),
	})

	if v, ok := hs.Get("key1"); !ok || v != "value1" {
		t.Errorf("expected value1, got %v", v)
	}

	_ = hs.Set("key2", "value2", &SetOptions{
		ExpireAt: time.Now().Add(time.Millisecond * 50).UnixMilli(),
	})

	time.Sleep(time.Millisecond * 100)

	if _, ok := hs.Get("key2"); ok {
		t.Errorf("expected key2 to be expired")
	}
}

func TestExpiringHashTable_Remove(t *testing.T) {
	hs := NewExpiringHashTable()

	_ = hs.Set("key1", "value1", &SetOptions{
		ExpireAt: time.Now().Add(time.Hour).UnixMilli(),
	})

	if _, ok := hs.Get("key1"); !ok {
		t.Errorf("expected key1 to be present")
	}

	hs.Remove("key1")

	if _, ok := hs.Get("key1"); ok {
		t.Errorf("expected key1 to be removed")
	}
}

func TestExpiringHashTable_KeysChunk(t *testing.T) {
	hs := NewExpiringHashTable()

	for i := 0; i < 100; i++ {
		_ = hs.Set(fmt.Sprintf("key%d", i), i+1, &SetOptions{
			ExpireAt: time.Now().Add(time.Hour).UnixMilli(),
		})
	}

	keys := hs.KeysChunk(10)

	if len(keys) != 10 {
		t.Errorf("expected 10 keys, got %d", len(keys))
	}

	keys = hs.KeysChunk(200)

	if len(keys) != 100 {
		t.Errorf("expected 100 keys, got %d", len(keys))
	}
}

func TestExpiringHashTable_RemoveIfExpiredMany(t *testing.T) {
	hs := NewExpiringHashTable()

	for i := 0; i < 100; i++ {
		_ = hs.Set(fmt.Sprintf("key%d", i), i+1, &SetOptions{
			ExpireAt: time.Now().Add(time.Second).UnixMilli(),
		})
	}

	time.Sleep(time.Second * 2)

	hs.RemoveIfExpiredMany([]string{"key1", "key2"})

	if _, ok := hs.Get("key1"); ok {
		t.Errorf("expected delted key %s", "key1")
	}
}
