package main

import (
    loginproto "smartgo/proto/login_event"
    "smartgo/libs/socket"
    "smartgo/libs/utils"
)

func handleLoginByPassport(msg *loginproto.PLoginByPassport, session socket.Session)  {
    logInfo("Server: recv PLoginByPassport message")

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
}

func handleLoginByToken(msg *loginproto.PLoginByToken, session socket.Session)  {
    logInfo("Server: recv PLoginByToken message")

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
}

func handleLoginLogout(msg *loginproto.PLoginLogout, session socket.Session)  {
    logInfo("Server: recv PLoginLogout message")

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
}

func handleLoginPing(msg *loginproto.PLoginPing, session socket.Session)  {
    logInfo("Server: recv PLoginPing message")
    session.SetParam(SESSION_LAST_PING_TIME, utils.CurrentTimeMillSecond())
    uid, _ := getSessionParamAsString(session, SESSION_USER_ID)
    session.Send(&loginproto.PLoginPingRes{
        Uri:    PLoginPingRes_uri,
        Uid:    uid,
        Clientts: msg.Clientts,
        Serverts: utils.CurrentTimeMillSecond(),
    })
}

