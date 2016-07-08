package pool

type RawBuffer struct {
	Buffer
}

func newRawBuffer() *RawBuffer {
	return &RawBuffer{
		Buffer{
			Data:    make([]byte, 0, BUFFER_POOL_DEFAULT_BUFFER_SIZE),
			ReadPos: 0,
			isFreed: false,
			pool:    nil,
			next:    nil,
		},
	}
}

func (buf *RawBuffer) Free() {
	if buf.isFreed {
		panic("RawBuffer: double free")
	}

	buf.isFreed = true
}
