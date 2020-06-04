package entity

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Keys struct {
	ID        string    `json:"id" bson:"id"`
	Secret    string    `json:"secret" bson:"secret"`
	ExpiredAt time.Time `json:"expiredAt" bson:"expiredAt"`
}

type Security struct {
	WhileListNamees  []string `json:"while_list_addresses" bson:"while_list_addresses"`
	WhileListOrigins []string `json:"while_list_origins" bson:"while_list_origins"`
}

// Project
type Project struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Keys      Keys          `json:"keys" bson:"keys"`
	Security  Security      `json:"security" bson:"security"`
	User      bson.ObjectId `json:"user_id" bson:"user_id"`
	Status    bool          `json:"status" bson:"status"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// ProjectRecordUpdate
type ProjectRecordUpdate struct {
	*Project
}

// GetBSON insert record database
func (csi *Project) GetBSON() (interface{}, error) {
	d := Project{
		User:      csi.User,
		Name:      csi.Name,
		Keys:      csi.Keys,
		Security:  csi.Security,
		Status:    csi.Status,
		CreatedAt: csi.CreatedAt,
		UpdatedAt: csi.UpdatedAt,
	}
	return d, nil
}

// GetBSON upsert
func (csip ProjectRecordUpdate) GetBSON() (interface{}, error) {
	now := time.Now()
	set := bson.M{
		"name":      csip.Name,
		"user":      csip.User,
		"keys":      csip.Keys,
		"security":  csip.Security,
		"status":    csip.Status,
		"updatedAt": now,
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
