package project

import (
	"context"
	"errors"

	. "github.com/anhntbk08/gateway/internal/app/gateway/store"
	"github.com/anhntbk08/gateway/internal/app/gateway/store/entity"
	// . "github.com/anhntbk08/gateway/internal/common"
	// . "github.com/anhntbk08/gateway/internal/app/gateway/store/entity"
)

// +kit:endpoint:errorStrategy=project

type Service interface {
	Create(ctx context.Context, name string, user string) (project entity.Project, err error)
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

// func genProjectID(name, user) {

// }

func (s service) Create(ctx context.Context, name string, user string) (project entity.Project, err error) {
	return entity.Project{
			Name: name,
			Keys: entity.Keys{
				ID: "12313123123",
			},
		}, s.db.ProjectDao.Create(&entity.Project{
			Name: name,
			Keys: entity.Keys{
				ID: "12313123123",
			},
		})
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
