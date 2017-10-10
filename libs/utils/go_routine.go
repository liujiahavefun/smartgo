package utils

import (
    "runtime"
    "strings"
    "strconv"
    "fmt"
)

/*
 * 获取当前的goroutine的id，这个俺也是从网上抄的
 */
func GoID() int {
    var buf [64]byte
    n := runtime.Stack(buf[:], false)
    idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
    id, err := strconv.Atoi(idField)
    if err != nil {
        panic(fmt.Sprintf("cannot get goroutine id: %v", err))
    }
    return id
}