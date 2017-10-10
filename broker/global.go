package main

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
