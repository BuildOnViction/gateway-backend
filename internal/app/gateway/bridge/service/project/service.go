package project

import (
	"context"
	"errors"
	"fmt"

	emperrorErr "emperror.dev/errors"
	. "github.com/anhntbk08/gateway/internal/app/gateway/store"
	"github.com/anhntbk08/gateway/internal/app/gateway/store/entity"
	"github.com/rs/xid"
	// . "github.com/anhntbk08/gateway/internal/common"
	// . "github.com/anhntbk08/gateway/internal/app/gateway/store/entity"
)

// +kit:endpoint:errorStrategy=project

type Service interface {
	Create(ctx context.Context, name string) (project entity.Project, err error)
	// List(ctx context.Context) (projects []entity.Project, err error)
	// Update(ctx context.Context, project entity.Project) (success bool, err error)
	// Delete(ctx context.Context, id string) (success bool, err error)
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

func (s service) Create(ctx context.Context, name string) (entity.Project, error) {
	user := ctx.Value("User").(string)

	userDao, err := s.db.UserDao.IsExist(user)
	if err != nil {
		return entity.Project{}, emperrorErr.WithStack(errors.New("User not exists"))
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
	fmt.Println(err)

	return project, emperrorErr.WithStack(err)
}

func (s service) List(ctx context.Context) (projects []entity.Project, err error) {
	return []entity.Project{}, errors.New("Not implemented yet")
}

func (s service) Update(ctx context.Context, project entity.Project) (success bool, err error) {
	return false, errors.New("Not implemented yet")
}

func (s service) Delete(ctx context.Context, id string) (success bool, err error) {
	return false, errors.New("Not implemented yet")
}

func (s service) Statistic(ctx context.Context, id string) (success bool, err error) {
	return false, errors.New("Not implemented yet")
}
