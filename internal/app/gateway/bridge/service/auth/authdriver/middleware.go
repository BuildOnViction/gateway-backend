package authdriver

import (
	"context"

	"go.opencensus.io/trace"

	authService "github.com/anhntbk08/gateway/internal/app/gateway/bridge/service/auth"
)

// Middleware is a service middleware.
type Middleware func(authService.Service) authService.Service

// defaultMiddleware helps implementing partial middleware.
type defaultMiddleware struct {
	service authService.Service
}

func (m defaultMiddleware) RequestToken(ctx context.Context, request authService.RqTokenData) (authService.Token, error) {
	return m.service.RequestToken(ctx, request)
}

func (m defaultMiddleware) Login(ctx context.Context, request authService.Token) (bool, error) {
	return m.service.Login(ctx, request)
}

// LoggingMiddleware is a service level logging middleware.
func LoggingMiddleware(logger authService.Logger) Middleware {
	return func(next authService.Service) authService.Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   authService.Service
	logger authService.Logger
}

func (mw loggingMiddleware) RequestToken(ctx context.Context, request authService.RqTokenData) (authService.Token, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Info("Request token")

	token, err := mw.next.RequestToken(ctx, request)
	if err != nil {
		return token, err
	}

	logger.Info("Requested token", map[string]interface{}{"token": token.Token})

	return token, err
}

func (mw loggingMiddleware) Login(ctx context.Context, request authService.Token) (bool, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Info(request.Address + " trying to login in ")

	resp, err := mw.next.Login(ctx, request)
	if err != nil {
		return false, err
	}

	logger.Info("Logged in", map[string]interface{}{"address": request.Address})

	return resp, err
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
	return func(next authService.Service) authService.Service {
		return instrumentationMiddleware{
			Service: defaultMiddleware{next},
			next:    next,
		}
	}
}

type instrumentationMiddleware struct {
	authService.Service
	next authService.Service
}

func (mw instrumentationMiddleware) RequestToken(ctx context.Context, request authService.RqTokenData) (authService.Token, error) {
	token, err := mw.next.RequestToken(ctx, request)
	if err != nil {
		return token, err
	}

	if span := trace.FromContext(ctx); span != nil {
		span.AddAttributes(trace.StringAttribute("token", token.Token))
	}

	// stats.Record(ctx, CreatedTodoItemCount.M(1))

	return token, nil
}
