package project

import (
	"context"

	"emperror.dev/errors"

	emperrorErr "emperror.dev/errors"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/bus"
	. "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store"
	"github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
	common "github.com/anhntbk08/gateway/internal/common"
	"github.com/globalsign/mgo/bson"
	"github.com/rs/xid"
)

// +kit:endpoint:errorStrategy=project

type Service interface {
	Create(ctx context.Context, name string) (project entity.Project, err error)
	List(ctx context.Context) (projects []entity.Project, err error)
	Update(ctx context.Context, project entity.Project) (err error)
	Delete(ctx context.Context, id bson.ObjectId) (err error)
	GetOne(ctx context.Context, id bson.ObjectId) (project entity.Project, err error)
	// Statistic(ctx context.Context, id string) (success bool, err error)
}

type service struct {
	db  *Mongo
	bus *bus.Bus
}

func NewService(db *Mongo, bus *bus.Bus) Service {
	return &service{
		db:  db,
		bus: bus,
	}
}

func (s service) checkUserExist(user string, action string) (*entity.User, error) {
	userDao, err := s.db.UserDao.IsExist(user)
	if err != nil {
		return userDao, errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT." + action + ".NOT_EXIST",
				"User not exist",
			},
		}})
	}

	return userDao, err
}

func (s service) Create(ctx context.Context, name string) (entity.Project, error) {
	user := ctx.Value("User").(string)
	userDao, err := s.checkUserExist(user, "CREATING")
	if err != nil {
		return entity.Project{}, err
	}

	project := entity.Project{
		Name:   name,
		Secret: xid.New().String(),
		User:   userDao.ID,
	}

	err = s.db.ProjectDao.Create(&project)

	return project, emperrorErr.WithStack(err)
}

func (s service) List(ctx context.Context) ([]entity.Project, error) {
	user := ctx.Value("User").(string)
	userDao, err := s.checkUserExist(user, "LISTING")
	if err != nil {
		return []entity.Project{}, err
	}

	projects := []entity.Project{}
	err = s.db.ProjectDao.Get(bson.M{
		"user": userDao.ID,
	}, 0, 100, &projects)

	return projects, emperrorErr.WithStack(err)
}

func (s service) Update(ctx context.Context, project entity.Project) (err error) {
	user := ctx.Value("User").(string)
	userDao, err := s.checkUserExist(user, "UPDATING")
	if err != nil {
		return err
	}

	// check belonging
	oldProject := &entity.Project{}
	err = s.db.ProjectDao.GetOne(bson.M{
		"_id":  project.ID,
		"user": userDao.ID,
	}, &oldProject)

	if err != nil {
		return errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT.UPDATING.NOT_FOUND",
				"Resource unauthenticated or not found",
			},
		}})
	}

	project.User = userDao.ID
	project.Secret = oldProject.Secret
	err = s.db.ProjectDao.Update(bson.M{
		"_id": project.ID,
	}, bson.M{
		"$set": project,
	})

	// if list watch address changes
	// need to notify subscriber and readd
	//oldProject.Addresses.WatchSmartContracts
	return err
}

func (s service) Delete(ctx context.Context, id bson.ObjectId) (err error) {
	user := ctx.Value("User").(string)
	userDao, err := s.checkUserExist(user, "DELETING")
	if err != nil {
		return err
	}

	// check belonging
	res := &entity.Project{}
	err = s.db.ProjectDao.GetOne(bson.M{
		"_id":  id,
		"user": userDao.ID,
	}, &res)

	if err != nil {
		return errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT.DELETING.NOT_FOUND",
				"Resource unauthenticated or not found",
			},
		}})
	}

	err = s.db.ProjectDao.RemoveItem(bson.M{
		"_id": id,
	})

	return err
}

func (s service) GetOne(ctx context.Context, id bson.ObjectId) (project entity.Project, err error) {
	user := ctx.Value("User").(string)
	userDao, err := s.checkUserExist(user, "GETONE")
	if err != nil {
		return entity.Project{}, err
	}

	// check belonging
	res := &entity.Project{}
	err = s.db.ProjectDao.GetOne(bson.M{
		"_id":  id,
		"user": userDao.ID,
	}, &res)

	if err != nil {
		return entity.Project{}, errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT.GETONE.NOT_FOUND",
				"Resource unauthenticated or not found",
			},
		}})
	}

	return *res, err
}
