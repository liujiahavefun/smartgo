package pool

import (
	"encoding/binary"
	"io"
	"math"
	"unicode/utf8"
)

type Buffer struct {
	Data    []byte
	ReadPos int
	isFreed bool
	pool    *BufferPool
	next    *Buffer
}

func newBuffer(pool *BufferPool) *Buffer {
	return &Buffer{
		Data:    make([]byte, 0, pool.BufferSize()),
		ReadPos: 0,
		isFreed: false,
		pool:    pool,
		next:    nil,
	}
}

func (buf *Buffer) Free() {
	if buf.isFreed {
		panic("Buffer: double free")
	}

	buf.pool.PutInBuffer(buf)
}

func (buf *Buffer) Prepare(size int) {
	if cap(buf.Data) < size {
		buf.Data = make([]byte, size)
	} else {
		buf.Data = buf.Data[0:size]
	}
}

func (buf *Buffer) ResetSeeker() {
	buf.ReadPos = 0
}

func (buf *Buffer) RawData() []byte {
	return buf.Data
}

func (buf *Buffer) Serialize() ([]byte, error) {
	return buf.Data, nil
}

//implements io.Reader interface
func (buf *Buffer) Read(b []byte) (int, error) {
	if buf.ReadPos == len(buf.Data) {
		return 0, io.EOF
	}
	n := len(b)
	if n+buf.ReadPos > len(buf.Data) {
		n = len(buf.Data) - buf.ReadPos
	}
	copy(b, buf.Data[buf.ReadPos:])
	buf.ReadPos += n
	return n, nil
}

//implements io.Writer interface
func (buf *Buffer) Write(p []byte) (int, error) {
	buf.Data = append(buf.Data, p...)
	return len(p), nil
}

func (buf *Buffer) Slice(n int) []byte {
	if buf.ReadPos+n > len(buf.Data) {
		panic("Buffer: read buf of range")
	}
	r := buf.Data[buf.ReadPos : buf.ReadPos+n]
	buf.ReadPos += n
	return r
}

func (buf *Buffer) ReadBytes(n int) []byte {
	x := make([]byte, n)
	copy(x, buf.Slice(n))
	return x
}

func (buf *Buffer) ReadString() string {
	l := buf.ReadUint32BE()
	return string(buf.Slice(int(l)))
}

func (buf *Buffer) ReadRune() rune {
	x, size := utf8.DecodeRune(buf.Data[buf.ReadPos:])
	buf.ReadPos += size
	return x
}

func (buf *Buffer) ReadUint8() uint8 {
	return uint8(buf.Slice(1)[0])
}

func (buf *Buffer) ReadUint16LE() uint16 {
	return binary.LittleEndian.Uint16(buf.Slice(2))
}

func (buf *Buffer) ReadUint16BE() uint16 {
	return binary.BigEndian.Uint16(buf.Slice(2))
}

func (buf *Buffer) ReadUint32LE() uint32 {
	return binary.LittleEndian.Uint32(buf.Slice(4))
}

func (buf *Buffer) ReadUint32BE() uint32 {
	return binary.BigEndian.Uint32(buf.Slice(4))
}

func (buf *Buffer) ReadUint64LE() uint64 {
	return binary.LittleEndian.Uint64(buf.Slice(8))
}

func (buf *Buffer) ReadUint64BE() uint64 {
	return binary.BigEndian.Uint64(buf.Slice(8))
}

func (buf *Buffer) ReadFloat32LE() float32 {
	return math.Float32frombits(buf.ReadUint32LE())
}

func (buf *Buffer) ReadFloat32BE() float32 {
	return math.Float32frombits(buf.ReadUint32BE())
}

func (buf *Buffer) ReadFloat64LE() float64 {
	return math.Float64frombits(buf.ReadUint64LE())
}

func (buf *Buffer) ReadFloat64BE() float64 {
	return math.Float64frombits(buf.ReadUint64BE())
}

func (buf *Buffer) ReadVarint() int64 {
	v, n := binary.Varint(buf.Data[buf.ReadPos:])
	buf.ReadPos += n
	return v
}

func (buf *Buffer) ReadUvarint() uint64 {
	v, n := binary.Uvarint(buf.Data[buf.ReadPos:])
	buf.ReadPos += n
	return v
}

func (buf *Buffer) Append(p ...byte) {
	buf.Data = append(buf.Data, p...)
}

func (buf *Buffer) WriteBytes(d []byte) {
	buf.Append(d...)
}

func (buf *Buffer) WriteString(s string) {
	buf.WriteUint32BE(uint32(len([]byte(s))))
	buf.Append([]byte(s)...)
}

func (buf *Buffer) WriteRune(r rune) {
	p := []byte{0, 0, 0, 0}
	n := utf8.EncodeRune(p, r)
	buf.Append(p[:n]...)
}

func (buf *Buffer) WriteUint8(v uint8) {
	buf.Append(byte(v))
}

func (buf *Buffer) WriteUint16LE(v uint16) {
	buf.Append(byte(v), byte(v>>8))
}

func (buf *Buffer) WriteUint16BE(v uint16) {
	buf.Append(byte(v>>8), byte(v))
}

func (buf *Buffer) WriteUint32LE(v uint32) {
	buf.Append(byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

func (buf *Buffer) WriteUint32BE(v uint32) {
	buf.Append(byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (buf *Buffer) WriteUint64LE(v uint64) {
	buf.Append(
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56),
	)
}

func (buf *Buffer) WriteUint64BE(v uint64) {
	buf.Append(
		byte(v>>56),
		byte(v>>48),
		byte(v>>40),
		byte(v>>32),
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func (buf *Buffer) WriteFloat32LE(v float32) {
	buf.WriteUint32LE(math.Float32bits(v))
}

func (buf *Buffer) WriteFloat32BE(v float32) {
	buf.WriteUint32BE(math.Float32bits(v))
}

func (buf *Buffer) WriteFloat64LE(v float64) {
	buf.WriteUint64LE(math.Float64bits(v))
}

func (buf *Buffer) WriteFloat64BE(v float64) {
	buf.WriteUint64BE(math.Float64bits(v))
}

func (buf *Buffer) WriteUvarint(v uint64) {
	for v >= 0x80 {
		buf.Append(byte(v) | 0x80)
		v >>= 7
	}
	buf.Append(byte(v))
}

func (buf *Buffer) WriteVarint(v int64) {
	ux := uint64(v) << 1
	if v < 0 {
		ux = ^ux
	}
	buf.WriteUvarint(ux)
}
