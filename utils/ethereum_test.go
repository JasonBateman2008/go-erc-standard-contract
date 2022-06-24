package utils

import (
	"math/big"
	"testing"
)

func TestUnits(t *testing.T) {
	var res string
	res, err := ParseUnits2String(big.NewInt(1), 8)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)

	res, err = ParseUnits2String("1", 8)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
	res, err = ParseUnits2String(float64(1222), 18)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)

	res, err = ParseUnits2String(0.00039334, 18)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)

	res, err = ParseUnits2String(0.2, 18)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestAddress(t *testing.T) {
	// error test
	errAddr := "0xxxxxx"
	status := IsHexAddressValid(errAddr)
	if status {
		t.Errorf("%s address should not a address", errAddr)
	}

	address := "0x67497baefe2bdf028bb7fed35c7f211ce10469f6"
	eIP55Address := "0x67497bAEfe2BDF028bb7FEd35C7F211cE10469F6"
	addr, err := AddressFormatByEIP55(address)
	if err != nil {
		t.Error(err)
	}
	if addr != eIP55Address {
		t.Errorf("%s fomat eip55 should be %s", address, eIP55Address)
	}

	status, err = AddressEqual(address, eIP55Address)
	if err != nil {
		t.Error(err)
	}
	if !status {
		t.Errorf("%s %s should be equal ", address, eIP55Address)
	}

	status, err = IsAddressFormatByEIP55(eIP55Address)
	if err != nil {
		t.Error(err)
	}
	if !status {
		t.Errorf("%s should be Eip55 ", eIP55Address)
	}

	status, err = IsAddressFormatByEIP55(address)
	if err != nil {
		t.Error(err)
	}
	if status {
		t.Errorf("%s should be Eip55 ", address)
	}

	t.Log("Test Address module is OK")
}
