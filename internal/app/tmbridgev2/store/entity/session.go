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
