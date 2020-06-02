package entity

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
)

type UserDao struct {
	collectionName string
	dbName         string
}

// NewUserDao returns a new instance of UserDao.
func NewUserDao(dbname string) *UserDao {
	dbName := dbname
	collection := "users"
	return &UserDao{collection, dbName}
}

func (dao *UserDao) Upsert(cs *User) error {
	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"address": cs.Address}, &UserRecordUpdate{cs})

	if err != nil {
		return err
	}
	return nil
}

func (dao *UserDao) IsAuthen(token string) (*User, error) {
	res := &User{}
	err := db.GetOne(dao.dbName, dao.collectionName, bson.M{"session.token": token}, &res)

	if err != nil {
		return nil, err
	}

	if res.Session.ExpiredAt.Unix() < time.Now().Unix() {
		return nil, errors.New("TOKEN_EXPIRED")
	}

	return res, nil
}

func (dao *UserDao) IsExist(address string) (*User, error) {
	res := &User{}
	err := db.GetOne(dao.dbName, dao.collectionName, bson.M{"address": address}, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
