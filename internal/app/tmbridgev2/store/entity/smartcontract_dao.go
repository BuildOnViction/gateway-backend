package entity

import (
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type SmartContractDao struct {
	*DAO
}

// NewSmartContractDao returns a new instance of SmartContractDao.
func NewSmartContractDao(dbname string) *SmartContractDao {
	dbName := dbname
	collectionName := "smartcontracts"
	i1 := mgo.Index{
		Key:    []string{"address"},
		Unique: true,
	}

	err := db.Session.DB(dbName).C(collectionName).EnsureIndex(i1)
	if err != nil {
		panic(err)
	}

	return &SmartContractDao{
		DAO: &DAO{
			collectionName: collectionName,
			dbName:         dbName,
		},
	}
}

func (dao *SmartContractDao) BulkRemove(addresses []string) error {
	lAddresses := []string{}
	for i := 0; i < len(addresses); i++ {
		lAddresses = append(lAddresses, strings.ToLower(addresses[i]))
	}
	return dao.RemoveAll(bson.M{
		"address": bson.M{
			"$in": lAddresses,
		},
	})
}

func (dao *SmartContractDao) IsSyncing(address string) bool {
	res := SmartContract{}

	err := dao.GetOne(bson.M{
		"address":    strings.ToLower(address),
		"is_syncing": true,
	}, &res)

	if err != nil {
		return false
	}

	return true
}

func (dao *SmartContractDao) GetByAddress(address string) *SmartContract {
	res := SmartContract{}

	err := dao.GetOne(bson.M{
		"address": strings.ToLower(address),
	}, &res)

	if err != nil {
		return nil
	}

	return &res
}

func (dao *SmartContractDao) StartSync(address string) error {
	err := dao.Update(bson.M{
		"address": strings.ToLower(address),
	}, bson.M{
		"$set": bson.M{
			"is_syncing": true,
		},
	})

	return err
}

func (dao *SmartContractDao) StopSync(addressID bson.ObjectId, scannedTo int64) error {
	_, err := dao.Upsert(bson.M{
		"_id": addressID,
	}, bson.M{
		"$set": bson.M{
			"is_syncing":    false,
			"scanned_index": scannedTo,
		},
	})

	return err
}
