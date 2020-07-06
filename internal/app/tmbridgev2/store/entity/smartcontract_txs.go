package entity

import (
	"math/big"
	"time"

	"github.com/globalsign/mgo/bson"
)

// SmartContractTransaction
type SmartContractTransaction struct {
	ID                bson.ObjectId `json:"id" bson:"_id"`
	SmartContract     bson.ObjectId `json:"smartcontract, omitempty" bson:",omitempty"`
	Status            bool          `json:"status, omitempty" bson:",omitempty"`
	IsPending         bool          `json:"is_pending, omitempty" bson:",omitempty"`
	TxIndex           uint          `json:"tx_index, omitempty" bson:",omitempty"`
	Hash              string        `json:"hash, omitempty" bson:",omitempty"`
	BlockHash         string        `json:"block_hash, omitempty" bson:",omitempty"`
	BlockNumber       uint64        `json:"block_number, omitempty" bson:",omitempty"`
	CumulativeGasUsed uint64        `json:"cumulative_gas, omitempty" bson:",omitempty"`
	From              string        `json:"from, omitempty" bson:",omitempty"`
	Gas               uint64        `json:"gas, omitempty" bson:",omitempty"`
	GasPrice          *big.Int      `json:"gas_price, omitempty" bson:",omitempty"`
	Input             []byte        `json:"input, omitempty" bson:",omitempty"`
	Nonce             uint64        `json:"nonce, omitempty" bson:",omitempty"`
	Timestamp         time.Time     `json:"timestamp" bson:"createdAt"`
	To                string        `json:"to, omitempty" bson:",omitempty"`
	Value             *big.Int      `json:"value, omitempty" bson:",omitempty"`
	CreatedAt         time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// SmartContractTransactionRecordUpdate
type SmartContractTransactionRecordUpdate struct {
	*SmartContractTransaction
}

// GetBSON insert record database
func (csi *SmartContractTransaction) GetBSON() (interface{}, error) {
	d := SmartContractTransaction{
		ID:                bson.NewObjectId(),
		SmartContract:     csi.SmartContract,
		Status:            csi.Status,
		BlockHash:         csi.BlockHash,
		BlockNumber:       csi.BlockNumber,
		IsPending:         csi.IsPending,
		Hash:              csi.Hash,
		CumulativeGasUsed: csi.CumulativeGasUsed,
		From:              csi.From,
		To:                csi.To,
		Gas:               csi.Gas,
		GasPrice:          csi.GasPrice,
		Nonce:             csi.Nonce,
		Timestamp:         csi.Timestamp,
		Value:             csi.Value,
		Input:             csi.Input,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	return d, nil
}

// GetBSON upsert
func (csi SmartContractTransactionRecordUpdate) GetBSON() (interface{}, error) {
	now := time.Now()
	set := bson.M{
		"smart_contract":      csi.SmartContract,
		"status":              csi.Status,
		"block_hash":          csi.BlockHash,
		"block_number":        csi.BlockNumber,
		"is_pending":          csi.IsPending,
		"hash":                csi.Hash,
		"cumulative_gas_used": csi.CumulativeGasUsed,
		"from":                csi.From,
		"to":                  csi.To,
		"gas":                 csi.Gas,
		"gas_price":           csi.GasPrice,
		"nonce":               csi.Nonce,
		"timestamp":           csi.Timestamp,
		"value":               csi.Value,
		"input":               csi.Input,
		"updatedAt":           now,
	}
	setOnInsert := bson.M{
		"_id":       bson.NewObjectId(),
		"createdAt": now,
	}
	update := bson.M{
		"$set":         set,
		"$setOnInsert": setOnInsert,
	}

	return update, nil
}
