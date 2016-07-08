package net

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrorkrPacketInvalidHeader error = errors.New("Packet header invald")
	ErrorkrPacketTooLarge      error = errors.New("Packet body too large")
	ErrorkrPacketInvalidBody   error = errors.New("Packet invalid body")
	ErrorParameter             error = errors.New("Parameter error")
	ErrorNilData               error = errors.New("Nil data")
	ErrorIllegalData           error = errors.New("More than 8M data")
	ErrorNotImplemented        error = errors.New("Not implemented")
	ErrorConnClosed            error = errors.New("Connection closed")
)

const (
	HEADER_SIZE     = 4
	MIN_PACKET_SIZE = 4
	MAX_PACKET_SIZE = 4 * 1024
	WORKERS         = 512
	MAX_CONNECTIONS = 1000000
)

var (
	GLOBAL_BINARY_BYTE_ORDER = binary.BigEndian
)

func Undefined(msgType int32) error {
	return ErrorUndefined{
		msgType: msgType,
	}
}

type ErrorUndefined struct {
	msgType int32
}

func (eu ErrorUndefined) Error() string {
	return fmt.Sprintf("Undefined message %d", eu.msgType)
}
