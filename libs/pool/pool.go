package pool

import (
	"sync"
)

const (
	BUFFER_POOL_TOTAL_LIMIT         = 1024 * 1024 * 512
	BUFFER_POOL_BUFFER_LIMIT        = 1024 * 512
	BUFFER_POOL_DEFAULT_BUFFER_SIZE = 1024
)

type BufferPool struct {
	free  *Buffer
	lock  sync.Mutex
	total int64
	size  int
	num   int
}

func NewBufferPool(size int, num int) *BufferPool {
	pool := &BufferPool{
		size: size,
		num:  num,
		free: nil,
	}

	pool.grow()

	return pool
}

func (pool *BufferPool) BufferSize() int {
	return pool.size
}

func (pool *BufferPool) grow() {
	for i := 0; i < pool.num; i++ {
		buffer := newBuffer(pool)
		buffer.next = pool.free
		pool.free = buffer
	}

	pool.total = int64(pool.size) * int64(pool.num)
}

func (pool *BufferPool) GetInBuffer() (buf *Buffer) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	buffer := pool.free
	if buffer != nil {
		buffer.isFreed = false

		pool.free = buffer.next
		pool.total -= int64(cap(buffer.Data))
		return buffer
	}

	return newBuffer(pool)
}

func (pool *BufferPool) PutInBuffer(buf *Buffer) {
	if cap(buf.Data) >= BUFFER_POOL_BUFFER_LIMIT || pool.total+int64(cap(buf.Data)) >= BUFFER_POOL_TOTAL_LIMIT {
		return
	}

	pool.lock.Lock()
	defer pool.lock.Unlock()

	buf.Data = buf.Data[0:0]
	buf.ReadPos = 0
	buf.isFreed = true

	buf.next = pool.free
	pool.free = buf

	pool.total += int64(cap(buf.Data))
}
