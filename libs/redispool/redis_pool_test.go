package redispool_test

import (
	"fmt"
	"testing"
	"time"

	"smartgo/libs/redispool"
)

func Test_Base(t *testing.T) {
	var redisNetwork = "tcp"
	var redisAddress = "127.0.0.1:6379"
	fmt.Println("to connect redis on: ", redisAddress)
	t.Log("to connect redis on: ", redisAddress)
	rp, err := redispool.NewRedisPool(redisNetwork, redisAddress, 10, 300*time.Second)
	if err != nil {
		t.Errorf("connect redis failed: %v", err)
	} else {
		fmt.Println("connected redis on: ", redisAddress)
		t.Log("connected redis on: ", redisAddress)
	}

	rp.Cmd("ping")
}