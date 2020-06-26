package database

import (
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
)

// Mongo database access object
type Mongo struct {
	UserDao    *entity.UserDao
	SessionDao *entity.SessionDao
	ProjectDao *entity.ProjectDao
	AddressDao *entity.AddressDao
}

// NewMongo new database instace
func NewMongo(url string, dbname string) (*Mongo, error) {
	_, err := entity.InitSession(url, nil)
	if err != nil {
		return nil, err
	}

	return &Mongo{
		UserDao:    entity.NewUserDao(dbname),
		SessionDao: entity.NewSessionDao(dbname),
		ProjectDao: entity.NewProjectDao(dbname),
		AddressDao: entity.NewAddressDao(dbname),
	}, nil
}
