package entity

import (
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
)

// SmartContract
type SmartContract struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Address      string        `json:"address" bson:"address"`
	Type         string        `json:"type" bson:"type"` // TRC20, TRC21, TRC721
	IsSyncing    bool          `json:"is_syncing" bson:",omitempty"`
	ScannedIndex int64         `json:"scanned_index" bson:"scanned_index"`
	Description  string        `json:"description, omitempty" bson:",omitempty"`
	Project      bson.ObjectId `json:"project" bson:"project"`
	CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// SmartContractRecordUpdate
type SmartContractRecordUpdate struct {
	*SmartContract
}

// GetBSON insert record database
func (csi *SmartContract) GetBSON() (interface{}, error) {
	d := SmartContract{
		ID:           bson.NewObjectId(),
		Project:      csi.Project,
		Type:         csi.Type,
		IsSyncing:    csi.IsSyncing,
		ScannedIndex: csi.ScannedIndex,
		Address:      strings.ToLower(csi.Address),
		Description:  csi.Description,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return d, nil
}

// GetBSON upsert
func (csip SmartContractRecordUpdate) GetBSON() (interface{}, error) {
	now := time.Now()
	set := bson.M{
		"address":       strings.ToLower(csip.Address),
		"type":          csip.Type,
		"is_syncing":    csip.IsSyncing,
		"scanned_index": csip.ScannedIndex,
		"project":       csip.Project,
		"description":   csip.Description,
		"updatedAt":     now,
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
