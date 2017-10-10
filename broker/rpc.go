package main

import (
    "errors"
)

func loginByPassport(passport, password string, params map[string]string) (uid, token string, err error) {
    if passport == "liujia" && password == "123456" {
        return "1", "hello_liujia", nil
    }
    return "","", errors.New("invalid password")
}

func loginByToken(uid, token string, params map[string]string) (err error) {
    if uid == "1" && token == "hello_liujia" {
        return nil
    }
    return errors.New("invalid token")
}

func logout(uid, token string) (err error) {
    if uid == "1" && token == "hello_liujia" {
        return nil
    }
    return errors.New("invalid token")
}
