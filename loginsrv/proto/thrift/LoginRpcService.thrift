namespace go login.rpc

//用户登录服务
service RpcService {
    //用户名密码登录
    string loginByPasswd(1:string username, 2:string passwd, 3:map<string, string> params),
    //用户id和token登录
    bool loginByToen(1:string uid, 2:string token, 3:map<string, string> params),
    //登出
    bool logout(1:string uid, 2:string token, 3:map<string, string> params),
    //用户是否在线
    bool isOnline(1:string uid, 3:map<string, string> params),
    //
}