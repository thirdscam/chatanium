package Util

import (
	"fmt"
	"math/big"

	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/steebchen/prisma-client-go/runtime/types"
)

func StringToBigint(id string) types.BigInt {
	n := new(big.Int)
	n, ok := n.SetString(id, 10)
	if !ok {
		fmt.Println("SetString: error")
		Log.Error.Fatalf("Failed to convert ID: %v", id)
	}

	return types.BigInt(n.Int64())
}
