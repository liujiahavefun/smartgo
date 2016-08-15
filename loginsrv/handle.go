package main

import (
	"fmt"
	"smartgo/libs/uuid"
	"smartgo/loginsrv/proto/thrift/gen-go/login/rpc"
)

type LoginServiceHandler struct {
}

func NewLoginSrvHandler() *LoginServiceHandler {
	return &LoginServiceHandler{}
}

func (p *LoginServiceHandler) LoginByPasswd(username string, passwd string, fromtype string, params map[string]string) (r string, err error) {
	fmt.Print("LoginByPasswd()\n")

	userObj, err := getPersonByName(username)
	if err != nil {
		return "", ErrorUserNotFound
	}

	if passwd == userObj.Passwd {
		token := uuid.CreateToken()
		err = setLoginToken(userObj.Uid, token)
		if err != nil {
			fmt.Printf("set user token to redis failed, uid[%v] token[%v]", userObj.Uid, token)
		}
		return token, nil
	}

	return "", ErrorInvalidPassword
}

func (p *LoginServiceHandler) LoginByToken(uid string, token string, fromtype string, params map[string]string) (r string, err error) {
	fmt.Print("LoginByToken()\n")

	t, err := getLoginToken(uid)
	if err != nil {
		return "", err
	}

	if t != token {
		return "", ErrorInvalidToken
	}

	return token, nil
}

func (p *LoginServiceHandler) Logout(uid string, token string, params map[string]string) (r bool, err error) {
	fmt.Print("Logout()\n")

	err = delLoginToken(uid)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *LoginServiceHandler) IsOnline(uid string, params map[string]string) (r bool, err error) {
	fmt.Print("IsOnline()\n")

	excep := rpc.NewLoginException()
	excep.Errno = 1
	excep.Errmsg = "interface not implemented"
	err = excep

	return false, err
}

func (p *LoginServiceHandler) KickOff(uid string) (r bool, err error) {
	fmt.Print("KickOff()\n")

	excep := rpc.NewLoginException()
	excep.Errno = 1
	excep.Errmsg = "interface not implemented"
	err = excep

	return false, err
}

func (p *LoginServiceHandler) WhichNode(uid string) (r string, err error) {
	fmt.Print("WhichNode()\n")

	excep := rpc.NewLoginException()
	excep.Errno = 1
	excep.Errmsg = "interface not implemented"
	err = excep

	return "", err
}
