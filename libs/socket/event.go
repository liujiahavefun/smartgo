package socket

import (
	"fmt"

	_ "smartgo/proto/sessevent"
)

var (
	Event_SessionAccepted      = uint32(MessageMetaByName("session.SessionAccepted").ID)
	Event_SessionAcceptFailed  = uint32(MessageMetaByName("session.SessionAcceptFailed").ID)
	Event_SessionConnected     = uint32(MessageMetaByName("session.SessionConnected").ID)
	Event_SessionConnectFailed = uint32(MessageMetaByName("session.SessionConnectFailed").ID)
	Event_SessionClosed        = uint32(MessageMetaByName("session.SessionClosed").ID)
	Event_SessionError         = uint32(MessageMetaByName("session.SessionError").ID)
)

//会话事件
type SessionEvent struct {
	*Packet
	Ses Session
}

func (self SessionEvent) String() string {
	return fmt.Sprintf("SessionEvent msgid: %d data: %v", self.MsgID, self.Data)
}

func NewSessionEvent(msgid uint32, s Session, data []byte) *SessionEvent {
	return &SessionEvent{
		Packet: &Packet{MsgID: msgid, Data: data},
		Ses:    s,
	}
}

func newSessionEvent(msgid uint32, s Session, msg interface{}) *SessionEvent {
	pkt, _ := BuildPacket(msg)
	return &SessionEvent{
		Packet: pkt,
		Ses:    s,
	}
}
