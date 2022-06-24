package utils

import (
	"math/big"
	"strings"
)

func String2BigInt(src string) *big.Int {
	value := new(big.Int)
	if strings.Contains(src, "0x") {
		value.SetString(strings.TrimLeft(src, "0x"), 16)
	} else {
		value.SetString(src, 0)
	}
	return value
}
