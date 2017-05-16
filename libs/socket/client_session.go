package socket

import (
	"net"
)

type tcpClientSession struct {
	*sessionBase
}

func newClientSession(conn net.Conn, connector Connector, callbacks *SessionCallback) *tcpClientSession {
	otherPeer := newPeerBase()
	otherPeer.SetName(conn.RemoteAddr().String())
	otherPeer.SetMaxPacketSize(connector.MaxPacketSize())

	session := &tcpClientSession{
		sessionBase: newSessionBase(NewPacketStream(conn), connector, otherPeer, callbacks),
	}

	return session
}

