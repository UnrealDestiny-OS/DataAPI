package idle

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BalanceInjection struct {
	Wallet common.Address
	Amount *big.Int
	Date   *big.Int
}

func IsBalanceInjection(hash common.Hash) bool {
	return hash.String() == "0x9e9b7e29099f2d6721d8268a517ee7ada50ae5d99eab1c668a3ba067877023fb"
}

func SubscribeToContract(client *ethclient.Client, address string) (chan types.Log, ethereum.Subscription, error) {
	contractAddress := common.HexToAddress(address)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)

	return logs, sub, err
}
