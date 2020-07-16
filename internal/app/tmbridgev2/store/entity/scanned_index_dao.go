package entity

import (
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type ScannedIndexDao struct {
	*DAO
}

// NewScannedIndexDao returns a new instance of ScannedIndexDao.
func NewScannedIndexDao(dbname string) *ScannedIndexDao {
	dbName := dbname
	collectionName := "scanned_index"
	i1 := mgo.Index{
		Key:    []string{"type"},
		Unique: true,
	}

	err := db.Session.DB(dbName).C(collectionName).EnsureIndex(i1)
	if err != nil {
		panic(err)
	}

	return &ScannedIndexDao{
		DAO: &DAO{
			collectionName: collectionName,
			dbName:         dbName,
		},
	}
}

func (dao *ScannedIndexDao) BulkRemove(addresses []string) error {
	lAddresses := []string{}
	for i := 0; i < len(addresses); i++ {
		lAddresses = append(lAddresses, strings.ToUpper(addresses[i]))
	}
	return dao.RemoveAll(bson.M{
		"address": bson.M{
			"$in": lAddresses,
		},
	})
}

func (dao *ScannedIndexDao) IsSyncing(address string) bool {
	res := ScannedIndex{}

	err := dao.GetOne(bson.M{
		"address":    strings.ToUpper(address),
		"is_syncing": true,
	}, &res)

	if err != nil {
		return false
	}

	return true
}

func (dao *ScannedIndexDao) GetCurrentBlock(cointype string) uint64 {
	res := ScannedIndex{}
	err := dao.GetOne(bson.M{
		"type": strings.ToUpper(cointype),
	}, &res)

	if err != nil {
		return 0
	}

	return res.ScannedIndex
}

func (dao *ScannedIndexDao) SetCurrentBlock(cointype string, block uint64) error {
	index := &ScannedIndex{
		ScannedIndex: block,
		Type:         strings.ToUpper(cointype),
	}
	_, err := dao.Upsert(bson.M{
		"type": strings.ToUpper(cointype),
	}, &ScannedIndexRecordUpdate{index})

	return err
}
