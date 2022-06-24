package utils

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	chainModel "github.com/jason-bateman/go-erc-standard-contract/model"
)

func GetLatestBlockNumWithClient(client *ethclient.Client) (uint64, error) {
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

func GetLatestBlockNumWithRPC(rpc string) (uint64, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return 0, err
	}
	defer client.Close()
	return GetLatestBlockNumWithClient(client)
}

func GetAddressTxNonceWithClient(client *ethclient.Client, userAddress string) (*uint64, error) {
	txNonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(userAddress))
	if err != nil {
		return nil, err
	}

	return &txNonce, nil
}

func GetAddressTxNonceWithRPC(rpc string, userAddress string) (*uint64, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return GetAddressTxNonceWithClient(client, userAddress)
}

func GetChainIdWithClient(client *ethclient.Client) (uint64, error) {
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return 0, err
	}
	return chainId.Uint64(), nil
}

func GetChainIdWithRPC(rpc string) (uint64, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return 0, err
	}
	defer client.Close()
	return GetChainIdWithClient(client)
}

func GetNativeBalanceWithClient(client *ethclient.Client, address string) (float64, error) {

	operator := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), operator, nil)
	if err != nil {
		return 0, err
	}
	value, _ := FormatUnits2Float64(balance, 18)
	return value, nil
}

func GetNativeBalanceWithRPC(rpc string, address string) (float64, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return 0, err
	}
	defer client.Close()
	return GetNativeBalanceWithClient(client, address)
}

func MergeEventMessage(src, dest []*chainModel.EthereumEventMessage) []*chainModel.EthereumEventMessage {
	var lenSrc, lenDest int
	if src == nil {
		lenSrc = 0
	} else {
		lenSrc = len(src)
	}

	if dest == nil {
		lenDest = 0
	} else {
		lenDest = len(dest)
	}

	result := make([]*chainModel.EthereumEventMessage, lenSrc+lenDest)
	if lenSrc > 0 {
		copy(result, src)
	}
	if lenDest > 0 {
		copy(result[lenSrc:], dest)
	}

	return result
}
