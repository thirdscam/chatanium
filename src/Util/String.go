package Util

import (
	"fmt"
	"math/big"

	"github.com/thirdscam/chatanium/src/Util/Log"
)

func Str2Int64(id string) int64 {
	n := new(big.Int)
	n, ok := n.SetString(id, 10)
	if !ok {
		fmt.Println("SetString: error")
		Log.Error.Fatalf("Failed to convert ID: %v", id)
	}

	return n.Int64()
}
