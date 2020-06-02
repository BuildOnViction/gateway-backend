package entity

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Session
type Session struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Address   string        `json:"address" bson:"address"`
	Token     string        `json:"token" bson:"token"`
	Used      bool          `json:"used" bson:"used"`
	ExpiredAt time.Time     `json:"expiredAt" bson:"expiredAt"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// SessionRecordUpdate
type SessionRecordUpdate struct {
	*Session
}

// func (csi *Session) GetBSON() (interface{}, error) {
// 	d := Session{
// 		Address:   csi.Address,
// 		CreatedAt: csi.CreatedAt,
// 		UpdatedAt: csi.UpdatedAt,
// 	}
// 	return d, nil
// }

// GetBSON upsert
// func (csip SessionRecordUpdate) GetBSON() (interface{}, error) {
// 	now := time.Now()
// 	set := bson.M{
// 		"address":   csip.Address,
// 		"updatedAt": now,
// 	}
// 	setOnInsert := bson.M{
// 		"_id":       bson.NewObjectId(),
// 		"createdAt": now,
// 	}
// 	update := bson.M{
// 		"$set":         set,
// 		"$setOnInsert": setOnInsert,
// 	}

// 	return update, nil
// }
