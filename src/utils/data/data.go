package data

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func ArrayContainsString(array []string, e string) bool {
	for _, a := range array {
		if a == e {
			return true
		}
	}

	return false
}

func SearchHashOnArray(array []common.Hash, e string) bool {
	for _, a := range array {
		if a.String() == e {
			return true
		}
	}

	return false
}

func DecodeTxParams(abi abi.ABI, v map[string]interface{}, data []byte) (map[string]interface{}, error) {
	m, err := abi.MethodById(data[:4])
	if err != nil {
		return map[string]interface{}{}, err
	}
	if err := m.Inputs.UnpackIntoMap(v, data[4:]); err != nil {
		return map[string]interface{}{}, err
	}
	return v, nil
}
