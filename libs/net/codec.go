package net

import (
	"encoding/binary"
	"github.com/golang/glog"
	"io"
	"smartgo/libs/pool"
)

/* Application programmer can define a custom codec themselves */
type Codec interface {
	Decode(Connection) (*pool.Buffer, error)
	Encode(Message) ([]byte, error)
}

// only support fix-length-header
type FixedLengthHeaderCodec struct {
	HeaderSize    int
	MinPacketSize int
	MaxPacketSize int
	ByteOrder     binary.ByteOrder
}

func NewFixedLengthHeaderCodec(header, min, max int, bo binary.ByteOrder) FixedLengthHeaderCodec {
	return FixedLengthHeaderCodec{
		HeaderSize:    header,
		MinPacketSize: min,
		MaxPacketSize: max,
		ByteOrder:     bo,
	}
}

func (codec FixedLengthHeaderCodec) Decode(conn Connection) (buffer *pool.Buffer, err error) {
	buffer = globalPool.GetInBuffer()

	buffer.Prepare(codec.HeaderSize)
	if _, err = io.ReadFull(conn.GetRawConn(), buffer.Data); err != nil {
		glog.Errorf("read packet header failed, err: %T %v\n", err, err)
		return
	}

	bodyLen := codec.decodeHead(buffer.Data)
	if codec.MaxPacketSize > 0 && bodyLen > codec.MaxPacketSize {
		err = ErrorkrPacketTooLarge
		glog.Errorf("Packet body Size Too Large, size: %d\n", bodyLen)
		return
	}

	if bodyLen < codec.MinPacketSize {
		glog.Errorf("Packet Body Size Invalid, size: %d\n", bodyLen)
		err = ErrorkrPacketInvalidBody
		return
	}

	buffer.Prepare(bodyLen)
	if _, err = io.ReadFull(conn.GetRawConn(), buffer.Data); err != nil {
		glog.Errorf("read packet body failed, err: %v\n", err)
		return
	}

	glog.Errorf("recv msg: %v\n", buffer.Data)
	return
}

func (codec FixedLengthHeaderCodec) Encode(msg Message) ([]byte, error) {
	body, err := msg.MarshalBinary()
	if err != nil {
		return nil, err
	}

	bodyLen := uint32(len(body) + 4)
	major := msg.Major()
	minor := msg.Minor()

	packet := make([]byte, codec.HeaderSize+4+len(body))
	packet = packet[0:0]
	packet = append(packet, byte(bodyLen>>24), byte(bodyLen>>16), byte(bodyLen>>8), byte(bodyLen))
	packet = append(packet, byte(major>>8), byte(major))
	packet = append(packet, byte(minor>>8), byte(minor))
	packet = append(packet, body...)
	return packet, nil
}

func (codec FixedLengthHeaderCodec) decodeHead(header []byte) int {
	if codec.HeaderSize == 2 {
		return int(codec.ByteOrder.Uint16(header))
	}

	if codec.HeaderSize == 4 {
		return int(codec.ByteOrder.Uint32(header))
	}

	return 0
}
