package proto

import (
    "smartgo/libs/socket"
    sessproto "smartgo/proto/sessevent"
    testproto "smartgo/proto/test"
)

func init() {
    // session.proto
    socket.RegisterMessageMeta("session.SessionAccepted", (*sessproto.SessionAccepted)(nil), 2136350511)
    socket.RegisterMessageMeta("session.SessionAcceptFailed", (*sessproto.SessionAcceptFailed)(nil), 1213847952)
    socket.RegisterMessageMeta("session.SessionConnected", (*sessproto.SessionConnected)(nil), 4228538224)
    socket.RegisterMessageMeta("session.SessionConnectFailed", (*sessproto.SessionConnectFailed)(nil), 1278926828)
    socket.RegisterMessageMeta("session.SessionClosed", (*sessproto.SessionClosed)(nil), 2830250790)
    socket.RegisterMessageMeta("session.SessionError", (*sessproto.SessionError)(nil), 3227768243)

    // test.proto
    socket.RegisterMessageMeta("test.TestEchoACK", (*testproto.TestEchoACK)(nil), 509149489)
}