package net

import (
	"encoding"
)

type Message interface {
	Major() uint16
	Minor() uint16
	encoding.BinaryMarshaler
}

type MessageHandler interface {
	Process(client Connection) bool
}
