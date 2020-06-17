package entity

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

// func (dao *ProjectDao) Upsert(project *Project) error {
// 	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{"_id": project.ID}, &ProjectRecordUpdate{cs})

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (dao *ProjectDao) Delete(project *Project) error {
// 	_, err := db.RemoveItem(dao.dbName, dao.collectionName, bson.M{"_id": project.ID})

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
