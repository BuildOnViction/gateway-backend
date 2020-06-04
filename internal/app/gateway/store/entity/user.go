package entity

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type AuthenSession struct {
	Signature string    `json:"signature" bson:"signature"`
	Token     string    `json:"token" bson:"token"`
	ExpiredAt time.Time `json:"expiredAt" bson:"expiredAt"`
}

// User
type User struct {
	ID              bson.ObjectId `json:"id" bson:"_id"`
	Address         string        `json:"address" bson:"address"`
	Session         AuthenSession `json:"session" bson:"session"`
	MaximumProjects uint64        `json:"max_projects" bson:"max_projects"`
	PaymentPlan     uint8         `json:"payment_plan" bson:"payment_plan"`
	CreatedAt       time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// UserRecordUpdate
type UserRecordUpdate struct {
	*User
}

// GetBSON insert record database
func (csi *User) GetBSON() (interface{}, error) {
	d := User{
		Address:   csi.Address,
		Session:   csi.Session,
		CreatedAt: csi.CreatedAt,
		UpdatedAt: csi.UpdatedAt,
	}
	return d, nil
}

// GetBSON upsert
func (csip UserRecordUpdate) GetBSON() (interface{}, error) {
	now := time.Now()
	set := bson.M{
		"address":   csip.Address,
		"updatedAt": now,
		"session":   csip.Session,
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
