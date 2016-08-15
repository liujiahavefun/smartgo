package main

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	"os"
	"time"
)

const (
	TIMEOUT = 5
)

var (
	rc *redis.Client
)

func getLoginKey(uid string) string {
	return uid + "_token"
}

func errorHandler(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func initRedis() error {
	var err error
	rc, err = redis.DialTimeout("tcp", ":6379", time.Duration(TIMEOUT)*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func echo() {
	r := rc.Cmd("echo", "Hello world!")
	errorHandler(r.Err)
}

func setLoginToken(uid string, token string) error {
	if len(uid) != 32 {
		return ErrorInvalidToken
	}

	r := rc.Cmd("set", getLoginKey(uid), token)
	return r.Err
}

func getLoginToken(uid string) (string, error) {
	if len(uid) != 32 {
		return "", ErrorInvalidToken
	}

	token, err := rc.Cmd("get", getLoginKey(uid)).Str()
	return token, err
}

func delLoginToken(uid string) error {
	if len(uid) != 32 {
		return ErrorInvalidToken
	}

	r := rc.Cmd("del", getLoginKey(uid))
	return r.Err
}
