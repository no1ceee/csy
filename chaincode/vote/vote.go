package main

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type VoteContract struct {
	contractapi.Contract
}



type Vote struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Voter     string `json:"voter"`
	Timestamp string `json:"timestamp"`
}

func (s *VoteContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}
// 投票
func (vc *VoteContract) Vote(ctx contractapi.TransactionContextInterface, id string, name string) error {
	// 获取调用方身份
	voterID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client identity: %v", err)
	}

	// 获取当前时间（UNIX时间戳）
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get transaction timestamp: %v", err)
	}

	// 封装投票信息
	
	vote := Vote{
		ID:        id,
		Name:      name,
		Voter:     voterID,
		Timestamp: timestamp.String(),
	}
	voteAsBytes, _ := json.Marshal(vote)
	// 将投票信息存储到Ledger中
	err = ctx.GetStub().PutState(id, voteAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put state: %v", err)
	}

	log.Printf("voter %s voted for %s at %v", voterID, name, timestamp)

	return nil
}

// 查询投票结果
func (vc *VoteContract) QueryResult(ctx contractapi.TransactionContextInterface, id string) (*Vote, error) {
	// 查询Ledger中指定ID的投票信息
	voteBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("from world state: %v", err)
	}
	if voteBytes == nil {
		return nil, fmt.Errorf("the vote %s does not exist", id)
	}

	// unmarshal，并返回
	vote := new(Vote)
	err = json.Unmarshal(voteBytes, vote)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return vote, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&VoteContract{})
	if err != nil {
		log.Panicf("Error creating vote chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting vote chaincode: %v", err)
	}
}
