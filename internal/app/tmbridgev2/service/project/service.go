package project

import (
	"context"
	"fmt"

	"emperror.dev/errors"

	emperrorErr "emperror.dev/errors"
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
	// Statistic(ctx context.Context, id string) (success bool, err error)
}

type service struct {
	db *Mongo
}

func NewService(db *Mongo) Service {
	return &service{
		db: db,
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
		Name: name,
		Keys: entity.Keys{
			ID:     xid.New().String(),
			Secret: xid.New().String(),
		},
		User: userDao.ID,
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
		"user_id": userDao.ID,
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
	res := &entity.Project{}
	err = s.db.ProjectDao.GetOne(bson.M{
		"_id":     project.ID,
		"user_id": userDao.ID,
	}, &res)

	if err != nil {
		return errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT.UPDATING.UNAUTHENTICATED",
				"Resource unauthenticated",
			},
		}})
	}

	project.User = userDao.ID
	err = s.db.ProjectDao.Update(bson.M{
		"_id": project.ID,
	}, project)

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
		"_id":     id,
		"user_id": userDao.ID,
	}, &res)

	if err != nil {
		return errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT.DELETING.UNAUTHENTICATED",
				"Resource unauthenticated",
			},
		}})
	}

	err = s.db.ProjectDao.RemoveItem(bson.M{
		"_id": id,
	})
	fmt.Println("err ", err)
	return err
}

func (s service) Statistic(ctx context.Context, id string) (success bool, err error) {
	return false, errors.New("Not implemented yet")
}
