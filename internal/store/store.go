package store

import (
	"redis/pkg/ds"
)

type Store interface {
	HashTable() *ds.ExpiringHashTable
}

type store struct {
	hashTable *ds.ExpiringHashTable
}

func NewStore(ht *ds.ExpiringHashTable) Store {
	s := store{
		hashTable: ht,
	}
	return &s
}

func (s *store) HashTable() *ds.ExpiringHashTable {
	return s.hashTable
}
