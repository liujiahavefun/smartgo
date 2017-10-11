package main

import (
    loginproto "smartgo/proto/login_event"
    "smartgo/libs/socket"
    "smartgo/libs/utils"
)

func handleLoginByPassport(msg *loginproto.PLoginByPassport, session socket.Session)  {
    logInfo("Server: recv PLoginByPassport message, enter")

    passport := msg.Passport
    password := msg.Password
    params := map[string]string{
        "deviceid": msg.Deviceid,
        "devicetype": msg.Devicetype,
    }
    uid, token, err := loginByPassport(passport, password, params)
    if err == nil {
        gSessionMgr.onSessionLogin(session, uid, token)
        session.Send(&loginproto.PLoginByPassportRes{
            Uri:    PLoginByPassportRes_uri,
            Rescode:RES_OK,
            Uid:    uid,
            Token:  token,
        })
    }else {
        session.Send(&loginproto.PLoginByPassportRes{
            Uri:    PLoginByPassportRes_uri,
            Rescode:RES_INVALID_PASSWORD,
            Uid:    "",
            Token:  "",
        })
    }

    logInfo("Server: recv PLoginByPassport message, leave")
}

func handleLoginByToken(msg *loginproto.PLoginByToken, session socket.Session)  {
    logInfo("Server: recv PLoginByToken message, enter")

    uid := msg.Uid
    token := msg.Token
    params := map[string]string{
        "deviceid": msg.Deviceid,
        "devicetype": msg.Devicetype,
    }
    err := loginByToken(uid, token, params)
    if err == nil {
        gSessionMgr.onSessionLogin(session, uid, token)
        session.Send(&loginproto.PLoginByTokenRes{
            Uri:    PLoginByUidRes_uri,
            Rescode:RES_OK,
            Uid:    uid,
            Token:  token,
        })
    }else {
        session.Send(&loginproto.PLoginByPassportRes{
            Uri:    PLoginByUidRes_uri,
            Rescode:RES_INVALID_TOKEN,
            Uid:    "",
            Token:  "",
        })
    }

    logInfo("Server: recv PLoginByToken message, leave")
}

func handleLoginLogout(msg *loginproto.PLoginLogout, session socket.Session)  {
    logInfo("Server: recv PLoginLogout message, enter")

    uid := msg.Uid
    token := msg.Token

    err := logout(uid, token)
    if err == nil {
        session.Send(&loginproto.PLoginLogoutRes{
            Uri:    PLoginLogoutRes_uri,
            Rescode:RES_OK,
            Uid:    uid,
        })
        gSessionMgr.onSessionLogout(session)
        gSessionMgr.onSessionClose(session)
    }

    logInfo("Server: recv PLoginLogout message, leave")
}

func handleLoginPing(msg *loginproto.PLoginPing, session socket.Session)  {
    logInfo("Server: recv PLoginPing message, enter")

    session.SetParam(SESSION_LAST_PING_TIME, utils.CurrentTimeMillSecond())
    uid, _ := getSessionParamAsString(session, SESSION_USER_ID)
    session.Send(&loginproto.PLoginPingRes{
        Uri:    PLoginPingRes_uri,
        Uid:    uid,
        Clientts: msg.Clientts,
        Serverts: utils.CurrentTimeMillSecond(),
    })

    logInfo("Server: recv PLoginPing message, leave")
}

