package broker_test

import (
    "fmt"
    "time"
    "testing"

    sessproto "smartgo/proto/session_event"
    loginproto "smartgo/proto/login_event"
    "smartgo/libs/utils"
    "smartgo/libs/socket"
    "smartgo/libs/socket/example"
)

const (
    //错误码
    RES_OK                  = 0
    RES_FAIL                = 1
    RES_NO_USER             = 101
    RES_INVALID_PASSWORD    = 102
    RES_INVALID_TOKEN       = 103

    //服务主ID
    SVID_LOGIN int32 = 1

    //login服务子ID
    PLoginByPassport_uri int32     = (SVID_LOGIN << 16 | 1)
    PLoginByPassportRes_uri int32  = (SVID_LOGIN << 16 | 2)
    PLoginByUid_uri int32          = (SVID_LOGIN << 16 | 3)
    PLoginByUidRes_uri int32       = (SVID_LOGIN << 16 | 4)
    PLoginLogout_uri int32         = (SVID_LOGIN << 16 | 5)
    PLoginLogoutRes_uri int32      = (SVID_LOGIN << 16 | 6)
    PLoginPing_uri int32           = (SVID_LOGIN << 16 | 7)
    PLoginPingRes_uri int32        = (SVID_LOGIN << 16 | 8)
    PLoginKickOff_uri int32        = (SVID_LOGIN << 16 | 9)
)

var signal *test.SignalTester

var ip_port = "123.56.88.196:9100"
//var ip_port = "127.0.0.1:9100"

/*
 * 正常连接并登录(用户名密码方式)，然后ping三次后关掉连接
 */
func loginByPassportAndPing() {
    queue := socket.NewEventQueue()
    connector := socket.NewConnector(queue).Start(ip_port)

    //ping三次就好了
    ping_times := 0

    //session message
    socket.RegisterMessage(connector, "session_event.SessionAccepted", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionAccepted)
        if !ok {
            fmt.Println("Client: recv invalid SessionAccepted message")
            return
        }
        fmt.Println("Client: recv SessionAccepted message")
    })

    socket.RegisterMessage(connector, "session_event.SessionAcceptFailed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionAcceptFailed)
        if !ok {
            fmt.Println("Client: recv invalid SessionAcceptFailed message")
            return
        }
        fmt.Println("Client: recv SessionAcceptFailed message")
    })

    socket.RegisterMessage(connector, "session_event.SessionConnected", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionConnected)
        if !ok {
            fmt.Println("Client: recv invalid SessionConnected message")
            return
        }
        signal.Done(1)
        fmt.Println("Client: recv SessionConnected message")
        fmt.Println("Client: send PLoginByPassport message")
        connector.Session().SetParam("connected", true)
        connector.Session().Send(&loginproto.PLoginByPassport{
            Uri:PLoginByPassport_uri,
            Passport:"liujia",
            Password:"123456",
            Deviceid:"abcdefg",
            Devicetype:"test",
            Params:map[string]string{"111":"222"},
        })
    })

    socket.RegisterMessage(connector, "session_event.SessionConnectFailed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionConnectFailed)
        if !ok {
            fmt.Println("Client: recv invalid SessionConnectFailed message")
            return
        }
        fmt.Println("Client: recv SessionConnectFailed message")
    })

    socket.RegisterMessage(connector, "session_event.SessionError", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionError)
        if !ok {
            fmt.Println("Client: recv invalid SessionError message")
            return
        }
        fmt.Println("Client: recv SessionError message")
    })

    socket.RegisterMessage(connector, "session_event.SessionClosed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionClosed)
        if !ok {
            fmt.Println("Client: recv invalid SessionClosed message")
            return
        }
        fmt.Println("Client: recv SessionClosed message")

        timer, ok := connector.Session().GetParam("pingtimer").(*socket.Timer)
        if ok {
            fmt.Println("Client: stop timer")
            timer.Stop()
        }

        signal.Done(3)
    })

    //login message
    socket.RegisterMessage(connector, "login_event.PLoginByPassportRes", func(content interface{}, session socket.Session) {
        msg, ok := content.(*loginproto.PLoginByPassportRes)
        if !ok {
            fmt.Println("Client: recv invalid PLoginByPassportRes message")
            return
        }
        if msg.Rescode == RES_OK {
            signal.Done(2)
            fmt.Println("Client: login success, uid/token", msg.Uid, msg.Token)
            uid := msg.Uid
            token := msg.Token
            connector.Session().SetParam("uid", uid)
            connector.Session().SetParam("token", token)
            timer := socket.NewTimer(queue, 1*time.Second, func(t *socket.Timer) {
                fmt.Println("Client: timed to send ping", utils.GoID())
                connector.Session().Send(&loginproto.PLoginPing{
                    Uri:PLoginPing_uri,
                    Uid:uid,
                    Clientts:time.Now().UnixNano()/1000000,
                })
            })
            connector.Session().SetParam("pingtimer", timer)
        }else {
            fmt.Println("Client: login failed", msg.Rescode)
            connector.Session().SetParam("uid", "")
            connector.Session().SetParam("token", "")
            v := connector.Session().GetParam("pingtimer")
            if timer, ok := v.(*socket.Timer);ok {
                timer.Stop()
            }
        }
    })

    socket.RegisterMessage(connector, "login_event.PLoginPingRes", func(content interface{}, session socket.Session) {
        msg, ok := content.(*loginproto.PLoginPingRes)
        if !ok {
            fmt.Println("Client: recv invalid PLoginByPassportRes message")
            return
        }
        fmt.Println(msg)

        ping_times++
        if ping_times == 3 {
            connector.Stop()
        }
    })

    queue.StartLoop()

    signal.WaitAndExpect(1, "TestConnActiveClose not connected")
    signal.WaitAndExpect(2, "TestConnActiveClose not logined")
    signal.WaitAndExpect(3, "TestConnActiveClose not close")

    //queue.StopLoop(-1)
}

/*
 * 正常连接并登录(uid和token方式)，然后ping三次后关掉连接
 */
func loginByTokenAndPing() {
    queue := socket.NewEventQueue()
    connector := socket.NewConnector(queue).Start(ip_port)

    //ping三次就好了
    ping_times := 0

    //session message
    socket.RegisterMessage(connector, "session_event.SessionAccepted", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionAccepted)
        if !ok {
            fmt.Println("Client: recv invalid SessionAccepted message")
            return
        }
        fmt.Println("Client: recv SessionAccepted message")
    })

    socket.RegisterMessage(connector, "session_event.SessionAcceptFailed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionAcceptFailed)
        if !ok {
            fmt.Println("Client: recv invalid SessionAcceptFailed message")
            return
        }
        fmt.Println("Client: recv SessionAcceptFailed message")
    })

    socket.RegisterMessage(connector, "session_event.SessionConnected", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionConnected)
        if !ok {
            fmt.Println("Client: recv invalid SessionConnected message")
            return
        }
        signal.Done(1)
        fmt.Println("Client: recv SessionConnected message")
        fmt.Println("Client: send PLoginByToken message")
        connector.Session().SetParam("connected", true)
        connector.Session().Send(&loginproto.PLoginByToken{
            Uri:PLoginByPassport_uri,
            Uid:"1",
            Token:"hello_liujia",
            Deviceid:"abcdefg",
            Devicetype:"test",
            Params:map[string]string{"111":"222"},
        })
    })

    socket.RegisterMessage(connector, "session_event.SessionConnectFailed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionConnectFailed)
        if !ok {
            fmt.Println("Client: recv invalid SessionConnectFailed message")
            return
        }
        fmt.Println("Client: recv SessionConnectFailed message")
    })

    socket.RegisterMessage(connector, "session_event.SessionError", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionError)
        if !ok {
            fmt.Println("Client: recv invalid SessionError message")
            return
        }
        fmt.Println("Client: recv SessionError message")
    })

    socket.RegisterMessage(connector, "session_event.SessionClosed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionClosed)
        if !ok {
            fmt.Println("Client: recv invalid SessionClosed message")
            return
        }
        fmt.Println("Client: recv SessionClosed message")

        timer, ok := connector.Session().GetParam("pingtimer").(*socket.Timer)
        if ok {
            fmt.Println("Client: stop timer")
            timer.Stop()
        }

        signal.Done(3)
    })

    //login message
    socket.RegisterMessage(connector, "login_event.PLoginByPassportRes", func(content interface{}, session socket.Session) {
        msg, ok := content.(*loginproto.PLoginByPassportRes)
        if !ok {
            fmt.Println("Client: recv invalid PLoginByPassportRes message")
            return
        }
        if msg.Rescode == RES_OK {
            signal.Done(2)
            fmt.Println("Client: login by passport success, uid/token", msg.Uid, msg.Token)
            uid := msg.Uid
            token := msg.Token
            connector.Session().SetParam("uid", uid)
            connector.Session().SetParam("token", token)
            timer := socket.NewTimer(queue, 1*time.Second, func(t *socket.Timer) {
                fmt.Println("Client: timed to send ping", utils.GoID())
                connector.Session().Send(&loginproto.PLoginPing{
                    Uri:PLoginPing_uri,
                    Uid:uid,
                    Clientts:time.Now().UnixNano()/1000000,
                })
            })
            connector.Session().SetParam("pingtimer", timer)
        }else {
            fmt.Println("Client: login failed", msg.Rescode)
            connector.Session().SetParam("uid", "")
            connector.Session().SetParam("token", "")
            v := connector.Session().GetParam("pingtimer")
            if timer, ok := v.(*socket.Timer);ok {
                timer.Stop()
            }
        }
    })

    socket.RegisterMessage(connector, "login_event.PLoginByTokenRes", func(content interface{}, session socket.Session) {
        msg, ok := content.(*loginproto.PLoginByTokenRes)
        if !ok {
            fmt.Println("Client: recv invalid PLoginByTokenRes message")
            return
        }
        if msg.Rescode == RES_OK {
            signal.Done(2)
            fmt.Println("Client: login by token success, uid/token", msg.Uid, msg.Token)
            uid := msg.Uid
            token := msg.Token
            connector.Session().SetParam("uid", uid)
            connector.Session().SetParam("token", token)
            timer := socket.NewTimer(queue, 1*time.Second, func(t *socket.Timer) {
                fmt.Println("Client: timed to send ping", utils.GoID())
                connector.Session().Send(&loginproto.PLoginPing{
                    Uri:PLoginPing_uri,
                    Uid:uid,
                    Clientts:time.Now().UnixNano()/1000000,
                })
            })
            connector.Session().SetParam("pingtimer", timer)
        }else {
            fmt.Println("Client: login failed", msg.Rescode)
            connector.Session().SetParam("uid", "")
            connector.Session().SetParam("token", "")
            v := connector.Session().GetParam("pingtimer")
            if timer, ok := v.(*socket.Timer);ok {
                timer.Stop()
            }
        }
    })

    socket.RegisterMessage(connector, "login_event.PLoginPingRes", func(content interface{}, session socket.Session) {
        msg, ok := content.(*loginproto.PLoginPingRes)
        if !ok {
            fmt.Println("Client: recv invalid PLoginByPassportRes message")
            return
        }
        fmt.Println(msg)

        ping_times++
        if ping_times == 3 {
            connector.Stop()
        }
    })

    queue.StartLoop()

    signal.WaitAndExpect(1, "TestConnActiveClose not connected")
    signal.WaitAndExpect(2, "TestConnActiveClose not logined")
    signal.WaitAndExpect(3, "TestConnActiveClose not close")

    //queue.StopLoop(-1)
}

func noLoginAfterConnected() {
    queue := socket.NewEventQueue()
    connector := socket.NewConnector(queue).Start(ip_port)

    //ping三次就好了
    ping_times := 0

    //session message
    socket.RegisterMessage(connector, "session_event.SessionAccepted", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionAccepted)
        if !ok {
            fmt.Println("Client: recv invalid SessionAccepted message")
            return
        }
        fmt.Println("Client: recv SessionAccepted message")
    })

    socket.RegisterMessage(connector, "session_event.SessionAcceptFailed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionAcceptFailed)
        if !ok {
            fmt.Println("Client: recv invalid SessionAcceptFailed message")
            return
        }
        fmt.Println("Client: recv SessionAcceptFailed message")
    })

    socket.RegisterMessage(connector, "session_event.SessionConnected", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionConnected)
        if !ok {
            fmt.Println("Client: recv invalid SessionConnected message")
            return
        }
        signal.Done(1)
        fmt.Println("Client: recv SessionConnected message")

        /*
        fmt.Println("Client: send PLoginByPassport message")
        connector.Session().SetParam("connected", true)
        connector.Session().Send(&loginproto.PLoginByPassport{
            Uri:PLoginByPassport_uri,
            Passport:"liujia",
            Password:"123456",
            Deviceid:"abcdefg",
            Devicetype:"test",
            Params:map[string]string{"111":"222"},
        })
        */
    })

    socket.RegisterMessage(connector, "session_event.SessionConnectFailed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionConnectFailed)
        if !ok {
            fmt.Println("Client: recv invalid SessionConnectFailed message")
            return
        }
        fmt.Println("Client: recv SessionConnectFailed message")
    })

    socket.RegisterMessage(connector, "session_event.SessionError", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionError)
        if !ok {
            fmt.Println("Client: recv invalid SessionError message")
            return
        }
        fmt.Println("Client: recv SessionError message")
    })

    socket.RegisterMessage(connector, "session_event.SessionClosed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionClosed)
        if !ok {
            fmt.Println("Client: recv invalid SessionClosed message")
            return
        }
        fmt.Println("Client: recv SessionClosed message")
        signal.Done(2)
    })

    //login message
    socket.RegisterMessage(connector, "login_event.PLoginByPassportRes", func(content interface{}, session socket.Session) {
        msg, ok := content.(*loginproto.PLoginByPassportRes)
        if !ok {
            fmt.Println("Client: recv invalid PLoginByPassportRes message")
            return
        }
        if msg.Rescode == RES_OK {
            //signal.Done(2)
            fmt.Println("Client: login success, uid/token", msg.Uid, msg.Token)
            uid := msg.Uid
            token := msg.Token
            connector.Session().SetParam("uid", uid)
            connector.Session().SetParam("token", token)
            timer := socket.NewTimer(queue, 3*time.Second, func(t *socket.Timer) {
                fmt.Println("Client: timed to send ping", utils.GoID())
                connector.Session().Send(&loginproto.PLoginPing{
                    Uri:PLoginPing_uri,
                    Uid:uid,
                    Clientts:time.Now().UnixNano()/1000000,
                })
            })
            connector.Session().SetParam("pingtimer", timer)
        }else {
            fmt.Println("Client: login failed", msg.Rescode)
            connector.Session().SetParam("uid", "")
            connector.Session().SetParam("token", "")
            v := connector.Session().GetParam("pingtimer")
            if timer, ok := v.(*socket.Timer);ok {
                timer.Stop()
            }
        }
    })

    socket.RegisterMessage(connector, "login_event.PLoginPingRes", func(content interface{}, session socket.Session) {
        msg, ok := content.(*loginproto.PLoginPingRes)
        if !ok {
            fmt.Println("Client: recv invalid PLoginByPassportRes message")
            return
        }
        fmt.Println(msg)

        ping_times++
        if ping_times == 3 {
            //signal.Done(3)
            connector.Stop()
        }
    })

    queue.StartLoop()

    signal.WaitAndExpect(1, "TestConnActiveClose not connected")
    signal.WaitAndExpect(2, "TestConnActiveClose not closed")
}

func failedToPing() {
    queue := socket.NewEventQueue()
    connector := socket.NewConnector(queue).Start(ip_port)

    //ping三次就好了
    ping_times := 0

    //session message
    socket.RegisterMessage(connector, "session_event.SessionAccepted", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionAccepted)
        if !ok {
            fmt.Println("Client: recv invalid SessionAccepted message")
            return
        }
        fmt.Println("Client: recv SessionAccepted message")
    })

    socket.RegisterMessage(connector, "session_event.SessionAcceptFailed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionAcceptFailed)
        if !ok {
            fmt.Println("Client: recv invalid SessionAcceptFailed message")
            return
        }
        fmt.Println("Client: recv SessionAcceptFailed message")
    })

    socket.RegisterMessage(connector, "session_event.SessionConnected", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionConnected)
        if !ok {
            fmt.Println("Client: recv invalid SessionConnected message")
            return
        }
        signal.Done(1)
        fmt.Println("Client: recv SessionConnected message")

        fmt.Println("Client: send PLoginByPassport message")
        connector.Session().SetParam("connected", true)
        connector.Session().Send(&loginproto.PLoginByPassport{
            Uri:PLoginByPassport_uri,
            Passport:"liujia",
            Password:"123456",
            Deviceid:"abcdefg",
            Devicetype:"test",
            Params:map[string]string{"111":"222"},
        })
    })

    socket.RegisterMessage(connector, "session_event.SessionConnectFailed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionConnectFailed)
        if !ok {
            fmt.Println("Client: recv invalid SessionConnectFailed message")
            return
        }
        fmt.Println("Client: recv SessionConnectFailed message")
    })

    socket.RegisterMessage(connector, "session_event.SessionError", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionError)
        if !ok {
            fmt.Println("Client: recv invalid SessionError message")
            return
        }
        fmt.Println("Client: recv SessionError message")
    })

    socket.RegisterMessage(connector, "session_event.SessionClosed", func(content interface{}, session socket.Session) {
        _, ok := content.(*sessproto.SessionClosed)
        if !ok {
            fmt.Println("Client: recv invalid SessionClosed message")
            return
        }
        fmt.Println("Client: recv SessionClosed message")
        signal.Done(3)
    })

    //login message
    socket.RegisterMessage(connector, "login_event.PLoginByPassportRes", func(content interface{}, session socket.Session) {
        msg, ok := content.(*loginproto.PLoginByPassportRes)
        if !ok {
            fmt.Println("Client: recv invalid PLoginByPassportRes message")
            return
        }
        if msg.Rescode == RES_OK {
            signal.Done(2)
            fmt.Println("Client: login success, uid/token", msg.Uid, msg.Token)
            uid := msg.Uid
            token := msg.Token
            connector.Session().SetParam("uid", uid)
            connector.Session().SetParam("token", token)
            /*
            timer := socket.NewTimer(queue, 3*time.Second, func(t *socket.Timer) {
                fmt.Println("Client: timed to send ping", utils.GoID())
                connector.Session().Send(&loginproto.PLoginPing{
                    Uri:PLoginPing_uri,
                    Uid:uid,
                    Clientts:time.Now().UnixNano()/1000000,
                })
            })
            connector.Session().SetParam("pingtimer", timer)
            */
        }else {
            fmt.Println("Client: login failed", msg.Rescode)
            connector.Session().SetParam("uid", "")
            connector.Session().SetParam("token", "")
            v := connector.Session().GetParam("pingtimer")
            if timer, ok := v.(*socket.Timer);ok {
                timer.Stop()
            }
        }
    })

    socket.RegisterMessage(connector, "login_event.PLoginPingRes", func(content interface{}, session socket.Session) {
        msg, ok := content.(*loginproto.PLoginPingRes)
        if !ok {
            fmt.Println("Client: recv invalid PLoginByPassportRes message")
            return
        }
        fmt.Println(msg)

        ping_times++
        if ping_times == 3 {
            //signal.Done(3)
            connector.Stop()
        }
    })

    queue.StartLoop()

    signal.WaitAndExpect(1, "TestConnActiveClose not connected")
    signal.WaitAndExpect(2, "TestConnActiveClose not logined")
    signal.WaitAndExpect(3, "TestConnActiveClose not closed")
}

func sessProtoRegisterMessage() {
    fmt.Println("session_event RegisterMessage")
    // session.proto
    socket.RegisterMessageMeta("session_event.SessionAccepted", (*sessproto.SessionAccepted)(nil), 348117910)
    socket.RegisterMessageMeta("session_event.SessionAcceptFailed", (*sessproto.SessionAcceptFailed)(nil), 1978788392)
    socket.RegisterMessageMeta("session_event.SessionConnected", (*sessproto.SessionConnected)(nil), 3543838007)
    socket.RegisterMessageMeta("session_event.SessionConnectFailed", (*sessproto.SessionConnectFailed)(nil), 1720533237)
    socket.RegisterMessageMeta("session_event.SessionClosed", (*sessproto.SessionClosed)(nil), 90181607)
    socket.RegisterMessageMeta("session_event.SessionError", (*sessproto.SessionError)(nil), 1937281175)
}

func loginProtoRegisterMessage() {
    fmt.Println("login_event RegisterMessage")
    // login_event.proto
    socket.RegisterMessageMeta("login_event.PLoginByPassport", (*loginproto.PLoginByPassport)(nil), 3176521479)
    socket.RegisterMessageMeta("login_event.PLoginByPassportRes", (*loginproto.PLoginByPassportRes)(nil), 1894287111)
    socket.RegisterMessageMeta("login_event.PLoginByToken", (*loginproto.PLoginByToken)(nil), 3676221678)
    socket.RegisterMessageMeta("login_event.PLoginByTokenRes", (*loginproto.PLoginByTokenRes)(nil), 2978710459)
    socket.RegisterMessageMeta("login_event.PLoginLogout", (*loginproto.PLoginLogout)(nil), 182945212)
    socket.RegisterMessageMeta("login_event.PLoginLogoutRes", (*loginproto.PLoginLogoutRes)(nil), 3420227263)
    socket.RegisterMessageMeta("login_event.PLoginPing", (*loginproto.PLoginPing)(nil), 1951739067)
    socket.RegisterMessageMeta("login_event.PLoginPingRes", (*loginproto.PLoginPingRes)(nil), 3845948673)
    socket.RegisterMessageMeta("login_event.PLoginKickOff", (*loginproto.PLoginKickOff)(nil), 4292429949)
}

func TestLoginByPassportAndPing(t *testing.T) {
    signal = test.NewSignalTesterTimeout(t, 10)
    sessProtoRegisterMessage()
    loginProtoRegisterMessage()
    socket.Init()
    loginByPassportAndPing()
}

func TestLoginByTokenAndPing(t *testing.T) {
    signal = test.NewSignalTesterTimeout(t, 10)
    sessProtoRegisterMessage()
    loginProtoRegisterMessage()
    socket.Init()
    loginByTokenAndPing()
}

func TestConnectedNotLogin(t *testing.T) {
    signal = test.NewSignalTesterTimeout(t, 10)
    sessProtoRegisterMessage()
    loginProtoRegisterMessage()
    socket.Init()
    noLoginAfterConnected()
}

func TestFailedToPing(t *testing.T) {
    signal = test.NewSignalTesterTimeout(t, 10)
    sessProtoRegisterMessage()
    loginProtoRegisterMessage()
    socket.Init()
    failedToPing()
}