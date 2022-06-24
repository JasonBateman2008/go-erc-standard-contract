package utils

import "testing"

func TestGetNativeBalanceWithRPC(t *testing.T) {
	balance, err := GetNativeBalanceWithRPC("https://data-seed-prebsc-2-s1.binance.org:8545", "0x35552c16704d214347f29Fa77f77DA6d75d7C752")
	if err != nil {
		t.Errorf("GetNativeBalanceWithRPC err:%+v\n", err)
		return
	}
	t.Logf("balance is:%f,\n", balance)
}
