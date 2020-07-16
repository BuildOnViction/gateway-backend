package entity

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// ScannedIndex
type ScannedIndex struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Type         string        `json:"type" bson:"type"`
	ScannedIndex uint64        `json:"scanned_index" bson:"scanned_index"`
	Description  string        `json:"description, omitempty" bson:",omitempty"`
	CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// ScannedIndexRecordUpdate
type ScannedIndexRecordUpdate struct {
	*ScannedIndex
}

// GetBSON insert record database
func (csi *ScannedIndex) GetBSON() (interface{}, error) {
	d := ScannedIndex{
		ID:           bson.NewObjectId(),
		Type:         csi.Type,
		ScannedIndex: csi.ScannedIndex,
		Description:  csi.Description,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return d, nil
}

// GetBSON upsert
func (csip ScannedIndexRecordUpdate) GetBSON() (interface{}, error) {
	now := time.Now()
	set := bson.M{
		"type":          csip.Type,
		"scanned_index": csip.ScannedIndex.ScannedIndex,
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
