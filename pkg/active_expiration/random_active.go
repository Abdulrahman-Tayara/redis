package active_expiration

import (
	"redis/logger"
	"redis/pkg/ds"
	"time"
)

const (
	DefaultIterChunkSize = 100
)

type RandomActiveExpiration struct {
	data *ds.ExpiringHashTable

	iterChunkSize int

	closed bool
}

func NewRandomActiveExpiration(data *ds.ExpiringHashTable) *RandomActiveExpiration {
	return &RandomActiveExpiration{
		data:          data,
		iterChunkSize: DefaultIterChunkSize,
	}
}

func (r *RandomActiveExpiration) Start() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("random active expiration panic: %v", err)
			time.Sleep(time.Second * 2)
			r.Start()
		}
	}()

	for {
		if r.closed {
			break
		}

		keys := r.data.KeysChunk(r.iterChunkSize)

		r.data.RemoveIfExpiredMany(keys)
	}
}

func (r *RandomActiveExpiration) Close() {
	r.closed = true
}
