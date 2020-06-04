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

// func (dao *ProjectDao) Create(project *Project) error {
// 	if !project.User.Valid() {
// 		return errors.New("USER_IS_NOT_VALID")
// 	}
// 	_, err := db.Create(dao.dbName, dao.collectionName, project)

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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
