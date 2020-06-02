package database

import (
	"github.com/anhntbk08/gateway/internal/app/gateway/store/entity"
)

// Mongo database access object
type Mongo struct {
	UserDao    *entity.UserDao
	SessionDao *entity.SessionDao
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
	}, nil
}
