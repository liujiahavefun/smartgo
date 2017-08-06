package socket

import (
	"fmt"
)

func Init() {
	fmt.Println("event Init()")
	Event_SessionAccepted      = uint32(MessageMetaByName("session_event.SessionAccepted").ID)
	Event_SessionAcceptFailed  = uint32(MessageMetaByName("session_event.SessionAcceptFailed").ID)
	Event_SessionConnected     = uint32(MessageMetaByName("session_event.SessionConnected").ID)
	Event_SessionConnectFailed = uint32(MessageMetaByName("session_event.SessionConnectFailed").ID)
	Event_SessionClosed        = uint32(MessageMetaByName("session_event.SessionClosed").ID)
	Event_SessionError         = uint32(MessageMetaByName("session_event.SessionError").ID)
}

/*
var (
	Event_SessionAccepted      = uint32(MessageMetaByName("session.SessionAccepted").ID)
	Event_SessionAcceptFailed  = uint32(MessageMetaByName("session.SessionAcceptFailed").ID)
	Event_SessionConnected     = uint32(MessageMetaByName("session.SessionConnected").ID)
	Event_SessionConnectFailed = uint32(MessageMetaByName("session.SessionConnectFailed").ID)
	Event_SessionClosed        = uint32(MessageMetaByName("session.SessionClosed").ID)
	Event_SessionError         = uint32(MessageMetaByName("session.SessionError").ID)
)
*/

var (
	Event_SessionAccepted uint32
	Event_SessionAcceptFailed uint32
	Event_SessionConnected uint32
	Event_SessionConnectFailed uint32
	Event_SessionClosed uint32
	Event_SessionError uint32
)

//会话事件
type SessionEvent struct {
	*Packet
	Ses Session
}

func (self SessionEvent) String() string {
	return fmt.Sprintf("SessionEvent msgid: %d data: %v", self.MsgID, self.Data)
}

func NewSessionEvent(msgid uint32, sess Session, data []byte) *SessionEvent {
	return &SessionEvent{
		Packet: &Packet{MsgID: msgid, Data: data},
		Ses:    sess,
	}
}

func newSessionEvent(msgid uint32, sess Session, msg interface{}) *SessionEvent {
	pkt, _ := BuildPacket(msg)
	return &SessionEvent{
		Packet: pkt,
		Ses:    sess,
	}
}
