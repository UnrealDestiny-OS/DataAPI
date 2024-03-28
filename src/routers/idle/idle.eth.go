package idle

import (
	"context"
	"errors"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

type BalanceInjection struct {
	Wallet common.Address
	Amount *big.Int
	Date   *big.Int
}

type TakeFees struct {
	Wallet        common.Address
	WalletBalance *big.Int
	TakedFees     *big.Int
	GameRewards   *big.Int
}

type BalanceWithdraw struct {
	Wallet common.Address
	Amount *big.Int
	Date   *big.Int
}

func IsTakeFees(hash common.Hash) bool {
	return hash.String() == "0xd8f5815a7fb42d63779ba883dd563f4d249624467c7abbdd5b8976cb80cd3f8f"
}

func IsBalanceInjection(hash common.Hash) bool {
	return hash.String() == "0x9e9b7e29099f2d6721d8268a517ee7ada50ae5d99eab1c668a3ba067877023fb"
}

func IsBalanceWithdraw(hash common.Hash) bool {
	return hash.String() == "0x08b48106681a11cad8a74d362ab593c5e8673b9d713b13e2469b8ec8dfa887de"
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

func ProcessInjectionEvent(eventType string, abi abi.ABI, data []byte, event *BalanceInjection) error {
	var ok bool

	injectBalanceInterface, err := abi.Unpack(eventType, data)

	if err != nil {
		return errors.New("error unpacking the injection event data")
	}

	event.Wallet, ok = injectBalanceInterface[0].(common.Address)

	if !ok {
		return errors.New("error parsing the injection event data (Wallet)")
	}

	event.Amount, ok = injectBalanceInterface[1].(*big.Int)

	if !ok {
		return errors.New("error parsing the injection event data (Amount)")
	}

	event.Date, ok = injectBalanceInterface[2].(*big.Int)

	if !ok {
		return errors.New("error parsing the injection event data (Date)")
	}

	return nil
}

func ProcessTakeFeesEvent(abi abi.ABI, data []byte, event *TakeFees) error {
	var ok bool

	injectBalanceInterface, err := abi.Unpack("TakeFees", data)

	if err != nil {
		return errors.New("error unpacking the injection event data")
	}

	event.Wallet, ok = injectBalanceInterface[0].(common.Address)

	if !ok {
		return errors.New("error parsing the injection event data (Wallet)")
	}

	event.WalletBalance, ok = injectBalanceInterface[1].(*big.Int)

	if !ok {
		return errors.New("error parsing the injection event data (WalletBalance)")
	}

	event.TakedFees, ok = injectBalanceInterface[2].(*big.Int)

	if !ok {
		return errors.New("error parsing the injection event data (TakedFees)")
	}

	event.GameRewards, ok = injectBalanceInterface[3].(*big.Int)

	if !ok {
		return errors.New("error parsing the injection event data (GameRewards)")
	}

	return nil
}

func GetSignatureByMethod(method string) []byte {
	sign := []byte(method)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(sign)
	methodID := hash.Sum(nil)[:4]
	return methodID
}

func GetTrainerJoinMethodSignature() []byte {
	return GetSignatureByMethod("joinWithTrainer(address,uint256)")
}

func GetCollectTransactionPointsSignature() []byte {
	return GetSignatureByMethod("collectTransPoints(address,uint256)")
}

func GetCollectIdlePointsSignature() []byte {
	return GetSignatureByMethod("collectIDLEPoints(address,uint256)")
}

func GetBuyImprovementSignature() []byte {
	return GetSignatureByMethod("buyTrainerImprovement(address,uint256,uint256)")
}

func GetTrainerJoinTxData(userAddress string, trainer int) []byte {
	var data []byte

	methodID := GetTrainerJoinMethodSignature()
	paddedAddress := ParseAddressForTxData(userAddress)
	paddedAmount := ParseIntForTxData(trainer)

	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	return data
}

func GetCollectTransactionPointsData(userAddress string, trainer int) []byte {
	var data []byte

	methodID := GetCollectTransactionPointsSignature()
	paddedAddress := ParseAddressForTxData(userAddress)
	paddedAmount := ParseIntForTxData(trainer)

	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	return data
}

func GetCollectionIdlePointsData(userAddress string, trainer int) []byte {
	var data []byte

	methodID := GetCollectIdlePointsSignature()
	paddedAddress := ParseAddressForTxData(userAddress)
	paddedAmount := ParseIntForTxData(trainer)

	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	return data
}

func GetBuyImprovementData(userAddress string, trainer int, improvement int) []byte {
	var data []byte

	methodID := GetBuyImprovementSignature()
	paddedAddress := ParseAddressForTxData(userAddress)
	paddedAmount := ParseIntForTxData(trainer)
	paddedImprovement := ParseIntForTxData(improvement)

	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	data = append(data, paddedImprovement...)

	return data
}

func ParseAddressForTxData(address string) []byte {
	return common.LeftPadBytes(common.HexToAddress(address).Bytes(), 32)
}

func ParseIntForTxData(number int) []byte {
	amount := new(big.Int)
	amount.SetInt64(int64(number))
	return common.LeftPadBytes(amount.Bytes(), 32)
}

func GenerateTransactionExecutionSign(wallet string, chain int) string {
	return "Create transaction executor code with " + wallet + " on unrealdestiny.com using the chain " + strconv.Itoa(chain)
}
