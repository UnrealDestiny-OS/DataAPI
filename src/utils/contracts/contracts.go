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

type ProductionDeploymentAddresses struct {
	UDToken          string
	TrainersERC721   string
	TrainersDeployer string
	TrainersIDLE     string
}

var ProductionAddresses = ProductionDeploymentAddresses{UDToken: "0x70f87602111d878DC63a67E8750753d7a889Dfbd", TrainersERC721: "0x46E60d93bf8dfE6e2bc7f4aC789455A21352D44F", TrainersDeployer: "0xCa310C0FcfEb8166e3A39aee298d76bD8f9A8bF8", TrainersIDLE: "0x6b0eA35EAb6fca810Bba0B3bcda6BA81833FdD78"}

func (deployment *Deployment) JsonAbi() *string {
	b, err := json.Marshal(deployment.Abi)

	if err != nil {
		return nil
	}

	var parsed string = string(b)

	return &parsed
}

func LoadDeploymentsData(useProductionAddresses bool) *Deployments {
	jsonData, err := os.Open("./src/data/deployment.json")

	if err != nil {
		return nil
	}

	defer jsonData.Close()

	byteValue, _ := io.ReadAll(jsonData)

	var deployments Deployments

	json.Unmarshal(byteValue, &deployments)

	if useProductionAddresses {
		deployments.UDToken.Address = ProductionAddresses.UDToken
		deployments.TrainersERC721.Address = ProductionAddresses.TrainersERC721
		deployments.TrainersDeployer.Address = ProductionAddresses.TrainersDeployer
		deployments.TrainersIDLE.Address = ProductionAddresses.TrainersIDLE
	}

	return &deployments
}
