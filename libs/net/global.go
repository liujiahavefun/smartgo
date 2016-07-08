package net

import (
	"smartgo/libs/pool"
	. "smartgo/libs/utils"
)

func init() {
	netIdentifier = NewAtomicInt64(0)
	globalPool = pool.NewBufferPool(GLOBAL_POOL_DEFAULT_BUFFER_SIZE, GLOBAL_POOL_DEFAULT_BUFFER_NUM)
}

const (
	GLOBAL_POOL_DEFAULT_BUFFER_SIZE = 1024
	GLOBAL_POOL_DEFAULT_BUFFER_NUM  = 1024
)

var (
	netIdentifier *AtomicInt64
	globalPool    *pool.BufferPool
)
