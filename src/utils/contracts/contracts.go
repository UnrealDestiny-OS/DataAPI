package contracts

import (
	"encoding/json"
	"io"
	"os"
)

// Read, load and parse deployed contracts data

type IAbiInput struct {
	Indexed      bool   `json:"indexed"`
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

type IAbi struct {
	Inputs          []IAbiInput `json:"inputs"`
	Outputs         []IAbiInput `json:"outputs"`
	Anonymous       bool        `json:"anonymous"`
	Name            string      `json:"name"`
	Type            string      `json:"type"`
	StateMutability string      `json:"stateMutability"`
}

type Deployment struct {
	Address string `json:"address"`
	Abi     []IAbi `json:"abi"`
}

type Deployments struct {
	UDToken          Deployment `json:"UDToken"`
	TrainersERC721   Deployment `json:"TrainersERC721"`
	TrainersDeployer Deployment `json:"TrainersDeployer"`
	TrainersIDLE     Deployment `json:"TrainersIDLE"`
}

func (deployment *Deployment) JsonAbi() *string {
	b, err := json.Marshal(deployment.Abi)

	if err != nil {
		return nil
	}

	var parsed string = string(b)

	return &parsed
}

func LoadDeploymentsData() *Deployments {
	jsonData, err := os.Open("./src/data/deployment.json")

	if err != nil {
		return nil
	}

	defer jsonData.Close()

	byteValue, _ := io.ReadAll(jsonData)

	var deployments Deployments

	json.Unmarshal(byteValue, &deployments)

	return &deployments
}
