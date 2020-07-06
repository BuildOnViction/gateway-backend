package database

import (
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
	"github.com/anhntbk08/gateway/internal/common"
	"github.com/globalsign/mgo"
)

// Mongo database access object
type Mongo struct {
	UserDao            *entity.UserDao
	SessionDao         *entity.SessionDao
	ProjectDao         *entity.ProjectDao
	AddressDao         *entity.AddressDao
	SmartContractDao   *entity.SmartContractDao
	SmartContractTxDao *entity.SmartContractTransactionDao
}

// NewMongo new database instace
func NewMongo(url string, dbname string, logger common.Logger) (*Mongo, error) {
	mgo.SetLogger(logger)
	_, err := entity.InitSession(url, nil)
	if err != nil {
		return nil, err
	}

	return &Mongo{
		UserDao:            entity.NewUserDao(dbname),
		SessionDao:         entity.NewSessionDao(dbname),
		ProjectDao:         entity.NewProjectDao(dbname),
		AddressDao:         entity.NewAddressDao(dbname),
		SmartContractDao:   entity.NewSmartContractDao(dbname),
		SmartContractTxDao: entity.NewSmartContractTransactionDao(dbname),
	}, nil
}
