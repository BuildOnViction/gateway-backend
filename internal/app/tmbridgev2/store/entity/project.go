package entity

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Security struct {
	WhileListDomains []string `json:"while_list_domains" bson:"while_list_domains"`
	WhileListIps     []string `json:"while_list_ips" bson:"while_list_ips"`
}

type ProjectAddresses struct {
	MintingAddress      string   `json:"minting_address" bson:"minting_address"`
	WatchSmartContracts []string `json:"watch_smart_contracts" bson:"watch_smart_contracts"`
}

type Notification struct {
	WebHook string   `json:"web_hook" bson:"web_hook"`
	Emails  []string `json:"emails" bson:"emails"`
}

// Project
type Project struct {
	ID           bson.ObjectId    `json:"id" bson:"_id"`
	Name         string           `json:"name" bson:"name"`
	Secret       string           `json:"secret" bson:"secret"`
	Addresses    ProjectAddresses `json:"addresses" bson:"addresses"`
	Security     Security         `json:"security" bson:"security"`
	User         bson.ObjectId    `json:"user_id" bson:"user_id"`
	Notification Notification     `json:"notification" bson:"notification"`
	Status       bool             `json:"status" bson:"status"`
	CreatedAt    time.Time        `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time        `json:"updatedAt" bson:"updatedAt"`
}

// ProjectRecordUpdate
type ProjectRecordUpdate struct {
	*Project
}

// GetBSON insert record database
func (csi *Project) GetBSON() (interface{}, error) {
	d := Project{
		ID:           bson.NewObjectId(),
		User:         csi.User,
		Name:         csi.Name,
		Secret:       csi.Secret,
		Security:     csi.Security,
		Addresses:    csi.Addresses,
		Notification: csi.Notification,
		Status:       csi.Status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return d, nil
}

// GetBSON upsert
func (csip ProjectRecordUpdate) GetBSON() (interface{}, error) {
	now := time.Now()
	set := bson.M{
		"name":         csip.Name,
		"user":         csip.User,
		"secret":       csip.Secret,
		"security":     csip.Security,
		"status":       csip.Status,
		"addresses":    csip.Addresses,
		"notification": csip.Notification,
		"updatedAt":    now,
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
