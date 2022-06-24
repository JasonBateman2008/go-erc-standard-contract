package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math"
	"math/big"
	"reflect"
	"strings"
)

func IsHexAddressValid(address string) bool {
	return common.IsHexAddress(address)
}

func AddressFormatByEIP55(address string) (string, error) {
	if !IsHexAddressValid(address) {
		return "", fmt.Errorf("address %s is not a hex encode address", address)
	}
	return common.HexToAddress(address).Hex(), nil
}

func IsAddressFormatByEIP55(address string) (bool, error) {
	eip55Addr, err := AddressFormatByEIP55(address)
	if err != nil {
		return false, err
	}
	return eip55Addr == address, nil
}

func AddressEqual(address1, address2 string) (bool, error) {
	if !IsHexAddressValid(address1) {
		return false, fmt.Errorf("address %s is not a hex encode address", address1)
	}
	if !IsHexAddressValid(address1) {
		return false, fmt.Errorf("address %s is not a hex encode address", address2)
	}
	return strings.EqualFold(address1, address2), nil
}

func ParseUnits2String(iValue interface{}, decimals uint8) (string, error) {
	var (
		succeed bool
		value   *big.Float
	)
	value = new(big.Float)
	decimal := big.NewFloat(math.Pow(10, float64(decimals)))
	switch v := iValue.(type) {
	case string:
		value, succeed = value.SetString(v)
		if !succeed {
			return "", fmt.Errorf("wrong string content:%s", v)
		}
	case float64:
		value = value.SetFloat64(v)
	case *float64:
		value = value.SetFloat64(*v)
	case int64:
		value = value.SetInt64(v)
	case *int64:
		value = value.SetInt64(*v)
	case int:
		value = value.SetInt64(int64(v))
	case *int:
		value = value.SetInt64(int64(*v))
	case big.Int:
		value = value.SetInt(&v)
	case *big.Int:
		value = value.SetInt(v)
	default:
		return "", fmt.Errorf("unsupported type:%+v", reflect.TypeOf(iValue))
	}
	res := new(big.Int)
	new(big.Float).Mul(value, decimal).Int(res)
	return res.String(), nil

}

func FormatUnits2Float64(iWei interface{}, decimals uint8) (float64, error) {
	var wei *big.Int

	switch v := iWei.(type) {
	case string:
		wei = String2BigInt(v)
	case int64:
		wei = big.NewInt(v)
	case big.Int:
		wei = &v
	case *big.Int:
		wei = v
	default:
		return 0.0, fmt.Errorf("unsupported type:%+v", reflect.TypeOf(iWei))
	}

	decimal := big.NewFloat(math.Pow(10, float64(decimals)))
	weiFloat := new(big.Float).SetInt(wei)
	value, _ := new(big.Float).Quo(weiFloat, decimal).Float64()

	return value, nil
}
