namespace go login.rpc

/*
 * 登录服务异常、错误
 */
exception LoginException {
  1: i32 errno,
  2: string errmsg
}

/*
 * 登录服务RPC
 */
service LoginService {
    //用户名密码登录,成功返回token,否则异常
    string loginByPasswd(1:string username, 2:string passwd, 3:string type, 4:map<string, string> params) throws (1:LoginException excep),

    //用户id和token登录,成功返回token,否则异常
    string loginByToken(1:string uid, 2:string token, 3:string type, 4:map<string, string> params) throws (1:LoginException excep),

    //登出
    bool logout(1:string uid, 2:string token, 3:map<string, string> params),

    //用户是否在线
    bool isOnline(1:string uid, 2:map<string, string> params),

    //踢人
    bool kickOff(1:string uid),

    //用户在哪个nodeSrv
    string whichNode(1:string uid),
}