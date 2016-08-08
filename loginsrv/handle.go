package main

import (
	"fmt"
	"smartgo/loginsrv/proto/thrift/gen-go/login/rpc"
)

type LoginServiceHandler struct {
}

func NewLoginSrvHandler() *LoginServiceHandler {
	return &LoginServiceHandler{}
}

func (p *LoginServiceHandler) LoginByPasswd(username string, passwd string, fromtype string, params map[string]string) (r string, err error) {
	userObj, err := GetPersonByName(username)
	if err != nil {
		return "", ErrorUserNotFound
	}

	if passwd == userObj.Passwd {
		return "token", nil
	}
	return "", nil
}

func (p *LoginServiceHandler) LoginByToken(uid string, token string, fromtype string, params map[string]string) (r string, err error) {
	fmt.Print("LoginByToken()\n")

	excep := rpc.NewLoginException()
	excep.Errno = 1
	excep.Errmsg = "interface not implemented"
	err = excep

	return "", err
}

func (p *LoginServiceHandler) Logout(uid string, token string, params map[string]string) (r bool, err error) {
	fmt.Print("Logout()\n")
	return false, nil
}

func (p *LoginServiceHandler) IsOnline(uid string, params map[string]string) (r bool, err error) {
	fmt.Print("IsOnline()\n")
	return false, nil
}

func (p *LoginServiceHandler) KickOff(uid string) (r bool, err error) {
	fmt.Print("KickOff()\n")
	return false, nil
}

func (p *LoginServiceHandler) WhichNode(uid string) (r string, err error) {
	fmt.Print("WhichNode()\n")
	return "", nil
}
