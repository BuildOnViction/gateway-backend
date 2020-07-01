package entity

import "github.com/globalsign/mgo/bson"

type ProjectDao struct {
	*DAO
}

// NewProjectDao returns a new instance of ProjectDao.
func NewProjectDao(dbname string) *ProjectDao {
	dbName := dbname
	collectionName := "projects"
	return &ProjectDao{
		DAO: &DAO{
			collectionName: collectionName,
			dbName:         dbName,
		},
	}
}

func (dao *ProjectDao) ExistToken(apitoken string) (*Project, error) {
	res := &Project{}
	err := dao.GetOne(bson.M{"keys.id": apitoken}, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// func (dao *ProjectDao) Delete(project *Project) error {
// 	_, err := db.RemoveItem(dao.dbName, dao.collectionName, bson.M{"_id": project.ID})

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
