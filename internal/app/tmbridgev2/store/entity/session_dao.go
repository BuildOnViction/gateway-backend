package entity

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
)

type SessionDao struct {
	collectionName string
	dbName         string
}

// NewSessionDao returns a new instance of SessionDao.
func NewSessionDao(dbname string) *SessionDao {
	dbName := dbname
	collection := "sessions"
	return &SessionDao{collection, dbName}
}

func (dao *SessionDao) Create(cs *Session) error {
	cs.ID = bson.NewObjectId()
	err := db.Create(dao.dbName, dao.collectionName, cs)
	if err != nil {
		return err
	}
	return nil
}

func (dao *SessionDao) IsValidToken(address, token string) (*Session, error) {
	res := &Session{}
	err := db.GetOne(dao.dbName, dao.collectionName, bson.M{"address": address, "token": token, "used": false}, &res)

	if err != nil {
		return nil, err
	}

	if res.ExpiredAt.Unix() < time.Now().Unix() {
		return nil, errors.New("TOKEN_EXPIRED")
	}

	return res, nil
}

func (dao *SessionDao) Used(token string) (*Session, error) {
	err := db.Update(dao.dbName, dao.collectionName, bson.M{"token": token}, bson.M{
		"$set": bson.M{
			"used": true,
		},
	})
	return nil, err
}
