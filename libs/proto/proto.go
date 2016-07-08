package proto

import (
	"encoding/binary"
	"smartgo/libs/pool"
)

func init() {
	ByteOrder = binary.BigEndian
}

var (
	ByteOrder binary.ByteOrder
)

type Packet interface {
	MajorID() uint16
	MinorID() uint16
	MarshalBinary() ([]byte, error)
}

type UpPacket struct {
	major  uint16
	minor  uint16
	buffer *pool.Buffer
}

func NewUpPacket(buffer *pool.Buffer) *UpPacket {
	packet := &UpPacket{
		buffer: buffer,
	}

	packet.DecodeUpMessage()
	return packet
}

func (packet *UpPacket) DecodeUpMessage() {
	packet.buffer.ResetSeeker()

	if ByteOrder == binary.BigEndian {
		packet.major = packet.buffer.ReadUint16BE()
	} else {
		packet.major = packet.buffer.ReadUint16LE()
	}

	if ByteOrder == binary.BigEndian {
		packet.minor = packet.buffer.ReadUint16BE()
	} else {
		packet.minor = packet.buffer.ReadUint16LE()
	}
}

func (packet *UpPacket) MajorID() uint16 {
	return packet.major
}

func (packet *UpPacket) MinorID() uint16 {
	return packet.minor
}

func (packet *UpPacket) MarshalBinary() ([]byte, error) {
	return packet.buffer.Serialize()
}
