package proto

import (
    "smartgo/libs/socket"
    sessproto "smartgo/proto/session_event"
    testproto "smartgo/proto/test_event"
)

func init() {
    // session.proto
    socket.RegisterMessageMeta("session_event.SessionAccepted", (*sessproto.SessionAccepted)(nil), 2136350511)
    socket.RegisterMessageMeta("session_event.SessionAcceptFailed", (*sessproto.SessionAcceptFailed)(nil), 1213847952)
    socket.RegisterMessageMeta("session_event.SessionConnected", (*sessproto.SessionConnected)(nil), 4228538224)
    socket.RegisterMessageMeta("session_event.SessionConnectFailed", (*sessproto.SessionConnectFailed)(nil), 1278926828)
    socket.RegisterMessageMeta("session_event.SessionClosed", (*sessproto.SessionClosed)(nil), 2830250790)
    socket.RegisterMessageMeta("session_event.SessionError", (*sessproto.SessionError)(nil), 3227768243)

    // test.proto
    socket.RegisterMessageMeta("test_event.TestEchoACK", (*testproto.TestEchoACK)(nil), 509149489)
}