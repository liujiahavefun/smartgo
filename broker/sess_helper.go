package main

import (
    "smartgo/libs/socket"
    "errors"
)

func getSessionParamAsString(session socket.Session, key string) (val string, err error)  {
    if val, ok := session.GetParam(key).(string); ok {
        return val, nil
    }else {
        return "", errors.New("failed to get key")
    }
}

func getSessionParamAsBool(session socket.Session, key string) (val bool, err error)  {
    if val, ok := session.GetParam(key).(bool); ok {
        return val, nil
    }else {
        return false, errors.New("failed to get key")
    }
}

func getSessionParamAsInt64(session socket.Session, key string) (val int64, err error)  {
    if val, ok := session.GetParam(key).(int64); ok {
        return val, nil
    }else {
        return 0, errors.New("failed to get key")
    }
}

func getSessionParamAsTimer(session socket.Session, key string) (val *socket.Timer, err error)  {
    if val, ok := session.GetParam(key).(*socket.Timer); ok {
        return val, nil
    }else {
        return nil, errors.New("failed to get key")
    }
}
