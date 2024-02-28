package trainers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TrainerTransfer struct {
	From  common.Address
	To    common.Address
	Token *big.Int
}

type TrainerMinting struct {
	Model uint16
	Token *big.Int
	To    common.Address
}

// Meter testnet

func IsTransfer(hash common.Hash) bool {
	return hash == crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
}

func IsNewMint(hash common.Hash) bool {
	return hash.String() == "0xfcbd9479bb682193b3eb060a58134eda82cf0a141904328119cfb9c15e6e171e"
}

func SubcribeToTransfers(client *ethclient.Client, address string) (chan types.Log, ethereum.Subscription, error) {
	contractAddress := common.HexToAddress(address)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)

	return logs, sub, err
}
