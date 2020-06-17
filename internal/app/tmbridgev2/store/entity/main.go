package entity

import (
	"reflect"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Database struct contains the pointer to mgo.session
// It is a wrapper over mgo to help utilize mgo connection pool
type Database struct {
	Session *mgo.Session
}

type DAOInterface interface {
	Create(data ...interface{}) (err error)
	Watch(options mgo.ChangeStreamOptions) (*mgo.ChangeStream, *mgo.Session, error)
	Count(query interface{}) (int, error)
	GetByID(id bson.ObjectId, response interface{}) (err error)
	Get(query interface{}, offset, limit int, response interface{}) (err error)
	GetOne(query interface{}, response interface{}) (err error)
	Query(query interface{}, selector interface{}, offset, limit int, response interface{}) (err error)
	GetAndSort(query interface{}, sort []string, offset, limit int, response interface{}) (err error)
	GetEx(query interface{}, sort []string, offset, limit int, response interface{}) (count int, err error)
	GetSortOne(query interface{}, sort []string, response interface{}) (err error)
	Update(query interface{}, update interface{}) error
	Upsert(query interface{}, update interface{}) (interface{}, error)
	UpdateAll(query interface{}, update interface{}) error
	ChangeAll(query interface{}, update interface{}) (*mgo.ChangeInfo, error)
	FindAndModify(query interface{}, change mgo.Change, response interface{}) error
	AggregateEx(query []bson.M, response interface{}) error
	Aggregate(query []bson.M, response interface{}) error
	Remove(query []bson.M) error
	RemoveItem(query interface{}) error
	RemoveAll(query interface{}) error
	DropCollection() error
}

type DAO struct {
	DAOInterface
	collectionName string
	dbName         string
}

// Global instance of Database struct for singleton use
var db *Database

// InitSession initializes a new session with mongodb
func InitSession(url string, session *mgo.Session) (*mgo.Session, error) {
	if db == nil {
		if session == nil {
			db1, err := mgo.Dial(url)
			if err != nil {
				return nil, err
			}

			session = db1
		}

		db = &Database{session}
	}
	return db.Session, nil
}

func (d *Database) InitDatabase(session *mgo.Session) {
	d.Session = session
}

func (d *Database) GetCollection(dbName, collection string) *mgo.Collection {
	s := d.Session

	return s.DB(dbName).C(collection)
}

// Create is a wrapper for mgo.Insert function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Create(dbName, collection string, data ...interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Insert(data...)
	return
}

func (d *Database) Watch(dbName, collection string, options mgo.ChangeStreamOptions) (*mgo.ChangeStream, *mgo.Session, error) {
	sc := d.Session.Copy()

	pipeline := []bson.M{}

	ct, err := sc.DB(dbName).C(collection).Watch(pipeline, options)

	return ct, sc, err
}

func (d *Database) Count(dbName, collection string, query interface{}) (int, error) {
	sc := d.Session.Copy()
	defer sc.Close()

	return sc.DB(dbName).C(collection).Find(query).Count()
}

// GetByID is a wrapper for mgo.FindId function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) GetByID(dbName, collection string, id bson.ObjectId, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).FindId(id).One(response)
	return
}

// Get is a wrapper for mgo.Find function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Get(dbName, collection string, query interface{}, offset, limit int, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).Skip(offset).Limit(limit).All(response)
	return
}

// GetOne return one document
func (d *Database) GetOne(dbName, collection string, query interface{}, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).One(response)
	return
}

func (d *Database) Query(dbName, collection string, query interface{}, selector interface{}, offset, limit int, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).Skip(offset).Limit(limit).Select(selector).All(response)
	return
}

// GetAndSort is a wrapper for mgo.Find function with SORT function in pipeline.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) GetAndSort(dbName, collection string, query interface{}, sort []string, offset, limit int, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()

	err = sc.DB(dbName).C(collection).Find(query).Sort(sort...).Skip(offset).Limit(limit).All(response)
	return
}

// GetEx extend get function
func (d *Database) GetEx(dbName, collection string, query interface{}, sort []string, offset, limit int, response interface{}) (count int, err error) {
	sc := d.Session.Copy()
	defer sc.Close()
	cursor := sc.DB(dbName).C(collection).Find(query).Sort(sort...)
	c, _ := cursor.Count()
	err = cursor.Skip(offset).Limit(limit).All(response)
	return c, err
}

// GetSortOne is a wrapper for mgo.Find function with SORT function in pipeline.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) GetSortOne(dbName, collection string, query interface{}, sort []string, response interface{}) (err error) {
	sc := d.Session.Copy()
	defer sc.Close()
	err = sc.DB(dbName).C(collection).Find(query).Sort(sort...).One(response)
	return
}

// Update is a wrapper for mgo.Update function.
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Update(dbName, collection string, query interface{}, update interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	err := sc.DB(dbName).C(collection).Update(query, update)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) Upsert(dbName, collection string, query interface{}, update interface{}) (interface{}, error) {
	sc := d.Session.Copy()
	defer sc.Close()

	changed, err := sc.DB(dbName).C(collection).Upsert(query, update)
	if err != nil {
		return nil, err
	}
	return changed.UpsertedId, nil
}

func (d *Database) UpdateAll(dbName, collection string, query interface{}, update interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	_, err := sc.DB(dbName).C(collection).UpdateAll(query, update)
	if err != nil {
		return err
	}

	return nil
}

// ChangeAll update all document and return change information
func (d *Database) ChangeAll(dbName, collection string, query interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	sc := d.Session.Copy()
	defer sc.Close()

	changeInfo, err := sc.DB(dbName).C(collection).UpdateAll(query, update)
	if err != nil {
		return nil, err
	}

	return changeInfo, nil
}

func (d *Database) FindAndModify(dbName, collection string, query interface{}, change mgo.Change, response interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	_, err := sc.DB(dbName).C(collection).Find(query).Apply(change, response)
	if err != nil {
		return err
	}

	return nil
}

// AggregateEx add collation
func (d *Database) AggregateEx(dbName, collection string, query []bson.M, response interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	result := reflect.ValueOf(response).Interface()
	c := mgo.Collation{
		NumericOrdering: true,
		Locale:          "en_US",
	}
	err := sc.DB(dbName).C(collection).Pipe(query).Collation(&c).All(result)
	if err != nil {
		return err
	}

	return nil
}

// Aggregate is a wrapper for mgo.Pipe function.
// It is used to make mongo aggregate pipeline queries
// It creates a copy of session initialized, sends query over this session
// and returns the session to connection pool
func (d *Database) Aggregate(dbName, collection string, query []bson.M, response interface{}) error {
	return d.AggregateEx(dbName, collection, query, response)
}

// Remove removes one document matching a certain query
func (d *Database) Remove(dbName, collection string, query []bson.M) error {
	sc := d.Session.Copy()
	defer sc.Close()

	err := sc.DB(dbName).C(collection).Remove(query)
	if err != nil {
		return err
	}

	return nil
}

// Remove removes one document matching a certain query
func (d *Database) RemoveItem(dbName, collection string, query interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	err := sc.DB(dbName).C(collection).Remove(query)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAll removes all the documents from a collection matching a certain query
func (d *Database) RemoveAll(dbName, collection string, query interface{}) error {
	sc := d.Session.Copy()
	defer sc.Close()

	_, err := sc.DB(dbName).C(collection).RemoveAll(query)
	if err != nil {
		return err
	}

	return nil
}

// DropCollection drops all the documents in a collection
func (d *Database) DropCollection(dbName, collection string) error {
	sc := d.Session.Copy()
	defer sc.Close()

	err := sc.DB(dbName).C(collection).DropCollection()
	if err != nil {
		return err
	}

	return nil
}

func (d *DAO) Create(data ...interface{}) (err error) {
	return db.Create(d.dbName, d.collectionName, data...)
}

func (d *DAO) Watch(options mgo.ChangeStreamOptions) (*mgo.ChangeStream, *mgo.Session, error) {
	return db.Watch(d.dbName, d.collectionName, options)
}

func (d *DAO) Count(query interface{}) (int, error) {
	return db.Count(d.dbName, d.collectionName, query)
}

func (d *DAO) GetByID(id bson.ObjectId, response interface{}) (err error) {
	return db.GetByID(d.dbName, d.collectionName, id, response)
}

func (d *DAO) Get(query interface{}, offset, limit int, response interface{}) (err error) {
	return db.Get(d.dbName, d.collectionName, query, offset, limit, response)
}

func (d *DAO) GetOne(query interface{}, response interface{}) (err error) {
	return db.GetOne(d.dbName, d.collectionName, query, response)
}

func (d *DAO) Query(query interface{}, selector interface{}, offset, limit int, response interface{}) (err error) {
	return db.Query(d.dbName, d.collectionName, query, selector, offset, limit, response)
}

func (d *DAO) GetAndSort(query interface{}, sort []string, offset, limit int, response interface{}) (err error) {
	return db.GetAndSort(d.dbName, d.collectionName, query, sort, offset, limit, response)
}

func (d *DAO) GetEx(query interface{}, sort []string, offset, limit int, response interface{}) (count int, err error) {
	return db.GetEx(d.dbName, d.collectionName, query, sort, offset, limit, response)
}

func (d *DAO) GetSortOne(query interface{}, sort []string, response interface{}) (err error) {
	return db.GetSortOne(d.dbName, d.collectionName, query, sort, response)
}

func (d *DAO) Update(query interface{}, update interface{}) error {
	return db.Update(d.dbName, d.collectionName, query, update)
}

func (d *DAO) Upsert(query interface{}, update interface{}) (interface{}, error) {
	return db.Upsert(d.dbName, d.collectionName, query, update)
}

func (d *DAO) UpdateAll(query interface{}, update interface{}) error {
	return db.UpdateAll(d.dbName, d.collectionName, query, update)
}

// ChangeAll update all document and return change information
func (d *DAO) ChangeAll(query interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	return db.ChangeAll(d.dbName, d.collectionName, query, update)
}

func (d *DAO) FindAndModify(query interface{}, change mgo.Change, response interface{}) error {
	return db.FindAndModify(d.dbName, d.collectionName, query, change, response)
}

func (d *DAO) AggregateEx(query []bson.M, response interface{}) error {
	return db.AggregateEx(d.dbName, d.collectionName, query, response)
}

func (d *DAO) Aggregate(query []bson.M, response interface{}) error {
	return db.AggregateEx(d.dbName, d.collectionName, query, response)
}

// Remove removes one document matching a certain query
func (d *DAO) Remove(query []bson.M) error {
	return db.Remove(d.dbName, d.collectionName, query)
}

// Remove removes one document matching a certain query
func (d *DAO) RemoveItem(query interface{}) error {
	return db.RemoveItem(d.dbName, d.collectionName, query)
}

func (d *DAO) RemoveAll(query interface{}) error {
	return db.RemoveAll(d.dbName, d.collectionName, query)
}

func (d *DAO) DropCollection() error {
	return db.DropCollection(d.dbName, d.collectionName)
}
