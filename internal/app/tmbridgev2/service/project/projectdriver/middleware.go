package projectdriver

import (
	"context"
	"strings"

	"emperror.dev/errors"
	projectService "github.com/anhntbk08/gateway/internal/app/tmbridgev2/service/project"
	entity "github.com/anhntbk08/gateway/internal/app/tmbridgev2/store/entity"
	"github.com/anhntbk08/gateway/internal/common"
	"github.com/globalsign/mgo/bson"
	validation "github.com/go-ozzo/ozzo-validation"
)

// Middleware is a service middleware.
type Middleware func(projectService.Service) projectService.Service

// defaultMiddleware helps implementing partial middleware, we use for validating fields
type defaultMiddleware struct {
	service projectService.Service
}

func (m defaultMiddleware) Create(ctx context.Context, name string) (entity.Project, error) {
	err := validation.Validate(name,
		validation.Required,
		validation.Length(3, 100),
	)
	if err != nil {
		return entity.Project{}, errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"name": {
				"NAME.MALFORM_NAME",
				err.Error(),
			},
		}})
	}
	return m.service.Create(ctx, name)
}
func (m defaultMiddleware) List(ctx context.Context) ([]entity.Project, error) {
	return m.service.List(ctx)
}

func checkDuplicatedAddresses(addresses []string) error {
	rndmap := make(map[string]bool)

	for _, address := range addresses {
		rndmap[strings.ToLower(address)] = true
	}

	if len(rndmap) < len(addresses) {
		return errors.New("Duplicated address")
	}

	return nil
}

func (m defaultMiddleware) Update(ctx context.Context, project entity.Project) error {
	err := validation.ValidateStruct(&project,
		validation.Field(&project.Name, validation.Length(3, 100)),
	)

	if err != nil {
		return errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT.MALFORM_PROJECT",
				err.Error(),
			},
		}})
	}

	err = checkDuplicatedAddresses(project.Addresses.WatchSmartContracts)

	if err != nil {
		return errors.WithStack(common.ValidationError{Violates: map[string][]string{
			"project": {
				"PROJECT.MALFORM_PROJECT_WATCH_SMART_CONTRACT",
				err.Error(),
			},
		}})
	}

	return m.service.Update(ctx, project)
}

func (m defaultMiddleware) Delete(ctx context.Context, projectID bson.ObjectId) error {
	return m.service.Delete(ctx, projectID)
}

func (m defaultMiddleware) GetOne(ctx context.Context, projectID bson.ObjectId) (entity.Project, error) {
	return m.service.GetOne(ctx, projectID)
}

// LoggingMiddleware is a service level logging middleware.
func LoggingMiddleware(logger projectService.Logger) Middleware {
	return func(next projectService.Service) projectService.Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   projectService.Service
	logger projectService.Logger
}

func (mw loggingMiddleware) Create(ctx context.Context, name string) (entity.Project, error) {
	logger := mw.logger.WithContext(ctx)
	resp, err := mw.next.Create(ctx, name)
	if err != nil {
		return entity.Project{}, err
	}

	logger.Info("Created project ", map[string]interface{}{"name": resp.Name, "id": resp.ID})

	return resp, err
}

func (mw loggingMiddleware) List(ctx context.Context) ([]entity.Project, error) {
	resp, err := mw.next.List(ctx)
	if err != nil {
		return []entity.Project{}, err
	}

	return resp, err
}

func (mw loggingMiddleware) Update(ctx context.Context, project entity.Project) error {
	logger := mw.logger.WithContext(ctx)
	err := mw.next.Update(ctx, project)

	if err == nil {
		logger.Info("Updated project ", map[string]interface{}{"name": project.Name, "id": project.ID})
	}

	return err
}

func (mw loggingMiddleware) Delete(ctx context.Context, id bson.ObjectId) error {
	logger := mw.logger.WithContext(ctx)
	err := mw.next.Delete(ctx, id)

	if err == nil {
		logger.Info("Deleted project ", map[string]interface{}{"id": id})
	}

	return err
}

func (mw loggingMiddleware) GetOne(ctx context.Context, id bson.ObjectId) (entity.Project, error) {
	return mw.next.GetOne(ctx, id)
}

// InstrumentationMiddleware is a service level instrumentation middleware.
func InstrumentationMiddleware() Middleware {
	return func(next projectService.Service) projectService.Service {
		return instrumentationMiddleware{
			Service: defaultMiddleware{next},
			next:    next,
		}
	}
}

type instrumentationMiddleware struct {
	projectService.Service
	next projectService.Service
}
