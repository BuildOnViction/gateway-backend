package entity

import (
	"github.com/globalsign/mgo"
)

type SmartContractTransactionDao struct {
	*DAO
}

// NewSmartContractTransactionDao returns a new instance of SmartContractTransactionDao.
func NewSmartContractTransactionDao(dbname string) *SmartContractTransactionDao {
	dbName := dbname
	collectionName := "smartcontract_txs"
	i1 := mgo.Index{
		Key: []string{"smart_contract"},
	}

	i2 := mgo.Index{
		Key:    []string{"hash", "from", "to"},
		Unique: true,
	}

	err := db.Session.DB(dbName).C(collectionName).EnsureIndex(i1)
	if err != nil {
		panic(err)
	}

	err = db.Session.DB(dbName).C(collectionName).EnsureIndex(i2)
	if err != nil {
		panic(err)
	}

	return &SmartContractTransactionDao{
		DAO: &DAO{
			collectionName: collectionName,
			dbName:         dbName,
		},
	}
}

func (dao *SmartContractTransactionDao) InsertBulk(transactions []SmartContractTransaction) error {
	var ui []interface{}
	for _, t := range transactions {
		ui = append(ui, t)
	}

	err := dao.Create(ui...)

	return err
}
