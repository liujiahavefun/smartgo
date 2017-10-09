package socket

import (
	"fmt"
)

/*
 * session_event虽然是外部定义的proto，但是仅内部使用，这里提前获得一下消息id，为了提高效率。
 * 对外部消息，则无必要。因为对内部消息，发送的时候我自己肯定知道发的是啥消息。对外部消息，我收的时候是不知道是啥的，需要解包才知道。
 * 只有内部消息外部消息都用proto，纯粹是哥高瞻远瞩，为了统一内部消息处理和外部消息处理
 */
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
