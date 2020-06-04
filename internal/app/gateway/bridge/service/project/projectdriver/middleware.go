package projectdriver

import (
	"context"
	"errors"

	projectService "github.com/anhntbk08/gateway/internal/app/gateway/bridge/service/project"
	entity "github.com/anhntbk08/gateway/internal/app/gateway/store/entity"
)

// Middleware is a service middleware.
type Middleware func(projectService.Service) projectService.Service

// defaultMiddleware helps implementing partial middleware.
type defaultMiddleware struct {
	service projectService.Service
}

func (m defaultMiddleware) Create(ctx context.Context, name, user string) (entity.Project, error) {
	return m.service.Create(ctx, name, user)
}

func (m defaultMiddleware) Delete(ctx context.Context, id string) (bool, error) {
	return false, errors.New("NOT_IMPLEMENTED_YET")
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

func (mw loggingMiddleware) Create(ctx context.Context, name, user string) (entity.Project, error) {
	// logger := mw.logger.WithContext(ctx)

	// logger.Info(request.Address + " trying to create project ")

	resp, err := mw.next.Create(ctx, name, user)
	if err != nil {
		return entity.Project{}, err
	}

	// logger.Info("Logged in", map[string]interface{}{"address": request.Address})

	return resp, err
}

func (mw loggingMiddleware) Delete(ctx context.Context, id string) (bool, error) {
	return false, errors.New("NOT_IMPLEMENTED_YET")
}

// Business metrics
// nolint: gochecknoglobals,lll
// var (
// 	CreatedTodoItemCount  = stats.Int64("created_todo_item_count", "Number of todo items created", stats.UnitDimensionless)
// 	CompleteTodoItemCount = stats.Int64("complete_todo_item_count", "Number of todo items marked complete", stats.UnitDimensionless)
// )

// // nolint: gochecknoglobals
// var (
// 	CreatedTodoItemCountView = &view.View{
// 		Name:        "todo_item_created_count",
// 		Description: "Count of todo items created",
// 		Measure:     CreatedTodoItemCount,
// 		Aggregation: view.Count(),
// 	}

// 	CompleteTodoItemCountView = &view.View{
// 		Name:        "todo_item_complete_count",
// 		Description: "Count of todo items complete",
// 		Measure:     CompleteTodoItemCount,
// 		Aggregation: view.Count(),
// 	}
// )

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

// func (mw instrumentationMiddleware) RequestToken(ctx context.Context, request projectService.RqTokenData) (projectService.Token, error) {
// 	token, err := mw.next.RequestToken(ctx, request)
// 	if err != nil {
// 		return token, err
// 	}

// 	if span := trace.FromContext(ctx); span != nil {
// 		span.AddAttributes(trace.StringAttribute("token", token.Token))
// 	}

// 	// stats.Record(ctx, CreatedTodoItemCount.M(1))

// 	return token, nil
// }
