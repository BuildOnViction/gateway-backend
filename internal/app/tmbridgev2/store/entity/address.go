package entity

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Address
type Address struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	ProjectID      bson.ObjectId `json:"project_id" bson:"project_id"`
	CoinType       string        `json:"coin_type" bson:"coin_type"`
	Address        string        `json:"address" bson:"address"`
	DepositAddress string        `json:"deposit_address" bson:"deposit_address"`
	AccountIndex   uint64        `json:"account_index" bson:"account_index"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// AddressRecordUpdate
type AddressRecordUpdate struct {
	*Address
}

// GetBSON insert record database
func (csi *Address) GetBSON() (interface{}, error) {
	d := Address{
		ProjectID:      csi.ProjectID,
		CoinType:       csi.CoinType,
		Address:        csi.Address,
		DepositAddress: csi.DepositAddress,
		AccountIndex:   csi.AccountIndex,
		CreatedAt:      csi.CreatedAt,
		UpdatedAt:      csi.UpdatedAt,
	}
	return d, nil
}

// GetBSON upsert
func (csip AddressRecordUpdate) GetBSON() (interface{}, error) {
	now := time.Now()
	set := bson.M{
		"project":         csip.ProjectID,
		"coin_type":       csip.CoinType,
		"address":         csip.Address,
		"deposit_address": csip.DepositAddress,
		"account_index":   csip.AccountIndex,
		"updatedAt":       now,
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
