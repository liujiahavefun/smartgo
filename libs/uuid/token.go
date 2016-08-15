package uuid

import (
	"fmt"
	"math/rand"
	"smartgo/libs/utils"
	"time"
)

var (
	token_key = "liujianiubi"
)

func init() {
	token_key = utils.GetMacAddress()
}

func CreateToken() string {
	now := time.Now().UnixNano()
	r := rand.New(rand.NewSource(now))
	uuid := NewV5(NamespaceOID, fmt.Sprintf("%v%v%v", token_key, now, r.Int63()))
	return uuid.StringNoDash()
}
