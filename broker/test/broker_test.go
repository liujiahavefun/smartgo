package broker_test

import (
    "fmt"
    "time"
    "testing"
    sessproto "smartgo/proto/session_event"
    loginproto "smartgo/proto/login_event"
    testproto "smartgo/proto/test_event"
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
    SVID_SESION int32 = 1
    SVID_LOGIN  int32 = 16
    SVID_TEST  int32 = 2

    //session服务子ID，这个仅内部使用
    PSessionAccepted_uri int32          = (SVID_SESION << 16 | 1)
    PSessionAcceptedFailed_uri int32    = (SVID_SESION << 16 | 2)
    PSessionConnected_uri int32         = (SVID_SESION << 16 | 3)
    PSessionConnectedFailed_uri int32   = (SVID_SESION << 16 | 4)
    PSessionClosed_uri int32            = (SVID_SESION << 16 | 5)
    PSessionError_uri int32             = (SVID_SESION << 16 | 6)

    //login服务子ID
    PLoginByPassport_uri int32      = (SVID_LOGIN << 16 | 1)
    PLoginByPassportRes_uri int32   = (SVID_LOGIN << 16 | 2)
    PLoginByUid_uri int32           = (SVID_LOGIN << 16 | 3)
    PLoginByUidRes_uri int32        = (SVID_LOGIN << 16 | 4)
    PLoginLogout_uri int32          = (SVID_LOGIN << 16 | 5)
    PLoginLogoutRes_uri int32       = (SVID_LOGIN << 16 | 6)
    PLoginPing_uri int32            = (SVID_LOGIN << 16 | 7)
    PLoginPingRes_uri int32         = (SVID_LOGIN << 16 | 8)
    PLoginKickOff_uri int32         = (SVID_LOGIN << 16 | 9)

    //test。。。
    PTestEchoACK_uri int32      = (SVID_TEST << 16 | 1)
)

var signal *test.SignalTester

//var ip_port = "123.56.88.196:9100"
var ip_port = "127.0.0.1:9100"

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

func testDefaultHandler() {
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
            Uri:PLoginByUid_uri,
            Uid:"1",
            Token:"hello_liujia",
            Deviceid:"abcdefg",
            Devicetype:"test",
            Params:map[string]string{"111":"222"},
        })

        connector.Session().Send(&testproto.TestEchoACK{
            Content:"DEFAULT",
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
}

func sessProtoRegisterMessage() {
    fmt.Println("session_event RegisterMessage")
    // session.proto
    socket.RegisterMessageMeta("session_event.SessionAccepted", (*sessproto.SessionAccepted)(nil), uint32(PSessionAccepted_uri))
    socket.RegisterMessageMeta("session_event.SessionAcceptFailed", (*sessproto.SessionAcceptFailed)(nil), uint32(PSessionAcceptedFailed_uri))
    socket.RegisterMessageMeta("session_event.SessionConnected", (*sessproto.SessionConnected)(nil), uint32(PSessionConnected_uri))
    socket.RegisterMessageMeta("session_event.SessionConnectFailed", (*sessproto.SessionConnectFailed)(nil), uint32(PSessionConnectedFailed_uri))
    socket.RegisterMessageMeta("session_event.SessionClosed", (*sessproto.SessionClosed)(nil), uint32(PSessionClosed_uri))
    socket.RegisterMessageMeta("session_event.SessionError", (*sessproto.SessionError)(nil), uint32(PSessionError_uri))
}

func loginProtoRegisterMessage() {
    fmt.Println("login_event RegisterMessage")
    // login_event.proto
    socket.RegisterMessageMeta("login_event.PLoginByPassport", (*loginproto.PLoginByPassport)(nil), uint32(PLoginByPassport_uri))
    socket.RegisterMessageMeta("login_event.PLoginByPassportRes", (*loginproto.PLoginByPassportRes)(nil), uint32(PLoginByPassportRes_uri))
    socket.RegisterMessageMeta("login_event.PLoginByToken", (*loginproto.PLoginByToken)(nil), uint32(PLoginByUid_uri))
    socket.RegisterMessageMeta("login_event.PLoginByTokenRes", (*loginproto.PLoginByTokenRes)(nil), uint32(PLoginByUidRes_uri))
    socket.RegisterMessageMeta("login_event.PLoginLogout", (*loginproto.PLoginLogout)(nil), uint32(PLoginLogout_uri))
    socket.RegisterMessageMeta("login_event.PLoginLogoutRes", (*loginproto.PLoginLogoutRes)(nil), uint32(PLoginLogoutRes_uri))
    socket.RegisterMessageMeta("login_event.PLoginPing", (*loginproto.PLoginPing)(nil), uint32(PLoginPing_uri))
    socket.RegisterMessageMeta("login_event.PLoginPingRes", (*loginproto.PLoginPingRes)(nil), uint32(PLoginPingRes_uri))
    socket.RegisterMessageMeta("login_event.PLoginKickOff", (*loginproto.PLoginKickOff)(nil), uint32(PLoginKickOff_uri))
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


func TestDefaultHandler(t *testing.T) {
    signal = test.NewSignalTesterTimeout(t, 10)
    sessProtoRegisterMessage()
    loginProtoRegisterMessage()
    socket.Init()

    //发送方必须注册此消息
    socket.RegisterMessageMeta("test_event.TestEchoACK", (*testproto.TestEchoACK)(nil), uint32(PTestEchoACK_uri))

    testDefaultHandler()
}