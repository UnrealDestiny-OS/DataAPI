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

// Meter testnet
var CONTRACT_TRAINERS_ERC721 = "0xeD3683F77b0685E109C085d8F380252B9bACa623"

func IsTransfer(hash common.Hash) bool {
	return hash == crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
}

func SubcribeToTransfers(client *ethclient.Client) (chan types.Log, ethereum.Subscription, error) {
	contractAddress := common.HexToAddress(CONTRACT_TRAINERS_ERC721)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)

	return logs, sub, err

}
